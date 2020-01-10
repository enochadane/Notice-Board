package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/amthesonofGod/Notice-Board/entity"
<<<<<<< HEAD

	repository "github.com/amthesonofGod/Notice-Board/User/repository"
	service "github.com/amthesonofGod/Notice-Board/User/service"

	repositoryCamp "github.com/amthesonofGod/Notice-Board/company/repositoryCamp"
	serviceCamp "github.com/amthesonofGod/Notice-Board/company/serviceCamp"

	postRepos "github.com/amthesonofGod/Notice-Board/post/repository"
	postServ "github.com/amthesonofGod/Notice-Board/post/service"

	appRepos "github.com/amthesonofGod/Notice-Board/application/repository"
	appServ "github.com/amthesonofGod/Notice-Board/application/service"

	reqRepos "github.com/amthesonofGod/Notice-Board/request/repository"
	reqServ "github.com/amthesonofGod/Notice-Board/request/service"

=======

	"github.com/amthesonofGod/Notice-Board/model/repository"
	"github.com/amthesonofGod/Notice-Board/model/service"

	postRepos "github.com/amthesonofGod/Notice-Board/post/repository"
	postServ "github.com/amthesonofGod/Notice-Board/post/service"

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
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
<<<<<<< HEAD

	// dbconn.DropTableIfExists(&entity.Session{})
	// errs := dbconn.CreateTable(&entity.Application{}, &entity.Request{}).GetErrors()
	errs := dbconn.CreateTable(&entity.CompanySession{}, &entity.UserSession{}, &entity.Post{}, &entity.User{}, &entity.Company{}).GetErrors()

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
=======

	// dbconn.DropTableIfExists(&entity.Session{})
	errs := dbconn.CreateTable(&entity.Post{}, &entity.Session{}, &entity.User{}, &entity.Company{}).GetErrors()
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f

	if errs != nil {
		return errs
	}

<<<<<<< HEAD
	fmt.Println("DB connection established")

	defer dbconn.Close()

	createTables(dbconn)
=======
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
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f

	// if err := dbconn.Ping(); err != nil {
	// 	panic(err)
	// }

<<<<<<< HEAD
	companyRepo := repositoryCamp.NewCompanyGormRepo(dbconn)
	companySrv := serviceCamp.NewCompanyService(companyRepo)

	postRepo := postRepos.NewPostGormRepo(dbconn)
	postSrv := postServ.NewPostService(postRepo)

	userRepo := repository.NewUserGormRepo(dbconn)
	userSrv := service.NewUserService(userRepo)

	applicationRepo := appRepos.NewApplicationGormRepo(dbconn)
	applicationSrv := appServ.NewApplicationService(applicationRepo)

	requestRepo := reqRepos.NewRequestGormRepo(dbconn)
	requestSrv := reqServ.NewRequestService(requestRepo)

	requestHandler := handler.NewRequestHandler(tmpl, requestSrv, postSrv)

	applicationHandler := handler.NewApplicationHandler(tmpl, applicationSrv, postSrv)

=======
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

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv)

	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv)

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv)

<<<<<<< HEAD
=======
	// dbconn.Model(company).Find(company)
	// for _, p := range postSrv.Posts {
	// 	dbconn.Model(company).Association("Posts").Append(p)
	// }

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
	r := mux.NewRouter()

	// Server CSS, JS & Images Statically.
	// fs := http.FileServer(http.Dir("../../ui/assets"))
	// r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))

	r.HandleFunc("/", usrHandler.Index)
	r.HandleFunc("/login", usrHandler.Login)
	r.HandleFunc("/signup-account", usrHandler.CreateAccount)
	r.HandleFunc("/home", usrHandler.Home)

	r.HandleFunc("/cmp", cmpHandler.SignInUp)
	r.HandleFunc("/cmp-login", cmpHandler.Login)
	r.HandleFunc("/cmp-signup-account", cmpHandler.CreateAccount)
	r.HandleFunc("/cmp-home", cmpHandler.Home)
	r.HandleFunc("/cmp-profile", cmpHandler.ShowProfile)
	r.HandleFunc("/admin", cmpHandler.Admin)

<<<<<<< HEAD
	r.HandleFunc("/admin/posts/new", postHandler.CompanyPostsNew)
	r.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)

	r.HandleFunc("/job/apply", applicationHandler.Apply)
	r.HandleFunc("/applicatons", applicationHandler.Applications)

	r.HandleFunc("/event/join", requestHandler.Join)
	r.HandleFunc("/requests", requestHandler.Requests)

=======
	r.HandleFunc("/admin", cmpHandler.Admin)
	r.HandleFunc("/admin/posts/new", postHandler.CompanyPostsNew)
	r.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)

>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
	http.ListenAndServe(":8080", r)
}
