package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"

	"github.com/amthesonofGod/Notice-Board/model/repository"
	"github.com/amthesonofGod/Notice-Board/model/service"
	"github.com/amthesonofGod/Notice-Board/delivery/http/handler"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "noticeboard"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("../../ui/templates/*"))
}

func main()  {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	dbconn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	// tmpl := template.Must(template.ParseGlob("../../ui/templates/*"))


	companyRepo := repository.NewCompanyRepositoryImpl(dbconn)
	companySrv := service.NewCompanyServiceImpl(companyRepo)

	postRepo := repository.NewPostRepositoryImpl(dbconn)
	postSrv := service.NewPostServiceImpl(postRepo)

	userRepo := repository.NewUserRepositoryImpl(dbconn)
	userSrv := service.NewUserServiceImpl(userRepo)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv)


	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv)

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv)

	// mux := http.NewServeMux()
	r := mux.NewRouter()
	
	// Server CSS, JS & Images Statically.
	// fs := http.FileServer(http.Dir("../../ui/assets"))
	// r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))
	
	r.HandleFunc("/", Index)
	
	// r.HandleFunc("/signin", usrHandler.Signin)
	// r.HandleFunc("/signup", usrHandler.Signup)
	r.HandleFunc("/login", usrHandler.Login)
	r.HandleFunc("/signup_account", usrHandler.CreateAccount)
	r.HandleFunc("/home", usrHandler.Home)
	r.HandleFunc("/user-profile", usrHandler.ShowUserProfile)
	
	// r.HandleFunc("/cmp-signin", cmpHandler.Signin)
	// r.HandleFunc("/cmp-signup", cmpHandler.Signup)
	r.HandleFunc("/cmp", cmpHandler.SignInUp)
	r.HandleFunc("/cmp-login", cmpHandler.Login)
	r.HandleFunc("/cmp-signup-account", cmpHandler.CreateAccount)
	r.HandleFunc("/cmp-home", cmpHandler.Home)
	r.HandleFunc("/admin", cmpHandler.Admin)
	r.HandleFunc("/cmp-profile", cmpHandler.ShowProfile)

	

	r.HandleFunc("/admin/post-job", postHandler.CompanyPostsNew)
	r.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)

	http.ListenAndServe(":8080", r)
}

func Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl.ExecuteTemplate(w, "index_signin_signup.html", nil)

}
