package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type userSqlHandler struct {
	db *sql.DB
}

func (ush *userSqlHandler) GetUsers() []*User {
	users := []*User{}
	rows, err := ush.db.Query("SELECT id, name, email, pass_word FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.PassWord)
		users = append(users, &user)
	}

	return users
}

func (ush *userSqlHandler) AddUser(register Register) *User {
	stmt, err := ush.db.Prepare("INSERT INTO users (name, email, pass_word) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}

	log.Println(register.Name)
	log.Println(register.Email)
	log.Println(register.PassWord)

	rst, err := stmt.Exec(register.Name, register.Email, register.PassWord) //(?,?,?)에 대응하는것
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId() //방금 발급된 ID
	var user User
	user.Id = int(id) //int64타입을 int로 바꿈
	user.Name = register.Name
	user.Email = register.Email
	user.PassWord = register.PassWord

	return &user
}

func (ush *userSqlHandler) Close() {
	ush.db.Close()
}

func newUserSqlHandler(filepath string) UserDBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			pass_word TEXT
			)`)
	statement.Exec()
	return &userSqlHandler{db: database}
}
