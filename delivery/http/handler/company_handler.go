package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type CompanyHandler struct {
	tmpl       *template.Template
	companySrv model.CompanyService
	postSrv    model.PostService
}

func NewCompanyHandler(T *template.Template, CS model.CompanyService, PS model.PostService) *CompanyHandler {
	return &CompanyHandler{tmpl: T, companySrv: CS, postSrv: PS}
}

func (ch *CompanyHandler) Admin(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_index.layout", nil)
}

func (ch *CompanyHandler) SignInUp(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
}

// func (ch *CompanyHandler) Signup(w http.ResponseWriter, r *http.Request) {
// 	ch.tmpl.ExecuteTemplate(w, "cmp_signup.layout", nil)
// }

func (ch *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {
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
				http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)
				break

			} else {
				fmt.Println("No such Company!")
			}
		}
	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}
}

func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		cmp := entity.Company{}
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

		err := ch.companySrv.StoreCompany(cmp)

		if err != nil {
			panic(err)
		}

		fmt.Println(companies)

		fmt.Println(cmp)

		fmt.Println("Company added to db")

		http.Redirect(w, r, "/cmp-home", http.StatusSeeOther)

	} else {
		ch.tmpl.ExecuteTemplate(w, "company_signin_signup.html", nil)
	}

}

func (ch *CompanyHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, _ := ch.postSrv.Posts()
	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", posts)
}

func (ch *CompanyHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_profile.layout", nil)
}
