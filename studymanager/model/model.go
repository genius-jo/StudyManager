package model

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

type UserDBHandler interface {
	GetUsers() []*User
	AddUser(register Register) *User
	Close()
}

func NewUserDBHandler(filepath string) UserDBHandler {
	return newUserSqlHandler(filepath)
}
