package app

import (
	"net/http"
	"studymanager/model"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db model.UserDBHandler
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

func (a *AppHandler) Close() {
	a.db.Close()
}

func MakeHandler(filepath string) *AppHandler {

	router := mux.NewRouter()

	ah := &AppHandler{
		Handler: router,
		db:      model.NewUserDBHandler(filepath),
	}

	router.HandleFunc("/users", ah.getUserListHandler).Methods("GET")
	router.HandleFunc("/users", ah.addUserHandler).Methods("POST")
	router.HandleFunc("/", ah.indexHandler)

	return ah
}
