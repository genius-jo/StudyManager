package app

import (
	"fmt"
	"log"
	"net/http"
	"studymanager/model"
)

func (ah *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) loginHandler(w http.ResponseWriter, r *http.Request) {

	var login model.Login = model.Login{r.FormValue("email"), r.FormValue("pass_word")}
	user := ah.ldb.CheckUser(login)

	if user != nil {

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
			log.Println("session error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
