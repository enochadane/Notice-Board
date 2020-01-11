package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/amthesonofGod/Notice-Board/entity"

	repository "github.com/amthesonofGod/Notice-Board/user/repository"
	service "github.com/amthesonofGod/Notice-Board/user/service"

	repositoryCamp "github.com/amthesonofGod/Notice-Board/company/repositoryCamp"
	serviceCamp "github.com/amthesonofGod/Notice-Board/company/serviceCamp"

	postRepos "github.com/amthesonofGod/Notice-Board/post/repository"
	postServ "github.com/amthesonofGod/Notice-Board/post/service"

	appRepos "github.com/amthesonofGod/Notice-Board/application/repository"
	appServ "github.com/amthesonofGod/Notice-Board/application/service"

	reqRepos "github.com/amthesonofGod/Notice-Board/request/repository"
	reqServ "github.com/amthesonofGod/Notice-Board/request/service"

	"github.com/amthesonofGod/Notice-Board/delivery/http/handler"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func createTables(dbconn *gorm.DB) []error {

	// dbconn.DropTableIfExists(&entity.Session{})
	errs := dbconn.CreateTable( &entity.Request{}, &entity.Application{}).GetErrors()
	// errs := dbconn.CreateTable(&entity.CompanySession{}, &entity.UserSession{}, &entity.Post{}, &entity.User{}, &entity.Company{}).GetErrors()

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


	requestHandler := handler.NewRequestHandler(tmpl, requestSrv, postSrv, userSrv)

	applicationHandler := handler.NewApplicationHandler(tmpl, applicationSrv, userSrv, postSrv)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv)

	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv)

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv)

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

	r.HandleFunc("/admin/posts/new", postHandler.CompanyPostsNew)
	r.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)
	r.HandleFunc("/cmp/posts/update", postHandler.CompanyPostUpdate)
	r.HandleFunc("/cmp/posts/delete", postHandler.CompanyPostDelete)

	r.HandleFunc("/job/apply", applicationHandler.Apply)
	r.HandleFunc("/applications", applicationHandler.Applications)
	r.HandleFunc("/received/applications", applicationHandler.CompanyReceivedApplications)
	r.HandleFunc("/received/applications/details", applicationHandler.ApplicationDetails)
	r.HandleFunc("/user/applications/update", applicationHandler.ApplicationUpdate)
	r.HandleFunc("/user/applications/delete", applicationHandler.ApplicationDelete)

	r.HandleFunc("/event/join", requestHandler.Join)
	r.HandleFunc("/requests", requestHandler.Requests)
	r.HandleFunc("/received/requests", requestHandler.CompanyReceivedRequests)
	r.HandleFunc("/user/requests/update", requestHandler.RequestUpdate)
	r.HandleFunc("/user/requests/delete", requestHandler.RequestDelete)

	http.ListenAndServe(":8080", r)
}

