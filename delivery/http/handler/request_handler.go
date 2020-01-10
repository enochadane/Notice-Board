package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/request"
	// "github.com/satori/go.uuid"
)

// RequestHandler handles user event join requests
type RequestHandler struct {
	tmpl    *template.Template
	reqSrv  request.RequestService
	postSrv post.PostService
}

// NewRequestHandler initializes and returns new RequestHandler
func NewRequestHandler(T *template.Template, RQ request.RequestService, PS post.PostService) *RequestHandler {
	return &RequestHandler{tmpl: T, reqSrv: RQ, postSrv: PS}
}

// Requests handle requests on route /requests
func (rqh *RequestHandler) Requests(w http.ResponseWriter, r *http.Request) {

	reqs, errs := rqh.reqSrv.Requests()
	if len(errs) > 0 {
		panic(errs)
	}

	rqh.tmpl.ExecuteTemplate(w, "request_list.layout", reqs)
}

// Join hanlde requests on route /event/join
func (rqh *RequestHandler) Join(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		req := &entity.Request{}
		req.FullName = r.FormValue("fullname")
		req.Email = r.FormValue("email")
		req.Phone = r.FormValue("phone")

		_, errs := rqh.reqSrv.StoreRequest(req)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	} else {

		rqh.tmpl.ExecuteTemplate(w, "user_request.layout", nil)

	}
}

// RequestUpdate handle requests on /requests/update
func (rqh *RequestHandler) RequestUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		fmt.Println(uint(id))

		if err != nil {
			panic(err)
		}

		req, errs := rqh.reqSrv.Request(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		rqh.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", req)

	} else if r.Method == http.MethodPost {

		rqs := &entity.Request{}
		id, _ := strconv.Atoi(r.FormValue("id"))
		rqs.ID = uint(id)
		rqs.FullName = r.FormValue("fullname")
		rqs.Email = r.FormValue("email")
		rqs.Phone = r.FormValue("phone")
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

// RequestDelete handle requests on route /requests/delete
func (rqh *RequestHandler) RequestDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		_, errs := rqh.reqSrv.DeleteRequest(uint(id))

		if len(errs) > 0 {
			panic(err)
		}

	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}
