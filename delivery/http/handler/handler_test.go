package handler

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"bytes"

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
	use      = "postgres"
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
		host, port, use, password, dbname)

	dbconn, err := gorm.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	defer dbconn.Close()

	sessCamp := configSessCamp()
	sess := configSess()

	userSessionRepo := repository.NewSessionGormRepo(dbconn)
	userSessionsrv := service.NewSessionService(userSessionRepo)

	companySessionRepo := repositoryCamp.NewSessionGormRepoCamp(dbconn)
	companySessionSrv := serviceCamp.NewSessionServiceCamp(companySessionRepo)

	companyRepo := repositoryCamp.NewCompanyGormRepoMock(nil)
	companySrv := serviceCamp.NewCompanyService(companyRepo)

	postRepo := postRepos.NewPostGormRepoMock(nil)
	postSrv := postServ.NewPostService(postRepo)
	postHandler = NewCompanyPostHandler(tmpl, postSrv, companySrv)

	cmpHandler = NewCompanyHandler(tmpl, companySrv, postSrv, companySessionSrv, sessCamp)

	// cmpHandler.loggedInUserCamp = &entity.CompanyMock

	currentCompUser = &entity.CompanyMock

	userRepo := repository.NewUserGormRepo(dbconn)
	userSrv := service.NewUserService(userRepo)
	csrfSignKey := []byte(rtoken.GenerateRandomID(32))
	//T *template.Template, US user.UserService, PS post.PostService, sessServ user.SessionService, usrSess *entity.UserSession, csKey []byte
	usrHandler = NewUserHandler(tmpl, userSrv, postSrv, userSessionsrv, sess, csrfSignKey)

	applicationRepo := appRepos.NewApplicationGormRepo(dbconn)
	applicationSrv := appServ.NewApplicationService(applicationRepo)
	applicationHandler = NewApplicationHandler(tmpl, applicationSrv, userSrv, postSrv)

	requestRepo := reqRepos.NewRequestGormRepo(dbconn)
	requestSrv := reqServ.NewRequestService(requestRepo)
	requestHandler = NewRequestHandler(tmpl, requestSrv, postSrv, userSrv)

	//(T *template.Template, CS company.CompanyService, PS post.PostService, sessServ company.SessionServiceCamp, campSess *entity.CompanySession)

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

// test application

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

func TestApply(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/job/apply", nil)

	if er != nil {
		t.Fatal(er)
	}

	res := r.Result()

	applicationHandler.Apply(r, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

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

	read, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s", string(read))

}

func TestCompanyPosts(t *testing.T) {

	// r := httptest.NewRecorder()

	// req, er := http.NewRequest("GET", "/admin/cmp-posts", nil)

	// if er != nil {
	// 	t.Fatal(er)
	// }
	// postHandler.CompanyPosts(r, req)

	// res := r.Result()

	// if res.StatusCode != http.StatusOK {
	// 	t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	// }

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/cmp-posts", postHandler.CompanyPosts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	url := ts.URL

	resp, err := tc.Get(url + "/admin/cmp-posts")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("Mock Category 01")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestCompanyPostsNew(t *testing.T) {
	// r := httptest.NewRecorder()
	// req, er := http.NewRequest("GET", "/", nil)
	// if er != nil {
	// 	t.Fatal(er)
	// }
	// postHandler.CompanyPostsNew(r, req)
	// res := r.Result()
	// if res.StatusCode != http.StatusOK {
	// 	t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	// }

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/posts/new", postHandler.CompanyPosts)

	ser := httptest.NewTLSServer(mux)
	tc := ser.Client()

	defer ser.Close()
	sURL := ser.URL

	form := url.Values{}

	form.Add("title", entity.PostMock.Title)
	form.Add("description", entity.PostMock.Description)
	form.Add("postimg", entity.PostMock.Image)
	form.Add("category", entity.PostMock.Category)

	res, err := tc.PostForm(sURL+"/admin/posts/new", form)

	if err != nil {
		t.Fatal("error")
	}

	if res.StatusCode != http.StatusOK {

		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()
	// body, er := ioutil.ReadAll(res.Body)

	// if er != nil {
	// 	t.Fatal(err)
	// }

	// if !(bytes.Contains(body, []byte("new post"))) {
	// 	t.Errorf("wanted %q", string(body))
	// }

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

func TestApplicationsUpdate(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	applicationHandler.ApplicationUpdate(r, req)
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

//Test for User

//var u *UserHandler

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

	//expected
	res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

}

func TestApplications(t *testing.T) {
	// r := httptest.NewRecorder()
	// req, er := http.NewRequest("GET", "/", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/applications", applicationHandler.Applications)
	tes := httptest.NewTLSServer(mux)

	url := tes.URL
	cl := tes.Client()
	defer tes.Close()

	res, err := cl.Get(url + "/applications")

	if err != nil {
		t.Fatal(err)
	}

	// applicationHandler.Applications(r, req)
	// res := r.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wanted %d got %d", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()

	body, er := ioutil.ReadAll(res.Body)

	if er != nil {
		t.Fatal(er)
	}

	if !(bytes.Contains(body, []byte("applications"))) {
		t.Errorf("wanted %q", body)
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

// Test post

var p CompanyPostHandler

func TestCompanyPostsUpdate(t *testing.T) {

	r := httptest.NewRecorder()

	req, er := http.NewRequest("GET", "/", nil)

	if er != nil {
		t.Fatal(er)
	}

	postHandler.CompanyPostUpdate(r, req)
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

	postHandler.CompanyPostDelete(r, req)

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
