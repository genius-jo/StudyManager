package main

import (
	_ "encoding/json"
	"log"
	"net/http"
	_ "strconv"
	"studymanager/app"
	"studymanager/model"
	_ "studymanager/model"
	_ "time"

	"github.com/antage/eventsource"
	_ "github.com/antage/eventsource"

	_ "github.com/urfave/negroni"
)

/*
var msgCh chan model.Message //메세지를 넣는 채널
var es eventsource.EventSource

//모든 클라이언트에게 메세지를 보냄
func sendMessage(name, msg string) {

	//다른 스레드에 큐 형태로 보냄
	msgCh <- model.Message{name, msg}
}

//채널에서 pop해서 이벤트 소스로 보냄
func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond())) //유니크한 ID는 현재시간으로
	}
}
*/
func main() {

	//es = eventsource.New(nil, nil)
	//defer es.Close()

	//msgCh = make(chan model.Message)

	//app.MakeEs()
	app.Es = eventsource.New(nil, nil)
	defer app.Es.Close()

	app.MsgCh = make(chan model.Message)

	go app.ProcessMsgCh(app.Es)

	m := app.MakeHandler("./test.db")
	defer m.Close()

	log.Println("start")
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
