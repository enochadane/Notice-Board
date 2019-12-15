package entity

import (
	"time"
)

type User struct {
	Id			int
	Name 		string
	Email 		string
	Password 	string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

