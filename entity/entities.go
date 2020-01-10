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

<<<<<<< HEAD
// CompanySession represents company sessions
type CompanySession struct {
	gorm.Model
	UUID		string
	CompanyID	uint
=======
// Session represents sessions
type Session struct {
	gorm.Model
	UUID		string
	CompanyID	uint
	UserID		uint
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
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

<<<<<<< HEAD
// UserSession represents user sessions
type UserSession struct {
	gorm.Model
	UUID		string
	UserID		uint
}

=======
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
// Application represents job applications forwarded by application users
type Application struct {
	gorm.Model
	FullName 	string
	Email		string
	Phone		string
	Letter		string
	Resume		string
	PostID		uint
<<<<<<< HEAD
	UserID		uint
=======
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
}

// Request represents event join requests forwarded by application users
type Request struct {
	gorm.Model
	FullName	string
	Email		string
	Phone		string
	PostID		uint
<<<<<<< HEAD
	UserID		uint
=======
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
}
