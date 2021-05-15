package app

import (
	"encoding/json"
	"strconv"
	"studymanager/model"
	"time"
	_ "time"

	"github.com/antage/eventsource"
)

var MsgCh chan model.Message //메세지를 넣는 채널
var Es eventsource.EventSource

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

/*
func MakeEs() {
	Es = eventsource.New(nil, nil)
	defer Es.Close()

	MsgCh = make(chan model.Message)
}*/
