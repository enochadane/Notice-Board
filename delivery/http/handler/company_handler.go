package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"

	"github.com/amthesonofGod/Notice-Board/post"

	uuid "github.com/satori/go.uuid"

	"github.com/amthesonofGod/Notice-Board/rtoken"
	"github.com/amthesonofGod/Notice-Board/session"
)

func king(){}

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

// LoginC handle requests on /cmp-login
func (ch *CompanyHandler) LoginC(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		email := r.FormValue("companyemail")
		password := r.FormValue("companypassword")

		companies, err := ch.companySrv.Companies()
		if err != nil {
			panic(err)
		}

		for _, cmp := range companies {
			if cmp.Email == email && cmp.Password == password {
				fmt.Println("authentication successfull! ")

				// if errc == http.ErrNoCookie {
				// 	sID, _ := uuid.NewV4()
				// 	cookie = &http.Cookie{
				// 		Name:  "session",
				// 		Value: sID.String(),
				// 		Path:  "/",
				// 	}
				// }

				c := &cmp
				ch.loggedInUserCamp = c
				claims := rtoken.Claims(c.Email, ch.campSess.Expires)
				session.Create(claims, ch.campSess.UUID, ch.campSess.SigningKey, w)
				newSess, errs := ch.sessionService.StoreSessionCamp(ch.campSess)
				// session.CompanyName = cmp.Name

				if len(errs) > 0 {
					panic(errs)
				}
				ch.campSess = newSess
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

// CreateAccountC handle requests on /cmp-signup-account
func (ch *CompanyHandler) CreateAccountC(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if r.Method == http.MethodPost {

		cmp := &entity.Company{}
		cmp.Name = r.FormValue("companyname")
		cmp.Email = r.FormValue("companyemail")
		cmp.Password = r.FormValue("companypassword")
		// confirmpass := r.FormValue("confirmPassword")

		companies, _ := ch.companySrv.Companies()

		for _, company := range companies {

			if cmp.Email == company.Email {
				http.Redirect(w, r, "/cmp", http.StatusSeeOther)
				fmt.Println("This Email is already in use! ")
				return
			}
		}

		_, errs := ch.companySrv.StoreCompany(cmp)

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

		session := &entity.CompanySession{}
		session.UUID = cookie.Value
		session.CompanyID = cmp.ID

		_, errs = ch.companySrv.StoreSession(session)

		if len(errs) > 0 {
			panic(errs)
		}

		fmt.Println(cmp)

		fmt.Println("Company added to db")

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}

}

// HomeC handle requests on /cmp-home
func (ch *CompanyHandler) HomeC(w http.ResponseWriter, r *http.Request) {

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

// LogoutC Logs the company out
func (ch *CompanyHandler) LogoutC(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logged-in")
	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
		}
		// session := data.Session{Uuid: cookie.Value}
		// session.DeleteByUUID()
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/cmp", 302)
}
