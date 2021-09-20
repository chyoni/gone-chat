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

func badRequestResponse(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
}

func unauthorizedResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: "you are not authorized"})
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
	router.HandleFunc("/login", login).Methods("POST")
	utils.HandleError(http.ListenAndServe(":4000", handler))
}
