package api

import "github.com/chiwon99881/gone-chat/entity"

type requestMessagePayload struct {
	Message string `json:"message"`
	RoomID  string `json:"roomId"`
}

type requestCreateUserPayload struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Alias           string `json:"alias"`
}

type requestUpdateUserAliasPayload struct {
	Alias string `json:"alias"`
}

type requestUpdateUserPasswordPayload struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type requestCreateRoomPayload struct {
	Participants []uint `json:"participants"`
}

type requestLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseUpdateUserAliasPayload struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Alias    string `json:"alias"`
}

type responseLoginPayload struct {
	ID        uint              `json:"id"`
	Username  string            `json:"username"`
	Alias     string            `json:"alias"`
	Avatar    string            `json:"avatar"`
	CreatedAt int               `json:"created_at"`
	UpdatedAt int               `json:"updated_at"`
	Tokens    map[string]string `json:"token"`
}

type responseGetMePayload struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Alias     string `json:"alias"`
	Avatar    string `json:"avatar"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type responseGetRoomPayload struct {
	RoomID []uint `json:"room_id"`
}

type responseCreateRoomPayload struct {
	ID           uint   `json:"room_id"`
	CreatedAt    int    `json:"created_at"`
	UpdatedAt    int    `json:"updated_at"`
	Participants []uint `json:"participants"`
}

type responseCommonPayload struct {
	Message string `json:"message"`
}

type responseError struct {
	ErrMessage       string `json:"error_message"`
	TokenRefreshFlag bool   `json:"token_refresh_flag,omitempty"`
}

type responseUploadImage struct {
	ErrMessage string      `json:"error_message,omitempty"`
	FilePath   string      `json:"filepath,omitempty"`
	Uploader   interface{} `json:"uploader,omitempty"`
}

type responseGetUsersByRoom struct {
	RoomID uint           `json:"room_id"`
	Users  []*entity.User `json:"users"`
}

type cleanChat struct {
	RoomID  uint         `json:"roomId"`
	From    *entity.User `json:"from"`
	Message string       `json:"message"`
	Created int          `json:"created"`
}
