package api

type requestMessagePayload struct {
	Message string `json:"message"`
	RoomID  string `json:"roomId"`
}

type requestCreateUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type requestCreateRoomPayload struct {
	Participants []uint `json:"participants"`
}

type requestLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseError struct {
	ErrMessage string `json:"error_message"`
}
