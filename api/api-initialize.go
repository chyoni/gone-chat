package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/chiwon99881/gone-chat/auth"
	"github.com/chiwon99881/gone-chat/database"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dbOperator database.Repository = database.RepoOperator{}

func tokenInvalidResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: "token invalid", TokenRefreshFlag: true})
}

func badRequestResponse(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
}

func unauthorizedResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: "you are not authorized"})
}

func unprocessableEntityResponse(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" && r.URL.Path != "/user" {
			tokenAuth, err := auth.ExtractTokenMetadata(r)
			if err != nil {
				tokenInvalidResponse(w)
				return
			}
			userID, err := auth.FetchAuth(tokenAuth)
			if err != nil {
				tokenInvalidResponse(w)
				return
			}
			r.Header.Add("currentUser", strconv.Itoa(int(userID)))
		}
		next.ServeHTTP(w, r)
	})
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(ContentTypeMiddleware)
	router.Use(AuthMiddleware)
	fmt.Println("Http server listening on localhost:4000")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowedHeaders: []string{os.Getenv("CORS_ALLOWED_HEADERS")},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
	})

	handler := c.Handler(router)

	router.HandleFunc("/message", message).Methods("POST")
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/user/info/me", getMe).Methods("GET")
	router.HandleFunc("/user/{userID:[0-9]+}", deleteUser).Methods("DELETE")
	router.HandleFunc("/user/alias/{userID:[0-9]+}", updateUserAlias).Methods("POST")
	router.HandleFunc("/user/password/{userID:[0-9]+}", updateUserPassword).Methods("POST")
	router.HandleFunc("/user/room/{userID:[0-9]+}", getRooms).Methods("GET")
	router.HandleFunc("/room", createRoom).Methods("POST")
	router.HandleFunc("/room/{roomID:[0-9]+}/users", getsUsersByRoom).Methods("GET")
	router.HandleFunc("/room/{roomID:[0-9]+}/messages", getAllMessagesByRoom).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/image/upload", uploadImage).Methods("POST")
	utils.HandleError(http.ListenAndServe(":4000", handler))
}
