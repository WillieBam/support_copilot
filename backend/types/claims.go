package types

import jwt "github.com/golang-jwt/jwt/v5"

type Claims struct {
	FirebaseUID string `json:"firebase_uid"`
	Email       string `json:"email"`
	MfaVerified bool   `json:"mfa_verified"`
	jwt.RegisteredClaims
}
