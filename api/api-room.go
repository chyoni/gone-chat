package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chiwon99881/gone-chat/chat"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
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
	dbOperator.CreateChatRecord(
		utils.ToUintFromString(requestMessagePayload.RoomID),
		utils.ToUintFromString(userID),
		requestMessagePayload.Message,
	)
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
	var participants []uint
	for _, value := range room.Participants {
		participants = append(participants, value.ID)
	}
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(responseCreateRoomPayload{
		ID:           room.ID,
		CreatedAt:    room.CreatedAt,
		UpdatedAt:    room.UpdatedAt,
		Participants: participants,
	})
	// bearerToken := r.Header.Get("Authorization")
	// upgradeURL := fmt.Sprintf("http://127.0.0.1:4000/ws/%d", room.ID)

	// client := &http.Client{}
	// req, _ := http.NewRequest("GET", upgradeURL, nil)
	// req.Header.Set("Authorization", bearerToken)
	// resp, err := client.Do(req)

	// if err != nil || resp.StatusCode != 200 {
	// 	rw.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(rw).Encode(responseError{ErrMessage: err.Error()})
	// 	return
	// }
}

func getRooms(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing user id in parameter"))
		return
	}
	currentUser := r.Header.Get("currentUser")
	if userID != currentUser {
		unauthorizedResponse(rw)
		return
	}
	userIDByUint := utils.ToUintFromString(userID)
	userRooms, err := dbOperator.GetRoomsByUserID(userIDByUint)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	var myRooms []uint
	for _, values := range userRooms {
		myRooms = append(myRooms, values.RoomID)
	}
	responseGetRoomPayload := &responseGetRoomPayload{RoomID: myRooms}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseGetRoomPayload)
}

func getsUsersByRoom(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID, ok := params["roomID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing room_id in parameter"))
		return
	}
	roomIDByUint := utils.ToUintFromString(roomID)
	usersForRoom, err := dbOperator.GetUsersByRoomID(roomIDByUint)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseGetUsersByRoom{
		RoomID: usersForRoom.RoomID,
		Users:  usersForRoom.Users,
	})
}
