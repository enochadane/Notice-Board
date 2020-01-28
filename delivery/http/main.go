package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

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

	"github.com/amthesonofGod/Notice-Board/rtoken"

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

	// dbconn.DropTableIfExists(&entity.Role{}, &entity.CompanySession{}, &entity.UserSession{}, &entity.Post{}, &entity.User{}, &entity.Company{}, &entity.Request{}, &entity.Application{})
	errs := dbconn.CreateTable(&entity.Role{}, &entity.CompanySession{}, &entity.UserSession{}, &entity.Post{}, &entity.User{}, &entity.Company{}, &entity.Request{}, &entity.Application{}).GetErrors()

	if errs != nil {
		return errs
	}

	return nil
}

func main() {

	csrfSignKey := []byte(rtoken.GenerateRandomID(32))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbconn, err := gorm.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	defer dbconn.Close()

	// createTables(dbconn)

	// if err := dbconn.Ping(); err != nil {
	// 	panic(err)
	// }

	userSessionRepo := repository.NewSessionGormRepo(dbconn)
	userSessionsrv := service.NewSessionService(userSessionRepo)

	companySessionRepo := repositoryCamp.NewSessionGormRepoCamp(dbconn)
	companySessionSrv := serviceCamp.NewSessionServiceCamp(companySessionRepo)

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


	requestHandler := handler.NewRequestHandler(tmpl, requestSrv, postSrv, userSrv, csrfSignKey)

	applicationHandler := handler.NewApplicationHandler(tmpl, applicationSrv, userSrv, postSrv, csrfSignKey)

	//(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession)
	sessCamp := configSessCamp()

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv, companySessionSrv, sessCamp, csrfSignKey)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv, csrfSignKey)
	sess := configSess()

	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv, userSessionsrv, sess, csrfSignKey)


	//r := mux.NewRouter()

	// r := http.NewServeMux()

	// Server CSS, JS & Images Statically.
	fs := http.FileServer(http.Dir("../../ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	//r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))

	http.HandleFunc("/", usrHandler.Index)
	http.HandleFunc("/login", usrHandler.Login)
	http.HandleFunc("/signup", usrHandler.CreateAccount)
	http.Handle("/home", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(usrHandler.Home))))
	http.Handle("/logout", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(usrHandler.Logout))))

	http.HandleFunc("/admin", cmpHandler.SignInUp)
	http.HandleFunc("/admin/login", cmpHandler.Login)
	http.HandleFunc("/admin/signup", cmpHandler.CreateAccount)
	http.Handle("/admin/home", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.Home))))
	http.Handle("/admin/profile", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.ShowProfile))))
	//http.HandleFunc("/admin/dashboard", cmpHandler.Admin)
	http.Handle("/admin/logout",cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.Logout))))

	http.Handle("/admin/new-post", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostsNew))))
	http.Handle("/admin/posts", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPosts))))
	http.Handle("/admin/posts/update", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostUpdate))))
	http.Handle("/admin/posts/delete", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostDelete))))

	http.Handle("/apply", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.Apply))))
	http.Handle("/applications", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.Applications))))
	http.Handle("/admin/received/applications", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(applicationHandler.CompanyReceivedApplications))))
	http.Handle("/admin/received/applications/details", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationDetails))))
	http.Handle("/applications/update", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationUpdate))))
	http.Handle("/applications/delete", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationDelete))))

	http.Handle("/join", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.Join))))
	http.Handle("/requests", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.Requests))))
	http.Handle("/admin/received/requests", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(requestHandler.CompanyReceivedRequests))))
	http.Handle("/requests/update", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.RequestUpdate))))
	http.Handle("/requests/delete", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.RequestDelete))))

	//port := fmt.Sprintf(":%s", os.Getenv("HPORT"))

	http.ListenAndServe(":8080", nil)
}

func configSess() *entity.UserSession {
	tokenExpires := time.Now().Add(time.Minute * 30).Unix()
	sessionID := rtoken.GenerateRandomID(32)
	signingString, err := rtoken.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	signingKey := []byte(signingString)

	return &entity.UserSession{
		Expires:    tokenExpires,
		SigningKey: signingKey,
		UUID:       sessionID,
	}
}

func configSessCamp() *entity.CompanySession {
	tokenExpires := time.Now().Add(time.Minute * 30).Unix()
	sessionID := rtoken.GenerateRandomID(32)
	signingString, err := rtoken.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	signingKey := []byte(signingString)

	return &entity.CompanySession{
		Expires:    tokenExpires,
		SigningKey: signingKey,
		UUID:       sessionID,
	}
}
