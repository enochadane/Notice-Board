package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"io"

<<<<<<< HEAD
	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"

	"github.com/amthesonofGod/Notice-Board/post"

	uuid "github.com/satori/go.uuid"
=======
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/model"
	"github.com/amthesonofGod/Notice-Board/post"

	"github.com/satori/go.uuid"
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
)

// CompanyHandler handles company handler admin requests
type CompanyHandler struct {
<<<<<<< HEAD
	tmpl       *template.Template
	companySrv company.CompanyService
	postSrv    post.PostService
}

// NewCompanyHandler initializes and returns new NewCompanyHandler
func NewCompanyHandler(T *template.Template, CS company.CompanyService, PS post.PostService) *CompanyHandler {
=======
	tmpl		*template.Template
	companySrv 	model.CompanyService
	postSrv		post.PostService
}

// NewCompanyHandler initializes and returns new NewCompanyHandler
func NewCompanyHandler(T *template.Template, CS model.CompanyService, PS post.PostService) *CompanyHandler {
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
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
<<<<<<< HEAD

				if errc == http.ErrNoCookie {
					sID, _ := uuid.NewV4()
					cookie = &http.Cookie{
						Name:  "session",
						Value: sID.String(),
						Path:  "/",
					}
				}

				session := &entity.CompanySession{}
=======
				
				if errc == http.ErrNoCookie {
					sID, _ := uuid.NewV4()
					cookie = &http.Cookie{
						Name: "session",
						Value: sID.String(),
						Path: "/",
					}
				}

				session := &entity.Session{}
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
				session.UUID = cookie.Value
				session.CompanyID = cmp.ID
				// session.CompanyName = cmp.Name

				_, errs := ch.companySrv.StoreSession(session)

				if len(errs) > 0 {
					panic(errs)
				}

				fmt.Println(session.UUID)
<<<<<<< HEAD

				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

=======
				
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)
				
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
				break

			} else {
				fmt.Println("No such Company!")
			}
		}

<<<<<<< HEAD
=======
		
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
		io.WriteString(w, cookie.String())

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}
}

// CreateAccount handle requests on /cmp-signup-account
func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if r.Method == http.MethodPost {
<<<<<<< HEAD

=======
		
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
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
<<<<<<< HEAD
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
=======
		}
		fmt.Println(companies)
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f

		fmt.Println(cmp)

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
<<<<<<< HEAD

=======
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
	
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", posts)
}

// ShowProfile handle requests on /cmp-profile
func (ch *CompanyHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_profile.html", nil)
}

<<<<<<< HEAD
// Logout Logs the company out
=======
// Logout Logs the user out
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
func (ch *CompanyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("logged-in")
	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
<<<<<<< HEAD
			Name:  "logged-in",
=======
			Name: "logged-in",
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
			Value: "0",
		}
		// session := data.Session{Uuid: cookie.Value}
		// session.DeleteByUUID()
	}

	http.SetCookie(w, cookie)
<<<<<<< HEAD
	http.Redirect(w, r, "/cmp", 302)
}
=======
	http.Redirect(w, r, "/", 302)
}

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
