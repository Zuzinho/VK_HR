package handler

type authResponse struct {
	Token string `json:"token"`
}

type insertedResponse struct {
	InsertedID int `json:"inserted_id"`
}
