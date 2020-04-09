package main

import (
	"command/controller"
	"net/http"
)

func main() {

	cmdHandler := controller.CreateCmdHandler()

	http.Handle("/cmd", cmdHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
