package api

type requestMessagePayload struct {
	Message string `json:"message"`
	RoomID  string `json:"roomId"`
}

type requestCreateUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseError struct {
	errMessage string
}
