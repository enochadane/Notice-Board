package entity

import (
	"github.com/jinzhu/gorm"   /*To install GORM just use the following command = go  get “github.com/jinzhu/gorm” */
)

// Company represents Companies
type Company struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null; unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	Posts    []Post
	RoleID   uint
}

// CompanySession represents company sessions
type CompanySession struct {
	ID         uint
	UUID       string `gorm:"type:varchar(255);not null"`
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

// Post represents Posts
type Post struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description"`
	Image       string `gorm:"type:varchar(255)"`
	Category    string `json:"category" gorm:"type:varchar(255);not null"`
	CompanyID   uint
	Owner       string `json:"owner" gorm:"type:varchar(255);not null"`
}

// User represents Users
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null; unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	RoleID	 uint
}

// Role repesents application user roles
type Role struct {
	ID    		uint
	Name  		string `gorm:"type:varchar(255)"`
	Users 		[]User
	Companies 	[]Company
}

// UserSession represents user sessions
type UserSession struct {
	ID         uint
	UUID       string `gorm:"type:varchar(255);not null"`
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

// Application represents job applications forwarded by application users
type Application struct {
	gorm.Model
	FullName string `json:"fullname" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null; unique"`
	Phone    string `json:"phone" gorm:"type:varchar(255);not null; unique"`
	Letter   string `json:"letter"`
	Resume   string `json:"resume"`
	PostID   uint
	UserID   uint
}

// Request represents event join requests forwarded by application users
type Request struct {
	gorm.Model
	FullName string `json:"fullname" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null; unique"`
	Phone    string `json:"phone" gorm:"type:varchar(255);not null; unique"`
	PostID   uint
	UserID   uint
}

