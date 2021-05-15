package app

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))) //쿠키 스토어를 만드는데, 암호화 하는 키는 환경변수에서 SESSION_KEY

var getSessionID = func(r *http.Request) string {

	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]

	if val == nil { //로그인이 안 되어 있을때
		log.Println("ID nil")
		return ""
	}

	ival := val.(int)
	sval := strconv.Itoa(ival)

	return sval
}

func getSessionName(r *http.Request) string {

	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["name"]

	if val == nil { //로그인이 안 되어 있을때
		log.Println("name nil")
		return ""
	}

	sval := val.(string)

	return sval
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//요청한 url이 로그인일때는 next로
	if strings.Contains(r.URL.Path, "/login") || strings.Contains(r.URL.Path, "/users") {
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
