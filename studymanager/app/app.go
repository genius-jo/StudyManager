package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"studymanager/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))) //쿠키 스토어를 만드는데, 암호화 하는 키는 환경변수에서 SESSION_KEY
var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db  model.UserDBHandler
	db2 model.LoginDBHandler
}

func getSessionID(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]

	if val == nil { //로그인이 안 되어 있을때
		return ""
	}

	return val.(string)
}

func (ah *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) getUserListHandler(w http.ResponseWriter, r *http.Request) {
	list := ah.db.GetUsers()
	rd.JSON(w, http.StatusOK, list)
}

func (ah *AppHandler) addUserHandler(w http.ResponseWriter, r *http.Request) {

	var register model.Register = model.Register{r.FormValue("name"), r.FormValue("email"), r.FormValue("pass_word")}
	user := ah.db.AddUser(register)
	rd.JSON(w, http.StatusCreated, user)
}

func (ah *AppHandler) Close() {
	ah.db.Close()
	ah.db2.Close()
}

func (ah *AppHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var login model.Login = model.Login{r.FormValue("email"), r.FormValue("pass_word")}
	user := ah.db2.CheckUser(login)

	if user != nil {
		rd.JSON(w, http.StatusOK, user)

		//세션 쿠키에 유저 정보를 저장
		session, err := store.Get(r, "session") //세션을 가져온다
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["id"] = user.Id
		session.Values["name"] = user.Name

		err = session.Save(r, w) //쿠키에 저장
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//메인 화면으로 리다이렉트
		//http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)

	} else {
		errMsg := fmt.Sprintf("invalid login info")
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		//rd.JSON(w, http.StatusBadRequest, user)
	}

}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//요청한 url이 로그인일때는 next로
	if strings.Contains(r.URL.Path, "/login.html") || strings.Contains(r.URL.Path, "/login") || strings.Contains(r.URL.Path, "/users") {
		next(w, r)
		return
	}

	//로그인이 되어 있을때
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	//로그인이 되어 있지 않을때 로그인 화면으로 리다이렉트
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func MakeHandler(filepath string) *AppHandler {

	router := mux.NewRouter()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin), negroni.NewStatic(http.Dir("public"))) //미들웨어 추가, 미들웨어는 체인 형식

	n.UseHandler(router)

	ah := &AppHandler{
		Handler: n,
		db:      model.NewUserDBHandler(filepath),
		db2:     model.NewLoginDBHandler(filepath),
	}

	router.HandleFunc("/users", ah.getUserListHandler).Methods("GET")
	router.HandleFunc("/users", ah.addUserHandler).Methods("POST")
	router.HandleFunc("/", ah.indexHandler)

	router.HandleFunc("/login", ah.loginHandler).Methods("POST")

	return ah
}
