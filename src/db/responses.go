package db

type DatabaseResponse struct {
	Response   string `json:"response"`
	Proccessed bool   `json:"proccessed"`
}

type GenToken struct {
	Token string `json:"token"`
	DatabaseResponse
}
