package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/chiwon99881/gone-chat/auth"
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
)

func createUser(rw http.ResponseWriter, r *http.Request) {
	requestCreateUserPayload := &requestCreateUserPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestCreateUserPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	dbOperator.CreateUser(requestCreateUserPayload.Username, requestCreateUserPayload.Password, requestCreateUserPayload.Alias)
}

func updateUserAlias(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing user_id in params"))
		return
	}
	requestUpdateUserAliasPayload := &requestUpdateUserAliasPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestUpdateUserAliasPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	currentUser := r.Header.Get("currentUser")
	if userID != currentUser {
		unauthorizedResponse(rw)
		return
	}
	userIDAsUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		utils.HandleError(err)
	}
	updatedUser, err := dbOperator.UpdateUserAlias(uint(userIDAsUint), requestUpdateUserAliasPayload.Alias)
	if err != nil {
		json.NewEncoder(rw).Encode(responseCommonPayload{Message: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseUpdateUserAliasPayload{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Alias:    updatedUser.Alias})
}

func updateUserPassword(rw http.ResponseWriter, r *http.Request) {
	requestUpdateUserPasswordPayload := &requestUpdateUserPasswordPayload{}
	params := mux.Vars(r)
	userID, ok := params["userID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing user_id in params"))
		return
	}
	currentUser := r.Header.Get("currentUser")
	if userID != currentUser {
		unauthorizedResponse(rw)
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestUpdateUserPasswordPayload)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	currentPasswordAsBytes := utils.ToBytes(requestUpdateUserPasswordPayload.CurrentPassword)
	currentPasswordAsHash := utils.ToHexStringHash(currentPasswordAsBytes)
	idAsBytes, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		utils.HandleError(err)
	}
	check := dbOperator.CheckUserPassword(uint(idAsBytes), currentPasswordAsHash)
	if !check {
		badRequestResponse(rw, errors.New("current password is not correct"))
		return
	}
	newPwAsBytes := utils.ToBytes(requestUpdateUserPasswordPayload.NewPassword)
	newPwAsHash := utils.ToHexStringHash(newPwAsBytes)
	err = dbOperator.UpdatePassword(uint(idAsBytes), newPwAsHash)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseCommonPayload{Message: "password changed successfully"})
}

func login(rw http.ResponseWriter, r *http.Request) {
	requestLoginPayload := &requestLoginPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestLoginPayload)
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

func logout(rw http.ResponseWriter, r *http.Request) {
	au, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		unauthorizedResponse(rw)
		return
	}
	deleted, err := auth.DeleteAuth(au.AccessUUID)
	if deleted == 0 || err != nil {
		unauthorizedResponse(rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseCommonPayload{Message: "Successfully logged out"})
}
