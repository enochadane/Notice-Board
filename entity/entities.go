package entity

import (
	"github.com/jinzhu/gorm"
)

// Company represents Companies
type Company struct {
	gorm.Model
	Name 		string `gorm:"type:varchar(255);not null"`
	Email 		string `gorm:"type:varchar(255);not null"`
	Password 	string `gorm:"type:varchar(255);not null"`
	Posts		[]Post
}

// CompanySession represents company sessions
type CompanySession struct {
	gorm.Model
	UUID		string
	CompanyID	uint
}

// Post represents Posts
type Post struct {
	gorm.Model
	Title		string `gorm:"type:varchar(255);not null"`
	Description string
	Image		string `gorm:"type:varchar(255)"`
	Category	string `gorm:"type:varchar(255);not null"`
	CompanyID	uint
	Owner		string `gorm:"type:varchar(255);not null"`
}

// User represents Users
type User struct {
	gorm.Model
	Name 		string `gorm:"type:varchar(255);not null"`
	Email 		string `gorm:"type:varchar(255);not null"`
	Password 	string `gorm:"type:varchar(255);not null"`
}

// UserSession represents user sessions
type UserSession struct {
	gorm.Model
	UUID		string
	UserID		uint
}

// Application represents job applications forwarded by application users
type Application struct {
	gorm.Model
	FullName 	string
	Email		string
	Phone		string
	Letter		string
	Resume		string
	PostID		uint
	UserID		uint
}

// Request represents event join requests forwarded by application users
type Request struct {
	gorm.Model
	FullName	string
	Email		string
	Phone		string
	PostID		uint
	UserID		uint
}
