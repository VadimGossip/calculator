package domain

import "time"

type User struct {
	Id           int64     `json:"id"`
	Login        string    `json:"login" binding:"required,login"`
	Password     string    `json:"password" binding:"required,password"`
	Admin        bool      `json:"admin"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Credentials struct {
	Login    string `json:"login" binding:"required,login"`
	Password string `json:"password" binding:"required,password"`
}
