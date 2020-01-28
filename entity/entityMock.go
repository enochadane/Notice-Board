package entity

// CompanyMock ...
var CompanyMock = Company{
	
	Name:     "moti importer and exporter",
	Email:    "import@gmail.com",
	Password: "info lehulum",
	Posts:    []Post{},
	RoleID:   1,
}

//Rol ...
var Rol = Role{
	ID:        1,
	Name:      "henok",
	Users:     []User{},
	Companies: []Company{},
}

// CampSessionMock ...
var CampSessionMock = CompanySession{

	ID:         1,
	UUID:       "_session_one",
	Expires:    0,
	SigningKey: []byte("king"),
}

// PostMock ...
var PostMock = Post{

	Title:       "company post",
	Description: "we provide  quality information",
	Image:       "C/Users/chala/Downloads/Telegram Desktop/photo_2020-01-19_17-12-06.jpg",
	Category:    "category 02",
	CompanyID:   CompanyMock.ID,
	Owner:       CompanyMock.Name,
}

// UserMock ...
var UserMock = User{
	Name:     "NANI",
	Email:    "nani@gmail.com",
	Password: "123455",
}

// UserSessionMock ...
var UserSessionMock = UserSession{

	ID:         3,
	UUID:       "all at the same time",
	Expires:    0,
	SigningKey: []byte("moti"),
}

// AplMock ...
var AplMock = Application{

	FullName: "kani koni kana",
	Email:    "example.com",
	Phone:    "098765",
	Letter:   "first letter",
	Resume:   "resume again",
	PostID:   2,
	UserID:   2,
}

// ReqMock ...
var ReqMock = Request{

	FullName: "NANI john",
	Email:    "nani@gmail.com",
	Phone:    "098765",
	PostID:   2,
	UserID:   2,
}
