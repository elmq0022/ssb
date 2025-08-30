package schemas

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonToken struct {
	Token string `json:"token"`
}
