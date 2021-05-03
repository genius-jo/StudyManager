package app

import (
	"net/http"
	"studymanager/model"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rd *render.Render

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func getUserListHandler(w http.ResponseWriter, r *http.Request) {
	list := model.GetUsers()
	rd.JSON(w, http.StatusOK, list)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	var register model.Register = model.Register{r.FormValue("name"), r.FormValue("email"), r.FormValue("pass_word")}
	user := model.AddUser(register)
	rd.JSON(w, http.StatusCreated, user)

}

func MakeHandler() http.Handler {
	rd = render.New()
	router := mux.NewRouter()

	router.HandleFunc("/users", getUserListHandler).Methods("GET")
	router.HandleFunc("/users", addUserHandler).Methods("POST")
	router.HandleFunc("/", indexHandler)

	return router
}
