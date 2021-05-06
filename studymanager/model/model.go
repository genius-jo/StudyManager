package model

import (
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassWord string `json:"pass_word"`
}

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassWord string `json:"pass_word"`
}

type Login struct {
	Email    string `json:"email"`
	PassWord string `json:"pass_word"`
}

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Created   time.Time `json:"created_at"`
}

type UserDBHandler interface {
	GetUsers() []*User
	AddUser(register Register) *User
	Close()
}

type LoginDBHandler interface {
	CheckUser(login Login) *User
	Close()
}

type TodoDBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

func NewUserDBHandler(filepath string) UserDBHandler {
	return newUserSqlHandler(filepath)
}

func NewLoginDBHandler(filepath string) LoginDBHandler {
	return newLoginSqlHandler(filepath)
}

func NewTodoDBHandler(filepath string) TodoDBHandler {
	return newTodoSqlHandler(filepath)
}
