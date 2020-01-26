package handler

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/amthesonofGod/Notice-Board/entity"

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

	"github.com/amthesonofGod/Notice-Board/rtoken"

	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "kingo"
	dbname   = "noticeboard"
)

var tmpl *template.Template
var usrHandler *UserHandler

var cmpHandler *CompanyHandler
var postHandler *CompanyPostHandler
var applicationHandler *ApplicationHandler
var requestHandler *RequestHandler

func init() {

	tmpl = template.Must(template.ParseGlob("../../../ui/templates/*"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbconn, err := gorm.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	defer dbconn.Close()

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

	requestHandler = NewRequestHandler(tmpl, requestSrv, postSrv)

	applicationHandler = NewApplicationHandler(tmpl, applicationSrv, postSrv)

	//(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession)
	sessCamp := configSessCamp()

	postHandler = NewCompanyPostHandler(tmpl, postSrv, companySrv)
	sess := configSess()

	usrHandler = NewUserHandler(tmpl, userSrv, postSrv, userSessionsrv, sess)

	cmpHandler = NewCompanyHandler(tmpl, companySrv, postSrv, companySessionSrv, sessCamp)

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

//Test for User

//var u *UserHandler

func TestIndex(t *testing.T) {
	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	usrHandler.Index(r, req)
	res := r.Result()

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestLogin(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	usrHandler.Login(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}
}

func TestCreateAccount(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	usrHandler.CreateAccount(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}
}

func TestHome(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	usrHandler.Home(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}
}

func TestLogout(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	usrHandler.Logout(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

//Test for Company

var c *CompanyHandler

func TestAdmin(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	cmpHandler.Admin(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestSignInUp(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	cmpHandler.SignInUp(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestLoginC(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}
	cmpHandler.LoginC(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestCreateAccountC(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}
	cmpHandler.CreateAccountC(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestShowProfile(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	cmpHandler.ShowProfile(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestLogoutC(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	cmpHandler.LogoutC(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestHomeC(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	cmpHandler.HomeC(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

// test application

var ap *ApplicationHandler

func TestApplications(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	applicationHandler.Applications(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestApply(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	res := r.Result()

	applicationHandler.Apply(r, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestApplicationsUpdate(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	applicationHandler.ApplicationsUpdate(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestApplicationDelete(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	applicationHandler.ApplicationDelete(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

// Test post

var p CompanyPostHandler

func TestCompanyPosts(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}
	postHandler.CompanyPosts(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestCompanyPostsNew(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}
	postHandler.CompanyPostsNew(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}
}

func TestCompanyPostsUpdate(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	postHandler.CompanyPostsUpdate(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestCompanyPostsDelete(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	postHandler.CompanyPostsDelete(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

//test request

var re RequestHandler

func TestRequests(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	requestHandler.Requests(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestJoin(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	requestHandler.Join(r, req)

	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestRequestUpdate(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	requestHandler.RequestUpdate(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestRequestDelete(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	requestHandler.RequestDelete(r, req)
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}
