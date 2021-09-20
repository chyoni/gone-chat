package api

import (
	"encoding/json"
	"net/http"

	"github.com/chiwon99881/gone-chat/utils"
)

func createUser(rw http.ResponseWriter, r *http.Request) {
	requestCreateUserPayload := &requestCreateUserPayload{}
	err := json.NewDecoder(r.Body).Decode(requestCreateUserPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	dbOperator.CreateUser(requestCreateUserPayload.Username, requestCreateUserPayload.Password, requestCreateUserPayload.Alias)
}

func login(rw http.ResponseWriter, r *http.Request) {
	requestLoginPayload := &requestLoginPayload{}
	err := json.NewDecoder(r.Body).Decode(requestLoginPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	user, err := dbOperator.FindUserByUsername(requestLoginPayload.Username)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	pwAsBytes := utils.ToBytes(requestLoginPayload.Password)
	pwAsHash := utils.ToHexStringHash(pwAsBytes)
	if user.Password != pwAsHash {
		unauthorizedResponse(rw)
		return
	}
	// generate token
}
