package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"studymanager/model"
	"time"

	"github.com/antage/eventsource"
)

var MsgCh chan model.Message   //메세지를 넣는 채널
var Es eventsource.EventSource //이벤트 소스

//모든 클라이언트에게 메세지를 보냄
func sendMessage(name, msg string) {

	//다른 스레드에 큐 형태로 보냄
	MsgCh <- model.Message{name, msg}
}

//채널에서 pop해서 이벤트 소스로 보냄
func ProcessMsgCh(es eventsource.EventSource) {
	for msg := range MsgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond())) //유니크한 ID는 현재시간으로
	}
}

func (ah *AppHandler) indexChatHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/chat.html", http.StatusTemporaryRedirect)
}

func (ah *AppHandler) messageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}

func (ah *AppHandler) addNameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")

	//해당 유저가 들어왔다는거을 알림 (모든 유저들에게 메세지를 보냄)
	sendMessage("", fmt.Sprintf("add user: %s", username))
}

func (ah *AppHandler) leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	//해당 유저가 나갔다는것을 알림
	sendMessage("", fmt.Sprintf("left user: %s", username))
}
