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
	"github.com/amthesonofGod/Notice-Board/request"
	"github.com/amthesonofGod/Notice-Board/user"

	// "github.com/satori/go.uuid"
)

// RequestHandler handles user event join requests
type RequestHandler struct {
	tmpl           *template.Template
	reqSrv         request.RequestService
	postSrv        post.PostService
	userSrv        user.UserService
	sessionService user.SessionService
	csrfSignKey    []byte
}

// NewRequestHandler initializes and returns new RequestHandler
func NewRequestHandler(T *template.Template, RQ request.RequestService, PS post.PostService, US user.UserService, csKey []byte) *RequestHandler {
	return &RequestHandler{tmpl: T, reqSrv: RQ, postSrv: PS, userSrv: US, csrfSignKey:csKey}
}

// Requests handle requests on route /requests
func (rqh *RequestHandler) Requests(w http.ResponseWriter, r *http.Request) {

	handler := UserHandler{loggedInUser: currentUser}

	reqs, errs := rqh.reqSrv.Requests()
	if len(errs) > 0 {
		panic(errs)
	}

	userRequests := []entity.Request{}

	for _, req := range reqs {
		if handler.loggedInUser.ID == req.UserID {
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

	reqs, errs := rqh.reqSrv.Requests()
	if len(errs) > 0 {
		panic(errs)
	}

	userRequests := []entity.Request{}

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

	token, err := rtoken.CSRFToken(rqh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

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

		values := url.Values{}
		values.Add("id", idRaw)
		newRequestForm := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			Post 	 *entity.Post
			CSRF     string
		}{
			Values:   values,
			VErrors:  form.ValidationErrors{},
			Post:	  post,
			CSRF:     token,
		}

		rqh.tmpl.ExecuteTemplate(w, "user_request.layout", newRequestForm)

	}

	if r.Method == http.MethodPost {

		req := &entity.Request{}
		req.FullName = r.FormValue("fullname")
		req.Email = r.FormValue("email")
		req.Phone = r.FormValue("phone")

		handler := UserHandler{loggedInUser: currentUser}

		req.UserID = handler.loggedInUser.ID
		pstID, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			panic(err)
		}

		req.PostID = uint(pstID)

		fmt.Println(pstID)

		// Validate the form contents
		applicationForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		applicationForm.Required("fullname", "email", "phone")
		applicationForm.MatchesPattern("email", form.EmailRX)
		applicationForm.MatchesPattern("phone", form.PhoneRX)
		applicationForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !applicationForm.Valid() {
			rqh.tmpl.ExecuteTemplate(w, "user_request.layout", applicationForm)
			return
		}

		_, errs := rqh.reqSrv.StoreRequest(req)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	}
}

// RequestUpdate handle requests on /user/requests/update
func (rqh *RequestHandler) RequestUpdate(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(rqh.csrfSignKey)
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

		req, errs := rqh.reqSrv.Request(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		values := url.Values{}
		values.Add("id", idRaw)
		values.Add("fullname", req.FullName)
		values.Add("email", req.Email)
		values.Add("phone", req.Phone)
		upRequestForm := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			Request 	 *entity.Request
			CSRF     string
		}{
			Values:   values,
			VErrors:  form.ValidationErrors{},
			Request:	  req,
			CSRF:     token,
		}

		rqh.tmpl.ExecuteTemplate(w, "request_update.layout", upRequestForm)

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		updateRequestForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		updateRequestForm.Required("fullname", "phone", "email")
		updateRequestForm.CSRF = token

		rqs := &entity.Request{}
		reqID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		req, errs := rqh.reqSrv.Request(uint(reqID))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		rqs.FullName = r.FormValue("fullname")
		rqs.Email = r.FormValue("email")
		rqs.Phone = r.FormValue("phone")

		pstID := req.PostID
		usrID := req.UserID

		rqs.PostID = uint(pstID)
		rqs.UserID = uint(usrID)

		fmt.Println(rqs.PostID)
		fmt.Println(rqs.UserID)

		_, errs = rqh.reqSrv.UpdateRequest(rqs)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/requests", http.StatusSeeOther)
	}

}

// RequestDelete handle requests on route /user/requests/delete
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

	http.Redirect(w, r, "/requests", http.StatusSeeOther)
}
