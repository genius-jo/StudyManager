package app

import (
	_ "fmt"
	_ "log"
	"net/http"
	_ "os"
	_ "strconv"
	_ "strings"
	"studymanager/model"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	udb model.UserDBHandler
	ldb model.LoginDBHandler
	tdb model.TodoDBHandler
}

func (ah *AppHandler) Close() {
	ah.udb.Close()
	ah.ldb.Close()
	ah.tdb.Close()
}

func (ah *AppHandler) getUserListHandler(w http.ResponseWriter, r *http.Request) {
	list := ah.udb.GetUsers()
	rd.JSON(w, http.StatusOK, list)
}

func MakeHandler(filepath string) *AppHandler {

	router := mux.NewRouter()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin), negroni.NewStatic(http.Dir("public"))) //미들웨어 추가, 미들웨어는 체인 형식
	n.UseHandler(router)

	ah := &AppHandler{
		Handler: n,
		udb:     model.NewUserDBHandler(filepath),
		ldb:     model.NewLoginDBHandler(filepath),
		tdb:     model.NewTodoDBHandler(filepath),
	}

	//users관련
	router.HandleFunc("/users", ah.getUserListHandler).Methods("GET")
	router.HandleFunc("/users", ah.addUserHandler).Methods("POST")
	router.HandleFunc("/home", ah.indexHandler)

	//login관련
	router.HandleFunc("/login", ah.loginHandler) //.Methods("POST")
	router.HandleFunc("/logout", ah.logoutHandler)

	//todo관련
	router.HandleFunc("/todo", ah.indexTodoHandler).Methods("GET")
	router.HandleFunc("/todos", ah.getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", ah.addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", ah.removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", ah.completeTodoHandler).Methods("GET")

	//채팅 관련
	router.HandleFunc("/chat", ah.indexChatHandler).Methods("GET")
	router.HandleFunc("/messages", ah.messageHandler).Methods("POST")
	router.Handle("/stream", Es) //클라이언트가 /stream으로 요청할때 커넥트를 맺는다
	//router.HandleFunc("/user", ah.addNameHandler)
	//router.HandleFunc("/user", ah.leftUserHandler).Methods("DELETE")

	return ah
}
