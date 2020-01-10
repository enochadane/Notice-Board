package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/application"

	// "github.com/satori/go.uuid"
)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl	*template.Template
	appSrv	application.ApplicationService
	postSrv post.PostService
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP application.ApplicationService, PS post.PostService) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, postSrv: PS}
}

// Applications handle requests on route /applications
func(ap *ApplicationHandler) Applications(w http.ResponseWriter, r *http.Request) {

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	ap.tmpl.ExecuteTemplate(w, "application_list.layout", apps)
}

// Apply hanlde requests on route /job/apply
func (ap *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {

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

		_, errs := ap.appSrv.StoreApplication(app)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/applications", http.StatusSeeOther)
	} else {
		
		ap.tmpl.ExecuteTemplate(w, "user_application.layout", nil)
	}
}

// ApplicationsUpdate handle requests on /applications/update
func (ap *ApplicationHandler) ApplicationsUpdate(w http.ResponseWriter, r *http.Request) {

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

		ap.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", app)

	} else if r.Method == http.MethodPost {

		appc := &entity.Application{}
		id, _ := strconv.Atoi(r.FormValue("id"))
		appc.ID = uint(id)
		appc.FullName = r.FormValue("fullname")
		appc.Email = r.FormValue("email")
		appc.Phone = r.FormValue("phone")
		appc.Letter = r.FormValue("letter")
		// appc.Resume = r.FormValue("resume")

		// mf, _, err := r.FormFile("resume")

		// if err != nil {
		// 	panic(err)
		// }

		// defer mf.Close()

		// writeFile(&mf, ctg.Image)

		// _, errs := ach.categorySrv.UpdateCategory(ctg)

		// if len(errs) > 0 {
		// 	panic(errs)
		// }

		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
	}

}

// ApplicationDelete handle requests on route /applications/delete
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

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

