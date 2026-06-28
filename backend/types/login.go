package types

// LoginRequest represents the payload from the frontend after Firebase authentication
type LoginRequest struct {
	IDToken        string `json:"id_token"`
	InvitationCode string `json:"invitation_code,omitempty"`
}

type LoginResponse struct {
	Message string `json:"message"`
	UID     string `json:"uid"`
	Email   string `json:"email"`
	Name    string `json:"display_name"`
}
