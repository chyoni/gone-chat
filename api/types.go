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

type responseLoginPayload struct {
	ID        uint              `json:"id"`
	Username  string            `json:"username"`
	Alias     string            `json:"alias"`
	CreatedAt int               `json:"created_at"`
	UpdatedAt int               `json:"updated_at"`
	Tokens    map[string]string `json:"token"`
}

type responseError struct {
	ErrMessage string `json:"error_message"`
}
