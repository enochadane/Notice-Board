package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/amthesonofGod/Notice-Board/User"
	"github.com/amthesonofGod/Notice-Board/entity"
	uuid "github.com/satori/go.uuid"

	"github.com/amthesonofGod/Notice-Board/session"

	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/rtoken"
)

// UserHandler handles user requests
type UserHandler struct {
	tmpl           *template.Template
	userSrv        User.UserService
	postSrv        post.PostService
	sessionService User.SessionService
	userSess       *entity.UserSession
	loggedInUser   *entity.User
	//csrfSignKey    []byte
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")

// NewUserHandler initializes and returns new NewUserHandler
func NewUserHandler(T *template.Template, US User.UserService, PS post.PostService, sessServ User.SessionService, usrSess *entity.UserSession) *UserHandler {
	return &UserHandler{tmpl: T, userSrv: US, postSrv: PS, sessionService: sessServ, userSess: usrSess}
}

// Authenticated checks if a user is authenticated to access a given route
func (uh *UserHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := uh.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, uh.userSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Index handle requests on /
//Index ...
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)

}

// Login handle requests on /login

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		email := r.FormValue("useremail")
		password := r.FormValue("userpassword")

		users, _ := uh.userSrv.Users()

		for _, user := range users {
			if email == user.Email && password == user.Password {
				fmt.Println("authentication successfull! ")

				// if err == http.ErrNoCookie {
				// 	sID, _ := uuid.NewV4()
				// 	cookie = &http.Cookie{
				// 		Name:  "session",
				// 		Value: sID.String(),
				// 		Path:  "/",
				// 	}
				// }

				usr := &user
				uh.loggedInUser = usr
				claims := rtoken.Claims(usr.Email, uh.userSess.Expires)
				session.Create(claims, uh.userSess.UUID, uh.userSess.SigningKey, w)
				newSess, errs := uh.sessionService.StoreSession(uh.userSess)

				if len(errs) > 0 {
					panic(errs)

				}

				uh.userSess = newSess
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				break

			} else {
				fmt.Println("No such user!")
			}
		}

		//io.WriteString(w, cookie.String())

	} else {
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)
	}
}

func (uh *UserHandler) loggedIn(r *http.Request) bool {
	if uh.userSess == nil {
		return false
	}
	userSess := uh.userSess
	c, err := r.Cookie(userSess.UUID)
	if err != nil {
		return false
	}
	ok, err := session.Valid(c.Value, userSess.SigningKey)
	if !ok || (err != nil) {
		return false
	}
	return true
}

// CreateAccount handle requests on /signup-account
func (uh *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
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

		_, errs := uh.userSrv.StoreUser(usr)

		if len(errs) > 0 {
			panic(errs)
		}

		if err == http.ErrNoCookie {
			sID, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session",
				Value: sID.String(),
				Path:  "/",
			}
		}

		session := &entity.UserSession{}
		session.UUID = cookie.Value
		session.UserID = usr.ID

		_, errs = uh.userSrv.StoreSession(session)

		if len(errs) > 0 {
			panic(errs)
		}

		fmt.Println(usr)

		fmt.Println("User added to db")

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	} else {
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)
	}

}

// Home handle requests on /home
func (uh *UserHandler) Home(w http.ResponseWriter, r *http.Request) {

	//get cookie
	_, err := r.Cookie("session")
	if err != nil {
		fmt.Println("no cookie")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	posts, _ := uh.postSrv.Posts()

	uh.tmpl.ExecuteTemplate(w, "home.layout", posts)
}

// Logout Logs the user out
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// get cookie
	cookie, err := r.Cookie("session")

	if err != http.ErrNoCookie {
		_, errs := uh.userSrv.DeleteSession(cookie.Value)
		// session.DeleteSession
		if len(errs) > 0 {
			panic(errs)
		}
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}
