package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	db3 model.TodoDBHandler
}

/*
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
*/

var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]

	log.Println("---sessionID---")
	log.Println(val)
	log.Println("---sessionID---")

	if val == nil { //로그인이 안 되어 있을때
		log.Println("ID nil")
		return ""
	}

	ret := val.(int)

	log.Println("---sessionID int---")
	log.Println(ret)
	log.Println("---sessionID int---")

	sret := strconv.Itoa(ret)

	log.Println("---sessionID string---")
	log.Println(sret)
	log.Println("---sessionID string---")

	return sret //"12" //val.(string)
}

func (ah *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) indexTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("here todo")
	log.Println(r.Method)
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
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
	ah.db3.Close()
}

func (ah *AppHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var login model.Login = model.Login{r.FormValue("email"), r.FormValue("pass_word")}
	user := ah.db2.CheckUser(login)

	if user != nil {
		//rd.JSON(w, http.StatusOK, user)

		//세션 쿠키에 유저 정보를 저장
		session, err := store.Get(r, "session") //세션을 가져온다
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(user.Id)
		log.Println(user.Name)
		log.Println(user.Email)
		log.Println(user.PassWord)

		session.Values["id"] = user.Id
		session.Values["name"] = user.Name

		err = session.Save(r, w) //쿠키에 저장
		if err != nil {
			log.Println("session error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//log.Println("redirect")
		//todo로 리다이렉트
		//http.Redirect(w, r, "/todo", 301)

	} else {
		errMsg := fmt.Sprintf("invalid login info")
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		rd.JSON(w, http.StatusBadRequest, user)
	}

}

func (ah *AppHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session") //세션을 가져온다
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//세션 초기화
	session.Values["id"] = nil
	session.Values["name"] = nil

	err = session.Save(r, w) //쿠키에 저장
	if err != nil {
		log.Println("session error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//로그인화면으로 리다이렉트
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {

	sessionId := getSessionID(r)
	list := ah.db3.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)
}

func (ah *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := ah.db3.AddTodo(name, sessionId)
	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func (ah *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := ah.db3.RemoveTodo(id)

	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (ah *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := ah.db3.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func MakeHandler(filepath string) *AppHandler {

	router := mux.NewRouter()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger() /*negroni.HandlerFunc(CheckSignin),*/, negroni.NewStatic(http.Dir("public"))) //미들웨어 추가, 미들웨어는 체인 형식

	n.UseHandler(router)

	ah := &AppHandler{
		Handler: n,
		db:      model.NewUserDBHandler(filepath),
		db2:     model.NewLoginDBHandler(filepath),
		db3:     model.NewTodoDBHandler(filepath),
	}

	router.HandleFunc("/users", ah.getUserListHandler).Methods("GET")
	router.HandleFunc("/users", ah.addUserHandler).Methods("POST")
	router.HandleFunc("/", ah.indexHandler)

	router.HandleFunc("/login", ah.loginHandler) //.Methods("POST")
	router.HandleFunc("/logout", ah.logoutHandler)

	router.HandleFunc("/todo", ah.indexTodoHandler).Methods("GET")
	router.HandleFunc("/todos", ah.getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", ah.addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", ah.removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", ah.completeTodoHandler).Methods("GET")

	return ah
}
