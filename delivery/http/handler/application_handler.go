package handler

import (
	"fmt"
	"github.com/amthesonofGod/Notice-Board/form"
	"github.com/amthesonofGod/Notice-Board/rtoken"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/application"
	"github.com/amthesonofGod/Notice-Board/user"

)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl           *template.Template
	appSrv         application.ApplicationService
	userSrv        user.UserService
	sessionService user.SessionService
	postSrv        post.PostService
	csrfSignKey    []byte
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP application.ApplicationService, US user.UserService, PS post.PostService, csKey []byte) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, userSrv: US, postSrv: PS, csrfSignKey:csKey}
}

// Applications handle requests on route /applications
func (ap *ApplicationHandler) Applications(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(ap.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	userApplications := []entity.Application{}

	handler := UserHandler{loggedInUser: currentUser}

	for _, app := range apps {
		if handler.loggedInUser.ID == app.UserID {
			userApplications = append(userApplications, app)
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

	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Posts   []entity.Post
		Applications []entity.Application
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Posts:   appliedOnPosts,
		Applications: userApplications,
		CSRF:    token,
	}

	// tmplData := struct {
	// 	UserApplications	[]entity.Application
	// 	Post				*entity.Post
	// }{userApplications, post}

	//m := map[string]interface{}{
	//	"UserApplications": userApplications,
	//	"Post":   appliedOnPosts,
	//}
	
	ap.tmpl.ExecuteTemplate(w, "application_list.layout", tmplData)
}

// CompanyReceivedApplications handle requests on route /received/applications
func(ap *ApplicationHandler) CompanyReceivedApplications(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(ap.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	userApplications := []entity.Application{}

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

	fmt.Println(userApplications)

	//m := map[string]interface{}{
	//	"UserApplications": userApplications,
	//	"Post":   appliedOnPosts,

	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Posts   []entity.Post
		Applications []entity.Application
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Posts:   appliedOnPosts,
		Applications: userApplications,
		CSRF:    token,
	}
	
	ap.tmpl.ExecuteTemplate(w, "received_applications.layout", tmplData)
}

//ApplicationDetails handle requests on route /received/applications/details
func(ap *ApplicationHandler) ApplicationDetails(w http.ResponseWriter, r *http.Request) {

	//token, err := rtoken.CSRFToken(ap.csrfSignKey)
	//if err != nil {
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//}

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

		//tmplData := struct {
		//	Values  url.Values
		//	VErrors form.ValidationErrors
		//	Application *entity.Application
		//	CSRF    string
		//}{
		//	Values:  nil,
		//	VErrors: nil,
		//	Application: app,
		//	CSRF:    token,
		//}

		ap.tmpl.ExecuteTemplate(w, "received_applications_detail.layout", app)
	} else {
		
		ap.tmpl.ExecuteTemplate(w, "received_applications_detail.layout", nil)
	}
}

// Apply hanlde requests on route /job/apply
func (ap *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(ap.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

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


		values := url.Values{}
		values.Add("id", idRaw)
		applicationForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			Post    *entity.Post
			CSRF    string
		}{
			Values:  values,
			VErrors: nil,
			Post:   post,
			CSRF:    token,
		}

		ap.tmpl.ExecuteTemplate(w, "user_application.layout", applicationForm)

	}

	if r.Method == http.MethodPost {

		fmt.Println("post invoked")
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

		handler := UserHandler{loggedInUser: currentUser}

		app.UserID = handler.loggedInUser.ID
		pstID, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			panic(err)
		}

		app.PostID = uint(pstID)
		
		fmt.Println(pstID)

		
		fmt.Println(pstID)
		fmt.Println(r.FormValue("id"))

		// Validate the form contents
		applicationForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		// applicationForm.Required("fullname", "email", "letter", "resume")
		// applicationForm.MatchesPattern("useremail", form.EmailRX)
		// applicationForm.MatchesPattern("phone", form.PhoneRX)
		applicationForm.MinLength("letter", 20)
		applicationForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !applicationForm.Valid() {
			ap.tmpl.ExecuteTemplate(w, "user_application.layout", applicationForm)
			return
		}

		_, errs := ap.appSrv.StoreApplication(app)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/applications", http.StatusSeeOther)
	}
}

// ApplicationUpdate handle requests on /applications/update
func (ap *ApplicationHandler) ApplicationUpdate(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(ap.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		fmt.Println(uint(id))

		if err != nil {
			panic(err)
		}

		app, errs := ap.appSrv.Application(uint(id))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		values := url.Values{}
		values.Add("id", idRaw)
		values.Add("fullname", app.FullName)
		values.Add("email", app.Email)
		values.Add("phone", app.Phone)
		values.Add("letter", app.Letter)
		values.Add("oldresume", app.Resume)
		upApplicationForm := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			Application 	 *entity.Application
			CSRF     string
		}{
			Values:   values,
			VErrors:  form.ValidationErrors{},
			Application:	  app,
			CSRF:     token,
		}

		ap.tmpl.ExecuteTemplate(w, "application_update.layout", upApplicationForm)

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		updateApplicationForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		updateApplicationForm.Required("fullname", "letter")
		updateApplicationForm.MinLength("letter", 20)
		updateApplicationForm.CSRF = token

		appID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		app, errs := ap.appSrv.Application(uint(appID))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		appc := &entity.Application{}
		appc.ID = uint(appID)

		appc.UserID = app.UserID
		appc.PostID = app.PostID

		fmt.Println(appc.UserID)
		fmt.Println(appc.PostID)

		appc.FullName = r.FormValue("fullname")
		appc.Email = r.FormValue("email")
		appc.Phone = r.FormValue("phone")
		appc.Letter = r.FormValue("letter")

		appc.Resume = r.FormValue("oldresume")

		mf, fh, err := r.FormFile("resume")

		if err == nil {
			appc.Resume = fh.Filename

			writeFile(&mf, appc.Resume)
		}

		if mf != nil {
			defer mf.Close()
		}


		_, errs = ap.appSrv.UpdateApplication(appc)

		if len(errs) > 0 {
			panic(errs)
		}

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
