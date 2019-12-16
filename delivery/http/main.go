package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"

	"NoticeBoard/model/repository"
	"NoticeBoard/model/service"
	"NoticeBoard/delivery/http/handler"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "godisgood"
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

	userRepo := repository.NewUserRepositoryImpl(dbconn)
	userSrv := service.NewUserServiceImpl(userRepo)

	usrHandler := handler.NewUserHandler(tmpl, userSrv)

	companyRepo := repository.NewCompanyRepositoryImpl(dbconn)
	companySrv := service.NewCompanyServiceImpl(companyRepo)

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv)

	postRepo := repository.NewPostRepositoryImpl(dbconn)
	postSrv := service.NewPostServiceImpl(postRepo)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv)

	mux := http.NewServeMux()
	
	// Server CSS, JS & Images Statically.
	fs := http.FileServer(http.Dir("../../ui/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	
	mux.HandleFunc("/", Index)
	
	mux.HandleFunc("/signin", usrHandler.Signin)
	mux.HandleFunc("/signin/signup", usrHandler.Signup)
	mux.HandleFunc("/login", usrHandler.Login)
	mux.HandleFunc("/signup_account", usrHandler.CreateAccount)
	mux.HandleFunc("/home", usrHandler.Home)
	
	mux.HandleFunc("/cmp-signin", cmpHandler.Signin)
	mux.HandleFunc("/cmp-signup", cmpHandler.Signup)
	mux.HandleFunc("/cmp-login", cmpHandler.Login)
	mux.HandleFunc("/cmp-signup-account", cmpHandler.CreateAccount)
	mux.HandleFunc("/cmp-home", cmpHandler.Home)
	mux.HandleFunc("/admin", cmpHandler.Admin)

	mux.HandleFunc("/admin/post-job", postHandler.CompanyPostsNew)
	mux.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)

	http.ListenAndServe(":8080", mux)
}

func Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl.ExecuteTemplate(w, "index.layout", nil)

}
