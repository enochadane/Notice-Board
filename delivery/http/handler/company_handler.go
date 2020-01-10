package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"io"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/model"
	"github.com/amthesonofGod/Notice-Board/post"

	"github.com/satori/go.uuid"
)

// CompanyHandler handles company handler admin requests
type CompanyHandler struct {
	tmpl		*template.Template
	companySrv 	model.CompanyService
	postSrv		post.PostService
}

// NewCompanyHandler initializes and returns new NewCompanyHandler
func NewCompanyHandler(T *template.Template, CS model.CompanyService, PS post.PostService) *CompanyHandler {
	return &CompanyHandler{tmpl: T, companySrv: CS, postSrv: PS}
}

// Admin handle requests on route /admin
func (ch *CompanyHandler) Admin(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_index.layout", nil)
}

// SignInUp hanlde requests on route /cmp
func (ch *CompanyHandler) SignInUp(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
}

// Login handle requests on /cmp-login
func (ch *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {

	cookie, errc := r.Cookie("session")

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
				
				if errc == http.ErrNoCookie {
					sID, _ := uuid.NewV4()
					cookie = &http.Cookie{
						Name: "session",
						Value: sID.String(),
						Path: "/",
					}
				}

				session := &entity.Session{}
				session.UUID = cookie.Value
				session.CompanyID = cmp.ID
				// session.CompanyName = cmp.Name

				_, errs := ch.companySrv.StoreSession(session)

				if len(errs) > 0 {
					panic(errs)
				}

				fmt.Println(session.UUID)
				
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)
				
				break
			
			} else {
				fmt.Println("No such Company!")
			}
		}

		
		io.WriteString(w, cookie.String())

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}
}

// CreateAccount handle requests on /cmp-signup-account
func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	
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
		fmt.Println(companies)

		fmt.Println(cmp)

		fmt.Println("Company added to db")

		http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

		
	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}

}

// Home handle requests on /cmp-home
func (ch *CompanyHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, _ := ch.postSrv.Posts()
	// for _, post := range posts {
	// 	cmp, _ := ch.companySrv.Company(post.CompanyID)
	// 	// if len(err) > 0 {
	// 	// 	panic(err)
	// 	// }

	// 	m = make(map[entity.Post]string)
	// 	m[post] = cmp.Name
	// 	// fmt.Println(m[post])

	// 	// fmt.Println(cmp.Name)

	// }
	
	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", posts)
}

// ShowProfile handle requests on /cmp-profile
func (ch *CompanyHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_profile.html", nil)
}

// Logout Logs the user out
func (ch *CompanyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logged-in")
	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name: "logged-in",
			Value: "0",
		}
		// session := data.Session{Uuid: cookie.Value}
		// session.DeleteByUUID()
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

