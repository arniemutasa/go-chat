package main

import (
	"log"
	"net/http"

	"github.com/arniemutasa/go-chat/internal/handlers"
)

func main() {

	mux := routes()

	log.Println("Starting Channel Listener")

	go handlers.ListenToWebsocketChannel()

	log.Println("Starting Server at Port 8080")

	_ = http.ListenAndServe(":8080", mux)
}
