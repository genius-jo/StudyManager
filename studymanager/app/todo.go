package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Success struct {
	Success bool `json:"success"`
}

func (ah *AppHandler) indexTodoHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {

	sessionId := getSessionID(r)
	list := ah.tdb.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)

}

func (ah *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {

	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := ah.tdb.AddTodo(name, sessionId)
	rd.JSON(w, http.StatusCreated, todo)

}

func (ah *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := ah.tdb.RemoveTodo(id)

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
	ok := ah.tdb.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}

}
