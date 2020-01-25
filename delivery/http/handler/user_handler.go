package handler

import (
	"time"
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/amthesonofGod/Notice-Board/user"
	"github.com/amthesonofGod/Notice-Board/entity"
	
	// uuid "github.com/satori/go.uuid"

	"github.com/amthesonofGod/Notice-Board/session"

	"github.com/amthesonofGod/Notice-Board/post"
	
	"github.com/amthesonofGod/Notice-Board/rtoken"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user requests
type UserHandler struct {
	tmpl           *template.Template
	userSrv        user.UserService
	postSrv        post.PostService
	sessionService user.SessionService
	userSess       *entity.UserSession
	loggedInUser   *entity.User
	//csrfSignKey    []byte
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")

// NewUserHandler initializes and returns new NewUserHandler
func NewUserHandler(T *template.Template, US user.UserService, PS post.PostService, sessServ user.SessionService, usrSess *entity.UserSession) *UserHandler {
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
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)

}

// func CheckPasswordHash(password, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }

// Login handle requests on /login
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	cookie, errc := r.Cookie("session")

	expireToken := time.Now().Add(time.Minute*30).Unix()
	// expireCookie := time.Now().Add(time.Minute*30)


	if r.Method == http.MethodPost {

		email := r.FormValue("useremail")
		password := r.FormValue("userpassword")
		users, _ := uh.userSrv.Users()

		
		for _, user := range users {
			if email == user.Email {
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == bcrypt.ErrMismatchedHashAndPassword {
					fmt.Println("Your email address or password is wrong")
					return
				}

				// match := CheckPasswordHash(password, user.Password)
				// fmt.Println("Match:   ", match)


				claims := rtoken.Claims(email, expireToken)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(email))

				if errc == http.ErrNoCookie {
					// sID, _ := uuid.NewV4()
					cookie = &http.Cookie{
						Name:  "session",
						Value: signedToken,
						Path:  "/",
					}
				}

				session := &entity.UserSession{}
				session.UUID = cookie.Value
				session.UserID = user.ID

				_, errs := uh.userSrv.StoreSession(session)

				if len(errs) > 0 {
					panic(errs)
				}

				fmt.Println(user.Password)
				fmt.Println(password)

				fmt.Println("authentication successfull!")

				http.SetCookie(w, cookie)
				fmt.Println(cookie.Value)
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				break
			} else {
				fmt.Println("user not found")
				// http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}

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

	cookie, errc := r.Cookie("session")
	
	expireToken := time.Now().Add(time.Minute*30).Unix()
	// expireCookie := time.Now().Add(time.Minute*30)

	if r.Method == http.MethodPost {

		usr := &entity.User{}
		usr.Name = r.FormValue("username")
		usr.Email = r.FormValue("useremail")
		password := r.FormValue("userpassword")
		// confirmpass := r.FormValue("confirmPassword")

		users, _ := uh.userSrv.Users()

		for _, user := range users {

			if usr.Email == user.Email {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				fmt.Println("This Email is already in use! ")
				return
			}
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			// singnUpForm.VErrors.Add("password", "Password Could not be stored")
			// uh.tmpl.ExecuteTemplate(w, "signup.layout", singnUpForm)
			panic(err)
		}

		
		usr.Password = string(hashedPassword)

		fmt.Println(usr.Password)

		_, errs := uh.userSrv.StoreUser(usr)

		if len(errs) > 0 {
			panic(errs)
		}

		claims := rtoken.Claims(usr.Email, expireToken)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(usr.Email))

		if errc == http.ErrNoCookie {
			// sID, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session",
				Value: signedToken,
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

// Logout hanldes the POST /logout requests
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	// session.Remove(userSess.UUID, w)
	// uh.sessionService.DeleteSession(userSess.UUID)

	cookie, _ := r.Cookie("session")

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
		Path:  "/",

	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
