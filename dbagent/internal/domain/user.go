package domain

import "time"

type User struct {
	Id           int64
	Login        string
	Password     string
	Admin        bool
	RegisteredAt time.Time
}

type Credentials struct {
	Login    string `json:"login" binding:"required,login"`
	Password string `json:"password" binding:"required,password"`
}
