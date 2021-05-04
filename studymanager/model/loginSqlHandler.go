package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type loginSqlHandler struct {
	db *sql.DB
}

//login한 정보가 DB에 있는 정보와 맞는지 확인
func (lsh *loginSqlHandler) CheckUser(login Login) *User {
	var user User
	rows, err := lsh.db.Query("SELECT id, name, email, pass_word FROM users WHERE email=?", login.Email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.PassWord)
	}

	log.Printf("%s %s", user.Email, user.Name)

	if user.PassWord != login.PassWord {
		return nil
	}

	return &user
}

func (ush *loginSqlHandler) Close() {
	ush.db.Close()
}

func newLoginSqlHandler(filepath string) LoginDBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	return &loginSqlHandler{db: database}
}
