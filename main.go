package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/janczizikow/battleship-backend/rooms"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")
	reader(ws)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	roomHandler := rooms.NewHandler(logger)
	r := mux.NewRouter()
	r.HandleFunc("/rooms", roomHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/rooms/{roomCode}/join", roomHandler.Join).Methods(http.MethodPost)
	r.HandleFunc("/ws", wsEndpoint).Methods(http.MethodGet)

	serv := http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(serv.ListenAndServe())
}
