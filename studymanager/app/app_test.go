package app

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	_ "strings"
	"studymanager/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	os.Remove("./test.db") //DB를 지우고 시작
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah) //ah을 인자로 바로 써줄 수 있다
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/users", url.Values{"name": {"jo"}, "email": {"jo@jo.com"}, "pass_word": {"abcd"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user) //json값을 user로 읽음
	assert.NoError(err)
	assert.Equal(user.Name, "jo")
	assert.Equal(user.Email, "jo@jo.com")
	assert.Equal(user.PassWord, "abcd")
	id1 := user.Id

	//하나 더 추가하기
	resp, err = http.PostForm(ts.URL+"/users", url.Values{"name": {"genius"}, "email": {"genius@genius.com"}, "pass_word": {"abcd"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&user) //json값을 user로 읽음
	assert.NoError(err)
	assert.Equal(user.Name, "genius")
	assert.Equal(user.Email, "genius@genius.com")
	assert.Equal(user.PassWord, "abcd")
	id2 := user.Id

	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*model.User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(len(users), 2)
	for _, t := range users {
		if t.Id == id1 {
			assert.Equal(t.Name, "jo")
			assert.Equal(t.Email, "jo@jo.com")
			assert.Equal(t.PassWord, "abcd")
		} else if t.Id == id2 {
			assert.Equal(t.Name, "genius")
			assert.Equal(t.Email, "genius@genius.com")
			assert.Equal(t.PassWord, "abcd")
		} else {
			assert.Error(fmt.Errorf("ID should be id1 or id2"))
		}
	}

}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah) //ah을 인자로 바로 써줄 수 있다
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/login", url.Values{"email": {"jo@jo.com"}, "pass_word": {"abcd"}})
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user) //json값을 user로 읽음
	assert.NoError(err)
	assert.Equal(user.Name, "jo")
	assert.Equal(user.Email, "jo@jo.com")
	assert.Equal(user.PassWord, "abcd")
}

func TestTodos(t *testing.T) {
	getSessionID = func(r *http.Request) string { //test할 세션을 미리 넣음
		return "testsessionId"
	}
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")

	id1 := todo.Id
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2")
	id2 := todo.Id

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, t := range todos {
		if t.Id == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.Id == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testId error"))
		}
	}

	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, t := range todos {
		if t.Id == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil) //delete는 이런식으로 해줘야 한다
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		assert.Equal(t.Id, id2)
	}

}
