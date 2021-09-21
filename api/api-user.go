package api

import (
	"encoding/json"
	"net/http"

	"github.com/chiwon99881/gone-chat/auth"
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
	if user.Username != requestLoginPayload.Username {
		unauthorizedResponse(rw)
		return
	}
	pwAsBytes := utils.ToBytes(requestLoginPayload.Password)
	pwAsHash := utils.ToHexStringHash(pwAsBytes)
	if user.Password != pwAsHash {
		unauthorizedResponse(rw)
		return
	}
	td, err := auth.CreateToken(user.ID)
	if err != nil {
		unprocessableEntityResponse(rw, err)
		return
	}
	err = auth.CreateAuth(user.ID, td)
	if err != nil {
		unprocessableEntityResponse(rw, err)
		return
	}
	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}
	me, _ := dbOperator.GetUser(user.ID)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseLoginPayload{
		ID:        me.ID,
		Username:  me.Username,
		Alias:     me.Alias,
		CreatedAt: me.CreatedAt,
		UpdatedAt: me.UpdatedAt,
		Tokens:    tokens,
	})
}
