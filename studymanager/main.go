package main

import (
	"log"
	"net/http"
	"studymanager/app"
	"studymanager/model"

	"github.com/antage/eventsource"
)

func main() {

	//이벤트 소스
	app.Es = eventsource.New(nil, nil)
	defer app.Es.Close()

	//메세지 채널
	app.MsgCh = make(chan model.Message)

	go app.ProcessMsgCh(app.Es)

	m := app.MakeHandler("./test.db")
	defer m.Close()

	log.Println("StudyManager start")
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
