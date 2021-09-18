package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type MessagePayload struct {
	Message string `json:"message"`
	RoomID  string `json:"roomId"`
}

func message(rw http.ResponseWriter, r *http.Request) {
	messagePayload := &MessagePayload{}
	err := json.NewDecoder(r.Body).Decode(messagePayload)
	utils.HandleError(err)

	chat.HandleMessage(messagePayload.Message, messagePayload.RoomID)
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("Server listening on localhost:4000")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("CORS_ALLOWED")},
	})

	handler := c.Handler(router)

	router.HandleFunc("/ws/{roomID:[0-9]+}", chat.UpgradeWithRoom).Methods("GET")
	router.HandleFunc("/message", message).Methods("POST")
	utils.HandleError(http.ListenAndServe(":4000", handler))
}
