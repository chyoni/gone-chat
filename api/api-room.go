package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/chat"
)

func message(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("currentUser")
	requestMessagePayload := &requestMessagePayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestMessagePayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	chat.HandleMessage(requestMessagePayload.Message, requestMessagePayload.RoomID, userID)
}

func createRoom(rw http.ResponseWriter, r *http.Request) {
	requestCreateRoomPayload := &requestCreateRoomPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestCreateRoomPayload)
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
	bearerToken := r.Header.Get("Authorization")
	upgradeURL := fmt.Sprintf("http://127.0.0.1:4000/ws/%d", room.ID)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", upgradeURL, nil)
	req.Header.Set("Authorization", bearerToken)
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(responseError{ErrMessage: "upgrade fail cause: should pass token"})
		return
	}
}
