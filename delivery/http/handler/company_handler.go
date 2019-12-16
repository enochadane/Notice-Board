package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type CompanyHandler struct {
	tmpl		*template.Template
	companySrv 	model.CompanyService
}

func NewCompanyHandler(T *template.Template, CS model.CompanyService) *CompanyHandler {
	return &CompanyHandler{tmpl: T, companySrv: CS}
}

func (ch *CompanyHandler) Admin(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_index.layout", nil)
}

func (ch *CompanyHandler) Signin(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_signin.layout", nil)
}

func (ch *CompanyHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_signup.layout", nil)
}

func (ch *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

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
		ch.tmpl.ExecuteTemplate(w, "cmp_signin.layout", nil)
	}
}

func (ch *CompanyHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == http.MethodPost {
		
		cmp := entity.Company{}
		cmp.Name = r.FormValue("username")
		cmp.Email = r.FormValue("useremail")
		cmp.Password = r.FormValue("password")
		confirmpass := r.FormValue("confirmPassword")

		companies, _ := ch.companySrv.Companies()

		// for _, company := range companies {
			
		// 	if cmp.Email == company.Email {
		// 		http.Redirect(w, r, "/cmp_signup", http.StatusSeeOther)
		// 		fmt.Println("This Email is already in use! ")
		// 		return
		// 	}
		// }

		if cmp.Password == confirmpass {

			err := ch.companySrv.StoreCompany(cmp)

			if err != nil {
				panic(err)
			}

			fmt.Println(companies)

			fmt.Println(cmp)

			fmt.Println("Company added to db")

			http.Redirect(w, r, "/cmp-login", http.StatusSeeOther)

		} else {
			fmt.Println("Password doesn't match! ")
		}
		
	} else {
		ch.tmpl.ExecuteTemplate(w, "cmp_signup.layout", nil)
	}

}

func (ch *CompanyHandler) Home(w http.ResponseWriter, r *http.Request) {
	ch.tmpl.ExecuteTemplate(w, "cmp_home.layout", nil)
}