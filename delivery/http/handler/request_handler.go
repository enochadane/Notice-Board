package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/request"
	"github.com/amthesonofGod/Notice-Board/model"

	// "github.com/satori/go.uuid"
)

// RequestHandler handles user event join requests
type RequestHandler struct {
	tmpl		*template.Template
	reqSrv		request.RequestService
	postSrv 	post.PostService
	userSrv 	model.UserService
}

// NewRequestHandler initializes and returns new RequestHandler
func NewRequestHandler(T *template.Template, RQ request.RequestService, PS post.PostService, US model.UserService) *RequestHandler {
	return &RequestHandler{tmpl: T, reqSrv: RQ, postSrv: PS, userSrv: US}
}

// Requests handle requests on route /requests
func(rqh *RequestHandler) Requests(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")

	s, errs := rqh.userSrv.Session(cookie.Value)
	if len(errs) > 0 {
		panic((errs))
	}

	reqs, errs := rqh.reqSrv.Requests()
	if len(errs) > 0 {
		panic(errs)
	}

	userRequests := []entity.Request{}

	for _, req := range reqs {
		if s.UserID == req.UserID {
			userRequests = append(userRequests, req)
		}
	}

	posts, errs := rqh.postSrv.Posts()
	if len(errs) > 0 {
		panic(errs)
	}

	requestedEvents := []entity.Post{}

	for _, post := range posts {
		for _, myreqs := range userRequests {
			if myreqs.PostID == post.ID && post.Category == "Event" {
				requestedEvents = append(requestedEvents, post)
			}
		}
	}

	m := map[string]interface{}{
		"UserRequests": userRequests,
		"Events": requestedEvents,
	}

	rqh.tmpl.ExecuteTemplate(w, "request_list.layout", m)
}

// CompanyReceivedRequests handle requests on route /received/requests
func(rqh *RequestHandler) CompanyReceivedRequests(w http.ResponseWriter, r *http.Request) {

	// cookie, _ := r.Cookie("session")

	// s, errs := rqh.companySrv.Session(cookie.Value)
	// if len(errs) > 0 {
	// 	panic((errs))
	// }

	reqs, errs := rqh.reqSrv.Requests()
	if len(errs) > 0 {
		panic(errs)
	}

	userRequests := []entity.Request{}

	// for _, req := range reqs {
	// 	if s.UserID == req.UserID {
	// 		userRequests = append(userRequests, req)
	// 	}
	// }

	posts, errs := rqh.postSrv.Posts()
	if len(errs) > 0 {
		panic(errs)
	}

	requestedEvents := []entity.Post{}

	for _, post := range posts {
		for _, req := range reqs {
			if req.PostID == post.ID && post.Category == "Event" {
				userRequests = append(userRequests, req)
				requestedEvents = append(requestedEvents, post)
			}
		}
	}

	m := map[string]interface{}{
		"UserRequests": userRequests,
		"Events": requestedEvents,
	}

	rqh.tmpl.ExecuteTemplate(w, "received_requests.layout", m)
}

// Join hanlde requests on route /event/join
func (rqh *RequestHandler) Join(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		post, errs := rqh.postSrv.Post(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		rqh.tmpl.ExecuteTemplate(w, "user_request.layout", post)

	}

	if r.Method == http.MethodPost {

		req := &entity.Request{}
		req.FullName = r.FormValue("fullname")
		req.Email = r.FormValue("email")
		req.Phone = r.FormValue("phone")

		cookie, _ := r.Cookie("session")
		s, errs := rqh.userSrv.Session(cookie.Value)

		req.UserID = s.UserID
		pstID, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			panic(err)
		}

		req.PostID = uint(pstID)

		_, errs = rqh.reqSrv.StoreRequest(req)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	} else {

		rqh.tmpl.ExecuteTemplate(w, "user_request.layout", nil)

	}
}

// RequestUpdate handle requests on /user/requests/update
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

		rqh.tmpl.ExecuteTemplate(w, "request_update.layout", req)

	} else if r.Method == http.MethodPost {

		rqs := &entity.Request{}
		id, _ := strconv.Atoi(r.FormValue("id"))
		rqs.ID = uint(id)
		rqs.FullName = r.FormValue("fullname")
		rqs.Email = r.FormValue("email")
		rqs.Phone = r.FormValue("phone")

		pstID, err := strconv.Atoi(r.FormValue("postid"))
		if err != nil {
			panic(err)
		}
		usrID, err := strconv.Atoi(r.FormValue("userid"))
		if err != nil {
			panic(err)
		}

		rqs.PostID = uint(pstID)
		rqs.UserID = uint(usrID)

		fmt.Println(rqs.PostID)
		fmt.Println(rqs.UserID)

		_, errs := rqh.reqSrv.UpdateRequest(rqs)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/requests", http.StatusSeeOther)
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
