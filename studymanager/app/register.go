package app

import (
	"net/http"
	"studymanager/model"
)

func (ah *AppHandler) addUserHandler(w http.ResponseWriter, r *http.Request) {

	var register model.Register = model.Register{r.FormValue("name"), r.FormValue("email"), r.FormValue("pass_word")}

	user := ah.udb.AddUser(register)
	rd.JSON(w, http.StatusCreated, user)

}
