package explorer

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
)

func room(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/chat.htm")
	utils.HandleError(err)
	err = tmpl.Execute(rw, nil)
	utils.HandleError(err)
}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/room", room).Methods("GET")

	fmt.Println("Explorer listening on localhost:5000")
	err := http.ListenAndServe(":5000", router)
	utils.HandleError(err)
}
