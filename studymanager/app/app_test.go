package app

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	_ "strings"
	"studymanager/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
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
