package ws

import (
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(ContentTypeMiddleware)
	fmt.Println("Websocket server listening on localhost:4040")

	router.HandleFunc("/ws/{roomID:[0-9]+}", chat.UpgradeWithRoom).Methods("GET")
	utils.HandleError(http.ListenAndServe(":4040", router))
}
