package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/amthesonofGod/Notice-Board/entity"

	api "github.com/amthesonofGod/Notice-Board/rest-api"

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
	password = "kingo"
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

	requestHandler := handler.NewRequestHandler(tmpl, requestSrv, postSrv, userSrv)

	applicationHandler := handler.NewApplicationHandler(tmpl, applicationSrv, userSrv, postSrv)

	//(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession)
	sessCamp := configSessCamp()

	cmpHandler := handler.NewCompanyHandler(tmpl, companySrv, postSrv, companySessionSrv, sessCamp)

	postHandler := handler.NewCompanyPostHandler(tmpl, postSrv, companySrv)
	sess := configSess()

	usrHandler := handler.NewUserHandler(tmpl, userSrv, postSrv, userSessionsrv, sess, csrfSignKey)

	//r := mux.NewRouter()

	r := http.NewServeMux()

	// Server CSS, JS & Images Statically.
	fs := http.FileServer(http.Dir("../../ui/assets"))
	r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	//r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))

	r.HandleFunc("/", usrHandler.Index)
	r.HandleFunc("/login", usrHandler.Login)
	r.HandleFunc("/signup", usrHandler.CreateAccount)
	r.Handle("/home", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(usrHandler.Home))))
	r.Handle("/logout", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(usrHandler.Logout))))

	r.HandleFunc("/admin", cmpHandler.SignInUp)
	r.HandleFunc("/admin/login", cmpHandler.LoginC)
	r.HandleFunc("/admin/signup", cmpHandler.CreateAccountC)
	r.Handle("/admin/home", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.HomeC))))
	r.Handle("/admin/profile", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.ShowProfile))))
	//r.HandleFunc("/admin/dashboard", cmpHandler.Admin)
	r.Handle("/admin/logout", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(cmpHandler.LogoutC))))

	r.Handle("/admin/posts/new", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostsNew))))
	r.Handle("/admin/posts", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPosts))))
	r.Handle("/admin/posts/update", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostUpdate))))
	r.Handle("/admin/posts/delete", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(postHandler.CompanyPostDelete))))

	r.Handle("/apply", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.Apply))))
	r.Handle("/applications", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.Applications))))
	r.Handle("/received/applications", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(applicationHandler.CompanyReceivedApplications))))
	r.Handle("/received/applications/details", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationDetails))))
	r.Handle("/applications/update", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationUpdate))))
	r.Handle("/applications/delete", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(applicationHandler.ApplicationDelete))))

	r.Handle("/join", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.Join))))
	r.Handle("/requests", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.Requests))))
	r.Handle("/received/requests", cmpHandler.Authenticated(cmpHandler.Authorized(http.HandlerFunc(requestHandler.CompanyReceivedRequests))))
	r.Handle("/requests/update", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.RequestUpdate))))
	r.Handle("/requests/delete", usrHandler.Authenticated(usrHandler.Authorized(http.HandlerFunc(requestHandler.RequestDelete))))

	aplAPI := api.NewApplicationHandler(tmpl, applicationSrv, userSrv, postSrv)

	router := httprouter.New()

	router.GET("/v1/admin/applications", aplAPI.Applications)                         //adminRoleHandler.GetRoles
	router.GET("/v1/admin/recieved/applications", aplAPI.CompanyReceivedApplications) //adminCommentHandler.GetSingleComment
	router.GET("/v1/admin/applications/:id", aplAPI.ApplicationDetails)               //adminCommentHandler.GetComments
	router.POST("/v1/admin/applications/apply/:id", aplAPI.Apply)                     //adminCommentHandler.PutComment
	router.PUT("/v1/admin/applications/update", aplAPI.ApplicationUpdate)             //adminCommentHandler.PostComment
	router.DELETE("/v1/admin/applications/delete", aplAPI.ApplicationDelete)          // adminCommentHandler.DeleteComment

	//port := fmt.Sprintf(":%s", os.Getenv("HPORT"))

	http.ListenAndServe(":8080", router)
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
