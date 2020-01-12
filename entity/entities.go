package entity

import (
	"github.com/jinzhu/gorm"
)

// Company represents Companies
type Company struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Posts    []Post
}

// CompanySession represents company sessions
type CompanySession struct {
	ID         uint
	CompanyID  uint
	UUID       string `gorm:"type:varchar(255);not null"`
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

// Session represents sessions
// type Session struct {
// 	ID         uint
// 	UUID       string `gorm:"type:varchar(255);not null"`
// 	Expires    int64  `gorm:"type:varchar(255);not null"`
// 	SigningKey []byte `gorm:"type:varchar(255);not null"`
// }

// Post represents Posts
type Post struct {
	gorm.Model
	Title       string `gorm:"type:varchar(255);not null"`
	Description string
	Image       string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(255);not null"`
	CompanyID   uint
	Owner       string `gorm:"type:varchar(255);not null"`
}

// PostSession represents post sessions
type PostSession struct {
	gorm.Model
	UUID	string
	PostID	uint
}

// User represents Users
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

// UserSession represents user sessions
type UserSession struct {
	ID         uint
	UserID     uint
	UUID       string `gorm:"type:varchar(255);not null"`
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

// Application represents job applications forwarded by application users
type Application struct {
	gorm.Model
	FullName string
	Email    string
	Phone    string
	Letter   string
	Resume   string
	PostID   uint
	UserID   uint
}

// Request represents event join requests forwarded by application users
type Request struct {
	gorm.Model
	FullName string
	Email    string
	Phone    string
	PostID   uint
	UserID   uint
}
