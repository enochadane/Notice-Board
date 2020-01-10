package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/amthesonofGod/Notice-Board/entity"

	"github.com/amthesonofGod/Notice-Board/model/repository"
	"github.com/amthesonofGod/Notice-Board/model/service"

	postRepos "github.com/amthesonofGod/Notice-Board/post/repository"
	postServ "github.com/amthesonofGod/Notice-Board/post/service"

	"github.com/amthesonofGod/Notice-Board/delivery/http/handler"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "kingo"
	dbname   = "noticeboard"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("../../ui/templates/*"))
}

func createTables(dbconn *gorm.DB) []error {

	// dbconn.DropTableIfExists(&entity.Session{})
	errs := dbconn.CreateTable(&entity.Post{}, &entity.Session{}, &entity.User{}, &entity.Company{}).GetErrors()

	if errs != nil {
		return errs
	}

	return nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbconn, err := gorm.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	defer dbconn.Close()

	createTables(dbconn)

	// if err := dbconn.Ping(); err != nil {
	// 	panic(err)
	// }

	companyRepo := repository.NewCompanyGormRepo(dbconn)
	companySrv := service.NewCompanyService(companyRepo)

	postRepo := postRepos.NewPostGormRepo(dbconn)
	postSrv := postServ.NewPostService(postRepo)

	userRepo := repository.NewUserGormRepo(dbconn)
	userSrv := service.NewUserService(userRepo)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv)

	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv)

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv)

	// dbconn.Model(company).Find(company)
	// for _, p := range postSrv.Posts {
	// 	dbconn.Model(company).Association("Posts").Append(p)
	// }

	r := mux.NewRouter()

	// Server CSS, JS & Images Statically.
	// fs := http.FileServer(http.Dir("../../ui/assets"))
	// r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))

	r.HandleFunc("/", usrHandler.Index)
	r.HandleFunc("/login", usrHandler.Login)
	r.HandleFunc("/signup_account", usrHandler.CreateAccount)
	r.HandleFunc("/home", usrHandler.Home)

	r.HandleFunc("/cmp", cmpHandler.SignInUp)
	r.HandleFunc("/cmp-login", cmpHandler.Login)
	r.HandleFunc("/cmp-signup-account", cmpHandler.CreateAccount)
	r.HandleFunc("/cmp-home", cmpHandler.Home)
	r.HandleFunc("/cmp-profile", cmpHandler.ShowProfile)

	r.HandleFunc("/admin", cmpHandler.Admin)
	r.HandleFunc("/admin/posts/new", postHandler.CompanyPostsNew)
	r.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)

	http.ListenAndServe(":8080", r)
}
