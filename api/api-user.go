package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	if requestCreateUserPayload.Password != requestCreateUserPayload.ConfirmPassword {
		badRequestResponse(rw, errors.New("password is not matched"))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	dbOperator.CreateUser(requestCreateUserPayload.Username, requestCreateUserPayload.Password, requestCreateUserPayload.Alias)
}

func getMe(rw http.ResponseWriter, r *http.Request) {
	currentUser := r.Header.Get("currentUser")
	userIDAsUint := utils.ToUintFromString(currentUser)
	me, err := dbOperator.GetUser(userIDAsUint)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseGetMePayload{
		ID:        me.ID,
		Username:  me.Username,
		Alias:     me.Alias,
		Avatar:    me.Avatar,
		CreatedAt: me.CreatedAt,
		UpdatedAt: me.UpdatedAt,
	})
}

func deleteUser(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing user_id in parameter"))
		return
	}
	currentUser := r.Header.Get("currentUser")
	if currentUser != userID {
		unauthorizedResponse(rw)
		return
	}
	userIDAsUint := utils.ToUintFromString(userID)
	err := dbOperator.DeleteUser(userIDAsUint)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseCommonPayload{Message: "user deleted successfully"})
}

func updateUserAlias(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userID"]
	if !ok {
		badRequestResponse(rw, errors.New("missing user_id in parameter"))
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
	userIDAsUint := utils.ToUintFromString(userID)
	updatedUser, err := dbOperator.UpdateUserAlias(userIDAsUint, requestUpdateUserAliasPayload.Alias)
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
		badRequestResponse(rw, errors.New("missing user_id in parameter"))
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
	userIDAsUint := utils.ToUintFromString(userID)
	check := dbOperator.CheckUserPassword(userIDAsUint, currentPasswordAsHash)
	if !check {
		badRequestResponse(rw, errors.New("current password is not correct"))
		return
	}
	newPwAsBytes := utils.ToBytes(requestUpdateUserPasswordPayload.NewPassword)
	newPwAsHash := utils.ToHexStringHash(newPwAsBytes)
	err = dbOperator.UpdatePassword(userIDAsUint, newPwAsHash)
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
		Avatar:    me.Avatar,
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

func uploadImage(rw http.ResponseWriter, r *http.Request) {
	sess := utils.ConnectAws()
	uploader := s3manager.NewUploader(sess)

	awsBucket := os.Getenv("AWS_BUCKET_NAME")
	file, header, err := r.FormFile("photo")
	if err != nil {
		utils.HandleError(err)
	}
	filename := header.Filename

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(responseUploadImage{
			ErrMessage: "failed to upload file",
			Uploader:   up,
		})
		return
	}
	filepath := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", awsBucket, os.Getenv("AWS_REGION"), filename)

	currentUser := r.Header.Get("currentUser")
	userID := utils.ToUintFromString(currentUser)
	err = dbOperator.UpdateUserAvatar(userID, filepath)
	if err != nil {
		badRequestResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseUploadImage{
		FilePath: filepath,
	})
}
