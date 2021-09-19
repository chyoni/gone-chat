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
		fmt.Fprintf(rw, "%s", responseError{errMessage: err.Error()})
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
		fmt.Fprintf(rw, "%s", responseError{errMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	dbOperator.CreateUser(requestCreateUserPayload.Username, requestCreateUserPayload.Password)
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
	utils.HandleError(http.ListenAndServe(":4000", handler))
}
