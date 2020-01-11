package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	// "io/ioutil"
	// "log"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/application"
	"github.com/amthesonofGod/Notice-Board/model"

)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl	*template.Template
	appSrv	application.ApplicationService
	userSrv model.UserService
	postSrv post.PostService
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP application.ApplicationService, US model.UserService, PS post.PostService) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, userSrv: US, postSrv: PS}
}

// Applications handle requests on route /applications
func(ap *ApplicationHandler) Applications(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")

	s, serr := ap.userSrv.Session(cookie.Value)

	if len(serr) > 0 {
		panic(serr)
	}

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	userApplications := []entity.Application{}

	// var pstid uint
	for _, app := range apps {
		if s.UserID == app.UserID {
			userApplications = append(userApplications, app)
			// pstid = app.PostID
		}
	}

	posts, perr := ap.postSrv.Posts()

	if len(perr) > 0 {
		panic(perr)
	}

	appliedOnPosts := []entity.Post{}

	for _, post := range posts {
		for _, myapps := range userApplications {
			if myapps.PostID == post.ID && post.Category == "Job" {
				appliedOnPosts = append(appliedOnPosts, post)
			}
		}
	}

	// tmplData := struct {
	// 	UserApplications	[]entity.Application
	// 	Post				*entity.Post
	// }{userApplications, post}

	m := map[string]interface{}{
		"UserApplications": userApplications,
		"Post":   appliedOnPosts,
	}
	
	ap.tmpl.ExecuteTemplate(w, "application_list.layout", m)
}

// CompanyReceivedApplications handle requests on route /received/applications
func(ap *ApplicationHandler) CompanyReceivedApplications(w http.ResponseWriter, r *http.Request) {

	// cookie, _ := r.Cookie("session")

	// s, serr := ap.userSrv.Session(cookie.Value)

	// if len(serr) > 0 {
	// 	panic(serr)
	// }

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	userApplications := []entity.Application{}

	// // var pstid uint
	// for _, app := range apps {
	// 	if s.UserID == app.UserID {
	// 		userApplications = append(userApplications, app)
	// 		// pstid = app.PostID
	// 	}
	// }

	posts, perr := ap.postSrv.Posts()

	if len(perr) > 0 {
		panic(perr)
	}

	appliedOnPosts := []entity.Post{}

	for _, post := range posts {
		for _, app := range apps {
			if app.PostID == post.ID && post.Category == "Job" {
				userApplications = append(userApplications, app)
				appliedOnPosts = append(appliedOnPosts, post)
			}
		}
	}

	m := map[string]interface{}{
		"UserApplications": userApplications,
		"Post":   appliedOnPosts,
	}
	
	ap.tmpl.ExecuteTemplate(w, "received_applications.layout", m)
}

//ApplicationDetails handle requests on route /received/applications/details
func(ap *ApplicationHandler) ApplicationDetails(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == http.MethodGet {
		
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		fmt.Println(id)

		app, errs := ap.appSrv.Application(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		fmt.Println(app)

		ap.tmpl.ExecuteTemplate(w, "received_applications_detail.layout", app)
	} else {
		
		ap.tmpl.ExecuteTemplate(w, "received_applications_detail.layout", nil)
	}
}

// Apply hanlde requests on route /job/apply
func (ap *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		// fmt.Println(id)

		if err != nil {
			panic(err)
		}

		post, errs := ap.postSrv.Post(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		ap.tmpl.ExecuteTemplate(w, "user_application.layout", post)

	} 

	if r.Method == http.MethodPost {

		app := &entity.Application{}
		app.FullName = r.FormValue("fullname")
		app.Email = r.FormValue("email")
		app.Phone = r.FormValue("phone")
		app.Letter = r.FormValue("letter")

		mf, fh, err := r.FormFile("resume")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		app.Resume = fh.Filename

		writeFile(&mf, fh.Filename)

		cookie, _ := r.Cookie("session")
		s, errs := ap.userSrv.Session(cookie.Value)

		app.UserID = s.UserID
		pstID, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			panic(err)
		}

		app.PostID = uint(pstID)

		// reqBody, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 		log.Fatal(err)
		// }

		// fmt.Printf("%s\n", reqBody)
		
		fmt.Println(pstID)
		fmt.Println(r.FormValue("id"))

		_, errs = ap.appSrv.StoreApplication(app)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/applications", http.StatusSeeOther)
	}
}

// ApplicationsUpdate handle requests on /user/applications/update
func (ap *ApplicationHandler) ApplicationUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		fmt.Println(uint(id))

		if err != nil {
			panic(err)
		}

		app, errs := ap.appSrv.Application(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		ap.tmpl.ExecuteTemplate(w, "application_update.layout", app)

	} else if r.Method == http.MethodPost {

		appc := &entity.Application{}
		id, _ := strconv.Atoi(r.FormValue("id"))
		appc.ID = uint(id)

		userid, _ := strconv.Atoi(r.FormValue("userid"))
		postid, _ := strconv.Atoi(r.FormValue("postid"))

		appc.UserID = uint(userid)
		appc.PostID = uint(postid)

		fmt.Println(appc.UserID)
		fmt.Println(appc.PostID)

		appc.FullName = r.FormValue("fullname")
		appc.Email = r.FormValue("email")
		appc.Phone = r.FormValue("phone")
		appc.Letter = r.FormValue("letter")

		appc.Resume = r.FormValue("oldresume")

		mf, fh, err := r.FormFile("resume")

		if err != nil {
			panic(err)
		}

		defer mf.Close()

		appc.Resume = fh.Filename

		writeFile(&mf, appc.Resume)

		_, errs := ap.appSrv.UpdateApplication(appc)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/applications", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/applications", http.StatusSeeOther)
	}

}

// ApplicationDelete handle requests on route /user/applications/delete
func (ap *ApplicationHandler) ApplicationDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		_, errs := ap.appSrv.DeleteApplication(uint(id))

		if len(errs) > 0 {
			panic(err)
		}

	}

	http.Redirect(w, r, "/applications", http.StatusSeeOther)
}

