package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/utils"
)

func message(rw http.ResponseWriter, r *http.Request) {
	requestMessagePayload := &requestMessagePayload{}
	err := json.NewDecoder(r.Body).Decode(requestMessagePayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	chat.HandleMessage(requestMessagePayload.Message, requestMessagePayload.RoomID)
}

func createRoom(rw http.ResponseWriter, r *http.Request) {
	requestCreateRoomPayload := &requestCreateRoomPayload{}
	err := json.NewDecoder(r.Body).Decode(requestCreateRoomPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	_, err = dbOperator.FindUserByID(requestCreateRoomPayload.Participants[0])
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	room, err := dbOperator.CreateRoom(requestCreateRoomPayload.Participants[0])
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	upgradeURL := fmt.Sprintf("http://127.0.0.1:4000/ws/%d", room.ID)
	_, err = http.Get(upgradeURL)
	if err != nil {
		utils.HandleError(err)
	}
}
