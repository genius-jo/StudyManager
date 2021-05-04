package main

import (
	"log"
	"net/http"
	"studymanager/app"

	"github.com/urfave/negroni"
)

func main() {

	m := app.MakeHandler("./test.db")
	defer m.Close()

	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("start")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
