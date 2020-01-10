package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/model"
	"github.com/amthesonofGod/Notice-Board/post"
)

// UserHandler ...
type UserHandler struct {
	tmpl    *template.Template
	userSrv model.UserService
	postSrv post.PostService
}

// NewUserHandler ...
func NewUserHandler(T *template.Template, US model.UserService, PS post.PostService) *UserHandler {
	return &UserHandler{tmpl: T, userSrv: US, postSrv: PS}
}

// Index ...
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)

}

//User ...

// Login ...
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

// CreateAccount ...
func (uh *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		usr := &entity.User{}
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

		_, errs := uh.userSrv.StoreUser(usr)

		if len(errs) > 0 {
			panic(errs)
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

// Home ...
func (uh *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, _ := uh.postSrv.Posts()
	uh.tmpl.ExecuteTemplate(w, "home.layout", posts)
}
