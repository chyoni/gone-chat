package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/database"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dbOperator database.Repository = database.RepoOperator{}

func message(rw http.ResponseWriter, r *http.Request) {
	requestMessagePayload := &requestMessagePayload{}
	err := json.NewDecoder(r.Body).Decode(requestMessagePayload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	chat.HandleMessage(requestMessagePayload.Message, requestMessagePayload.RoomID)
}

func createUser(rw http.ResponseWriter, r *http.Request) {
	requestCreateUserPayload := &requestCreateUserPayload{}
	err := json.NewDecoder(r.Body).Decode(requestCreateUserPayload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	dbOperator.CreateUser(requestCreateUserPayload.Username, requestCreateUserPayload.Password, requestCreateUserPayload.Alias)
}

func createRoom(rw http.ResponseWriter, r *http.Request) {
	requestCreateRoomPayload := &requestCreateRoomPayload{}
	err := json.NewDecoder(r.Body).Decode(requestCreateRoomPayload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
		return
	}
	_, err = dbOperator.FindUser(requestCreateRoomPayload.Participants[0])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
		return
	}
	room, err := dbOperator.CreateRoom(requestCreateRoomPayload.Participants[0])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
		return
	}
	upgradeURL := fmt.Sprintf("http://127.0.0.1:4000/ws/%d", room.ID)
	_, err = http.Get(upgradeURL)
	if err != nil {
		utils.HandleError(err)
	}
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(ContentTypeMiddleware)
	fmt.Println("Server listening on localhost:4000")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("CORS_ALLOWED")},
	})

	handler := c.Handler(router)

	router.HandleFunc("/ws/{roomID:[0-9]+}", chat.UpgradeWithRoom).Methods("GET")
	router.HandleFunc("/message", message).Methods("POST")
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/room", createRoom).Methods("POST")
	utils.HandleError(http.ListenAndServe(":4000", handler))
}
