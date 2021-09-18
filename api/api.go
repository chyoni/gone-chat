package api

import (
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("Server listening on localhost:4000")

	router.HandleFunc("/ws", chat.Upgrade).Methods("GET")
	utils.HandleError(http.ListenAndServe(":4000", router))
}
