package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/janczizikow/battleship-backend/rooms"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	roomHandler := rooms.NewHandler(logger)
	r := mux.NewRouter()
	r.HandleFunc("/rooms", roomHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/rooms/{roomCode}/join", roomHandler.Join).Methods(http.MethodPost)

	serv := http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(serv.ListenAndServe())
}
