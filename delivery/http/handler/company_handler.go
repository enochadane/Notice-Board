package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

<<<<<<< HEAD
	"github.com/amthesonofGod/Notice-Board/permission"

=======
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb
	// "time"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/form"
	"github.com/amthesonofGod/Notice-Board/post"

	"golang.org/x/crypto/bcrypt"

	"github.com/amthesonofGod/Notice-Board/rtoken"
	// "github.com/dgrijalva/jwt-go"
	"github.com/amthesonofGod/Notice-Board/session"
)

// CompanyHandler handles company handler admin requests
type CompanyHandler struct {
<<<<<<< HEAD
	tmpl             *template.Template
	companySrv       company.CompanyService
	postSrv          post.PostService
	sessionService   company.SessionServiceCamp
	campSess         *entity.CompanySession
	loggedInUserCamp *entity.Company
	companyRole      company.RoleService
	csrfSignKey      []byte
=======
	tmpl       			*template.Template
	companySrv 			company.CompanyService
	postSrv   			post.PostService
	sessionService   	company.SessionServiceCamp
	campSess         	*entity.CompanySession
	loggedInUserCamp 	*entity.Company
	companyRole			company.RoleService
	csrfSignKey    	  	[]byte
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb
}

var currentCompUser *entity.Company

type cntextKey string

var ctxCompanySessionKey = cntextKey("signed_in_company_session")

// NewCompanyHandler initializes and returns new NewCompanyHandler
func NewCompanyHandler(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession, csKey []byte) *CompanyHandler {
	return &CompanyHandler{tmpl: T, companySrv: CS, postSrv: PS, sessionService: sessServ, campSess: campSess, csrfSignKey:csKey}
}

// Authenticated ...
func (ch *CompanyHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := ch.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxCompanySessionKey, ch.campSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Authorized checks if a user has proper authority to access a give route
func (ch *CompanyHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ch.loggedInUserCamp == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := ch.companySrv.CompanyRoles(ch.loggedInUserCamp)
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
			ok, err := rtoken.ValidCSRF(r.FormValue("_csrf"), ch.csrfSignKey)
			if !ok || (err != nil) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
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

<<<<<<< HEAD
// LoginC handle requests on /cmp-login
func (ch *CompanyHandler) LoginC(w http.ResponseWriter, r *http.Request) {
=======
// Login handle requests on /cmp-login
func (ch *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb

	token, err := rtoken.CSRFToken(ch.csrfSignKey)
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
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", loginForm)
		return
	}

	if r.Method == http.MethodPost {

		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		cmp, errs := ch.companySrv.CompanyByEmail(r.FormValue("companyemail"))
		if len(errs) > 0 {
			fmt.Println("email errrrrr")
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", loginForm)
			return
		}

		//match := CheckPasswordHash(r.FormValue("companypassword"), cmp.Password)
		//fmt.Println("Match:   ", match)

		err := bcrypt.CompareHashAndPassword([]byte(cmp.Password), []byte(r.FormValue("companypassword")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			fmt.Println("pass err")
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", loginForm)
			return
		}

		ch.loggedInUserCamp = cmp
		currentCompUser = ch.loggedInUserCamp
		claims := rtoken.Claims(cmp.Email, ch.campSess.Expires)

		session.Create(claims, ch.campSess.UUID, ch.campSess.SigningKey, w)

		newSess, errs := ch.sessionService.StoreSessionCamp(ch.campSess)

		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to store session")
			ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", loginForm)
			return
		}

		ch.campSess = newSess

		roles, _ := ch.companySrv.CompanyRoles(cmp)
		if ch.checkAdmin(roles) {
			http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
			return
		}
<<<<<<< HEAD
=======

		http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb

		http.Redirect(w, r, "/admin/home", http.StatusSeeOther)

<<<<<<< HEAD
	}
}

// CreateAccountC handle requests on /admin/signup
func (ch *CompanyHandler) CreateAccountC(w http.ResponseWriter, r *http.Request) {

=======
// CreateAccount handle requests on /admin/signup
func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb
	token, err := rtoken.CSRFToken(ch.csrfSignKey)
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
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", signUpForm)
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
		singnUpForm.Required("companyname", "companyemail", "companypassword", "confirmpassword")
		singnUpForm.MatchesPattern("companyemail", form.EmailRX)
		singnUpForm.MinLength("companypassword", 8)
		singnUpForm.PasswordMatches("companypassword", "confirmpassword")
		singnUpForm.CSRF = token

		// If there are any errors, redisplay the signup form.
		if !singnUpForm.Valid() {
			ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", singnUpForm)
			return
		}

		eExists := ch.companySrv.EmailExists(r.FormValue("companyemail"))
		if eExists {
			singnUpForm.VErrors.Add("email", "Email Already Exists")
			ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", singnUpForm)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("companypassword")), 12)
		if err != nil {
			singnUpForm.VErrors.Add("password", "Password Could not be stored")
			ch.tmpl.ExecuteTemplate(w, "company_signin_signuphtml.html", singnUpForm)
			return
		}

		//role, errs := ch.companyRole.RoleByName("ADMIN")
		//
		//if len(errs) > 0 {
		//	singnUpForm.VErrors.Add("role", "could not assign role to the user")
		//	ch.tmpl.ExecuteTemplate(w, "signup.layout", singnUpForm)
		//	return
		//}

		company := &entity.Company{
			Name:     r.FormValue("companyname"),
			Email:    r.FormValue("companyemail"),
			Password: string(hashedPassword),
<<<<<<< HEAD
			RoleID:   0,
=======
			RoleID: 0,
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb
		}

		_, errs := ch.companySrv.StoreCompany(company)

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)

<<<<<<< HEAD
	}
}

// HomeC handle requests on /cmp-home
func (ch *CompanyHandler) HomeC(w http.ResponseWriter, r *http.Request) {
=======
	} 
}

// Home handle requests on /cmp-home
func (ch *CompanyHandler) Home(w http.ResponseWriter, r *http.Request) {
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb

	posts, _ := ch.postSrv.Posts()

	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", posts)
}

// ShowProfile handle requests on /cmp-profile
func (ch *CompanyHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_profile.html", nil)
}

// LogoutC Logs the company out
func (ch *CompanyHandler) LogoutC(w http.ResponseWriter, r *http.Request) {

	session.Remove(ch.campSess.UUID, w)
	ch.sessionService.DeleteSessionCamp(ch.campSess.UUID)

	http.Redirect(w, r, "/admin", 302)
<<<<<<< HEAD
}

func (ch *CompanyHandler) checkAdmin(rs []entity.Role) bool {
	for _, r := range rs {
		if strings.ToUpper(r.Name) == strings.ToUpper("Admin") {
			return true
		}
	}
	return false
=======
>>>>>>> 997df0981b2ffe30de1cb2328c8e127e034eedeb
}

func (ch *CompanyHandler) checkAdmin(rs []entity.Role) bool {
	for _, r := range rs {
		if strings.ToUpper(r.Name) == strings.ToUpper("Admin") {
			return true
		}
	}
	return false
}


