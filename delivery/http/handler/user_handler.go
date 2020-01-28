package handler

import (
	// "time"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/amthesonofGod/Notice-Board/permission"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/user"

	// uuid "github.com/satori/go.uuid"

	"github.com/amthesonofGod/Notice-Board/session"

	"github.com/amthesonofGod/Notice-Board/form"
	"github.com/amthesonofGod/Notice-Board/post"

	"github.com/amthesonofGod/Notice-Board/rtoken"
	// "github.com/dgrijalva/jwt-go"
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
	userRole       user.RoleService
	csrfSignKey    []byte
}

var currentUser *entity.User

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")

// NewUserHandler initializes and returns new NewUserHandler
func NewUserHandler(T *template.Template, US user.UserService, PS post.PostService, sessServ user.SessionService, usrSess *entity.UserSession, csKey []byte) *UserHandler {
	return &UserHandler{tmpl: T, userSrv: US, postSrv: PS, sessionService: sessServ, userSess: usrSess, csrfSignKey: csKey}
}

//CurrentUser ...
func (uh *UserHandler) CurrentUser() *entity.User {
	currentUser := uh.loggedInUser
	return currentUser
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

// Authorized checks if a user has proper authority to access a give route
func (uh *UserHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uh.loggedInUser == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := uh.userSrv.UserRoles(uh.loggedInUser)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
			if !permitted {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		if r.Method == http.MethodPost {
			ok, err := rtoken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
			if !ok || (err != nil) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Index handle requests on /
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)

}

//func CheckPasswordHash(password, hash string) bool {
//    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//    return err == nil
//}

// Login handle requests on /login
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		loginForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", loginForm)
		return
	}

	if r.Method == http.MethodPost {

		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		usr, errs := uh.userSrv.UserByEmail(r.FormValue("useremail"))
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", loginForm)
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(r.FormValue("userpassword")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", loginForm)
			return
		}

		//match := CheckPasswordHash(r.FormValue("userpassword"), usr.Password)
		//fmt.Println("Match:   ", match)

		uh.loggedInUser = usr
		currentUser = uh.loggedInUser
		claims := rtoken.Claims(usr.Email, uh.userSess.Expires)

		fmt.Println(usr.ID)

		session.Create(claims, uh.userSess.UUID, uh.userSess.SigningKey, w)

		newSess, errs := uh.sessionService.StoreSession(uh.userSess)

		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to store session")
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", loginForm)
			return
		}
		uh.userSess = newSess

		fmt.Println("authentication successfull!")

		http.Redirect(w, r, "/home", http.StatusSeeOther)

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

	token, err := rtoken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if r.Method == http.MethodGet {
		signUpForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", signUpForm)
		return
	}

	if r.Method == http.MethodPost {

		// Parsing the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		// Validate the form contents
		singnUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("username", "useremail", "userpassword", "confirmpassword")
		singnUpForm.MatchesPattern("useremail", form.EmailRX)
		singnUpForm.MinLength("userpassword", 8)
		singnUpForm.PasswordMatches("userpassword", "confirmpassword")
		singnUpForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !singnUpForm.Valid() {
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", singnUpForm)
			return
		}

		eExists := uh.userSrv.EmailExists(r.FormValue("useremail"))
		if eExists {
			singnUpForm.VErrors.Add("email", "Email Already Exists")
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", singnUpForm)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("userpassword")), 12)
		if err != nil {
			singnUpForm.VErrors.Add("password", "Password Could not be stored")
			uh.tmpl.ExecuteTemplate(w, "index_signin_signup.html", singnUpForm)
			return
		}

		//role, errs := uh.userRole.RoleByName("USER")
		//
		//if len(errs) > 0 {
		//	fmt.Println("could not assign role to the user")
		//	singnUpForm.VErrors.Add("role", "could not assign role to the user")
		//	uh.tmpl.ExecuteTemplate(w, "signup.layout", singnUpForm)
		//	return
		//}

		user := &entity.User{
			Name:     r.FormValue("username"),
			Email:    r.FormValue("useremail"),
			Password: string(hashedPassword),
			RoleID:   0,
		}

		_, errs := uh.userSrv.StoreUser(user)

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

// Home handle requests on /home
func (uh *UserHandler) Home(w http.ResponseWriter, r *http.Request) {

	posts, _ := uh.postSrv.Posts()

	uh.tmpl.ExecuteTemplate(w, "home.layout", posts)
}

// Logout hanldes the POST /logout requests
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.UserSession)
	session.Remove(uh.userSess.UUID, w)
	uh.sessionService.DeleteSession(uh.userSess.UUID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
