package domain

import "time"

type Token struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiresAt time.Time
}

type TokenResponse struct {
	Token string `json:"accessToken"`
}
