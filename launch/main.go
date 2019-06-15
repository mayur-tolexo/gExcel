package main

import (
	"log"
	"net/http"

	"github.com/mayur-tolexo/gExcel"
)

func main() {
	hub := gExcel.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
