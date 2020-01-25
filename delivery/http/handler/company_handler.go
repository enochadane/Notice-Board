package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	
	"golang.org/x/crypto/bcrypt"

	"github.com/amthesonofGod/Notice-Board/rtoken"
	"github.com/dgrijalva/jwt-go"
	"github.com/amthesonofGod/Notice-Board/session"
)

// CompanyHandler handles company handler admin requests
type CompanyHandler struct {
	tmpl       *template.Template
	companySrv company.CompanyService
	postSrv    post.PostService

	sessionService   company.SessionServiceCamp
	campSess         *entity.CompanySession
	loggedInUserCamp *entity.Company
}

// NewCompanyHandler initializes and returns new NewCompanyHandler
func NewCompanyHandler(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession) *CompanyHandler {
	return &CompanyHandler{tmpl: T, companySrv: CS, postSrv: PS, sessionService: sessServ, campSess: campSess}
}

// Authenticated ...
func (ch *CompanyHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := ch.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, ch.campSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Admin handle requests on route /admin
func (ch *CompanyHandler) Admin(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_index.layout", nil)
}

// SignInUp hanlde requests on route /cmp
func (ch *CompanyHandler) SignInUp(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
}

func (ch *CompanyHandler) loggedIn(r *http.Request) bool {
	if ch.campSess == nil {
		return false
	}
	campSess := ch.campSess
	c, err := r.Cookie(campSess.UUID)
	if err != nil {
		return false
	}
	ok, err := session.Valid(c.Value, campSess.SigningKey)
	if !ok || (err != nil) {
		return false
	}
	return true
}

// Login handle requests on /cmp-login
func (ch *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {

	cookie, errc := r.Cookie("session")
	
	expireToken := time.Now().Add(time.Minute*30).Unix()
	expireCookie := time.Now().Add(time.Minute*30)

	if r.Method == http.MethodPost {

		email := r.FormValue("companyemail")
		password := r.FormValue("companypassword")

		companies, err := ch.companySrv.Companies()
		if err != nil {
			panic(err)
		}

		for _, cmp := range companies {
			if cmp.Email == email {
				fmt.Println("authentication successfull! ")

				err := bcrypt.CompareHashAndPassword([]byte(cmp.Password), []byte(password))
				if err == bcrypt.ErrMismatchedHashAndPassword {
					fmt.Println("Your email address or password is wrong")
					return
				}

				
				claims := rtoken.Claims(email, expireToken)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(email))

				if errc == http.ErrNoCookie {
					// sID, _ := uuid.NewV4()
					cookie = &http.Cookie{
						Name:  "session",
						Value: signedToken,
						Expires: expireCookie,
						Path:  "/",
					}
				}

				session := &entity.CompanySession{}
				session.UUID = cookie.Value
				session.CompanyID = cmp.ID

				_, errs := ch.companySrv.StoreSession(session)

				if len(errs) > 0 {
					panic(errs)
				}

				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

				break

			} else {
				fmt.Println("No such Company!")
			}
		}

		//io.WriteString(w, cookie.String())

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}
}

// CreateAccount handle requests on /cmp-signup-account
func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	cookie, errc := r.Cookie("session")
	
	expireToken := time.Now().Add(time.Minute*30).Unix()
	expireCookie := time.Now().Add(time.Minute*30)

	
	if r.Method == http.MethodPost {

		cmp := &entity.Company{}
		cmp.Name = r.FormValue("companyname")
		cmp.Email = r.FormValue("companyemail")
		password := r.FormValue("companypassword")
		// confirmpass := r.FormValue("confirmPassword")

		companies, _ := ch.companySrv.Companies()

		for _, company := range companies {

			if cmp.Email == company.Email {
				http.Redirect(w, r, "/cmp", http.StatusSeeOther)
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

		cmp.Password = string(hashedPassword)

		_, errs := ch.companySrv.StoreCompany(cmp)

		if len(errs) > 0 {
			panic(errs)
		}

		claims := rtoken.Claims(cmp.Email, expireToken)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(cmp.Email))

		if errc == http.ErrNoCookie {
			// sID, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session",
				Value: signedToken,
				Expires: expireCookie,
				Path:  "/",
			}
		}

		session := &entity.CompanySession{}
		session.UUID = cookie.Value
		session.CompanyID = cmp.ID

		_, errs = ch.companySrv.StoreSession(session)

		if len(errs) > 0 {
			panic(errs)
		}

		fmt.Println(cmp)

		fmt.Println(cookie.Value)

		fmt.Println("Company added to db")

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}

}

// Home handle requests on /cmp-home
func (ch *CompanyHandler) Home(w http.ResponseWriter, r *http.Request) {

	// get cookie
	_, err := r.Cookie("session")
	if err != nil {
		fmt.Println("no cookie")
		http.Redirect(w, r, "/cmp", http.StatusSeeOther)
		return
	}

	posts, _ := ch.postSrv.Posts()

	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", posts)
}

// ShowProfile handle requests on /cmp-profile
func (ch *CompanyHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_profile.html", nil)
}

// Logout Logs the company out
func (ch *CompanyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
		Path:  "/",

	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/cmp", 302)
}
