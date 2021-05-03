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

var userMap map[int]*User

func init() {
	userMap = make(map[int]*User)
}

func GetUsers() []*User {
	list := []*User{}
	for _, v := range userMap {
		list = append(list, v)
	}
	return list
}

func AddUser(register Register) *User {
	id := len(userMap) + 1
	user := &User{id, register.Name, register.Email, register.PassWord}
	userMap[id] = user
	return user
}
