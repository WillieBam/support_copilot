package requests

type TokenExchangeRequest struct {
	FirebaseToken string `json:"firebase_token"`
}

type TokenExchangeResponse struct {
	Token string `json:"token"`
}
