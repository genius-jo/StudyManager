package main

import (
	"log"
	"net/http"
	"studymanager/app"

	_ "github.com/urfave/negroni"
)

func main() {

	m := app.MakeHandler("./test.db")
	defer m.Close()

	log.Println("start")
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
