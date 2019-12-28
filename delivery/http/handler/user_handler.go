package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type UserHandler struct {
	tmpl	*template.Template
	userSrv model.UserService
	postSrv model.PostService
}

func NewUserHandler(T *template.Template, US model.UserService, PS model.PostService) *UserHandler {
	return &UserHandler{tmpl: T, userSrv: US, postSrv: PS}
}

// func (uh *UserHandler) Signin(w http.ResponseWriter, r *http.Request) {
// 	uh.tmpl.ExecuteTemplate(w, "signin.layout", nil)
// }

// func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
// 	uh.tmpl.ExecuteTemplate(w, "signup.layout", nil)
// }

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		email := r.FormValue("useremail")
		password := r.FormValue("userpassword")

		users, _ := uh.userSrv.Users()
		
		for _, user := range users {
			fmt.Println(users)
			if email == user.Email && password == user.Password {
				fmt.Println("authentication successfull! ")
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				break
			
			} else {
				fmt.Println("No such user!")
			}
		}
	} else {
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)
	}
}

func (uh *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == http.MethodPost {
		
		usr := entity.User{}
		usr.Name = r.FormValue("username")
		usr.Email = r.FormValue("useremail")
		usr.Password = r.FormValue("userpassword")
		// confirmpass := r.FormValue("confirmPassword")

		users, _ := uh.userSrv.Users()

		for _, user := range users {
			
			if usr.Email == user.Email {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				fmt.Println("This Email is already in use! ")
				return
			}
		}

		// if usr.Password == confirmpass {

			err := uh.userSrv.StoreUser(usr)

			if err != nil {
				panic(err)
			}

			fmt.Println(users)

			fmt.Println(usr)

			fmt.Println("User added to db")

			http.Redirect(w, r, "/home", http.StatusSeeOther)

		// } else {
		// 	http.Redirect(w, r, "/signup", http.StatusSeeOther)
		// 	fmt.Println("Password doesn't match! ")
		// }
		
	} else {
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)
	}

}

func (uh *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, _ := uh.postSrv.Posts()
	uh.tmpl.ExecuteTemplate(w, "home.layout", posts)
}