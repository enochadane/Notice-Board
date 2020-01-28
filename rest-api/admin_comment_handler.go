package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/rest-api/applicationapi"
	"github.com/amthesonofGod/Notice-Board/user"
	"github.com/julienschmidt/httprouter"
)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl           *template.Template
	appSrv         applicationapi.ApplicationService
	userSrv        user.UserService
	sessionService user.SessionService
	postSrv        post.PostService
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP applicationapi.ApplicationService, US user.UserService, PS post.PostService) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, userSrv: US, postSrv: PS}
}

// Applications handle requests on route /applications
func (ap *ApplicationHandler) Applications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	apps, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	output, err := json.MarshalIndent(apps, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

//CompanyReceivedApplications handle requests on route /received/applications
func (ap *ApplicationHandler) CompanyReceivedApplications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a, errs := ap.appSrv.Applications()
	if len(errs) > 0 {
		panic(errs)
	}

	userApplications := []entity.Application{}

	posts, er := ap.postSrv.Posts()

	if len(er) > 0 {
		panic(er)
	}

	appliedOnPosts := []entity.Post{}

	for _, post := range posts {
		for _, apa := range a {
			if apa.PostID == post.ID && post.Category == "Job" {
				userApplications = append(userApplications, apa)
				appliedOnPosts = append(appliedOnPosts, post)
			}
		}
	}

	// m := map[string]interface{}{
	// 	"UserApplications": userApplications,
	// 	"Post":             appliedOnPosts,
	// }

	output, err := json.MarshalIndent(userApplications, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

//ApplicationDetails handle requests on route /received/applications/details
func (ap *ApplicationHandler) ApplicationDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		panic(err)
	}

	fmt.Println(id)

	app, errs := ap.appSrv.Application(uint(id))

	if len(errs) > 0 {
		panic(errs)
	}

	output, er := json.MarshalIndent(app, "", "\t\t")

	if er != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

//Apply hanlde requests on route /job/apply
func (ap *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// if r.Method == http.MethodGet {

	// 	// idRaw := r.URL.Query().Get("id")
	// 	// id, err := strconv.Atoi(idRaw)

	// 	id, err := strconv.Atoi(ps.ByName("id"))

	// 	// fmt.Println(id)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	post, errs := ap.postSrv.Post(uint(id))

	// 	if len(errs) > 0 {
	// 		panic(errs)
	// 	}

	// }

	app := entity.Application{}
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

	// handle := handler.UserHandler{loggedInUser: handler.UserHandler.currentUser}

	// app.UserID = handle.loggedInUser.ID
	pstID, err := strconv.Atoi(r.FormValue("id"))

	if err != nil {
		panic(err)
	}

	app.PostID = uint(pstID)

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	err = json.Unmarshal(body, app)

	fmt.Println(pstID)

	fmt.Println(pstID)
	fmt.Println(r.FormValue("id"))

	_, errs := ap.appSrv.StoreApplication(&app)

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/admin/comments/%d", app.ID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return

}

// ApplicationUpdate handle requests on /user/applications/update
func (ap *ApplicationHandler) ApplicationUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	application, errs := ap.appSrv.Application(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	l := r.ContentLength

	body := make([]byte, l)

	r.Body.Read(body)

	json.Unmarshal(body, &application)

	application, errs = ap.appSrv.UpdateApplication(application)

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(application, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

	// if r.Method == http.MethodGet {

	// 	// idRaw := r.URL.Query().Get("id")
	// 	// id, err := strconv.Atoi(idRaw)

	// 	id, err := strconv.Atoi(ps.ByName("id"))

	// 	fmt.Println(uint(id))

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	app, errs := ap.appSrv.Application(uint(id))

	// 	if len(errs) > 0 {
	// 		panic(errs)
	// 	}

	// 	ap.tmpl.ExecuteTemplate(w, "application_update.layout", app)

	// } else if r.Method == http.MethodPost {

	// 	appc := &entity.Application{}
	// 	id, _ := strconv.Atoi(r.FormValue("id"))
	// 	appc.ID = uint(id)

	// 	userid, _ := strconv.Atoi(r.FormValue("userid"))
	// 	postid, _ := strconv.Atoi(r.FormValue("postid"))

	// 	appc.UserID = uint(userid)
	// 	appc.PostID = uint(postid)

	// 	fmt.Println(appc.UserID)
	// 	fmt.Println(appc.PostID)

	// 	appc.FullName = r.FormValue("fullname")
	// 	appc.Email = r.FormValue("email")
	// 	appc.Phone = r.FormValue("phone")
	// 	appc.Letter = r.FormValue("letter")

	// 	appc.Resume = r.FormValue("oldresume")

	// 	mf, fh, err := r.FormFile("resume")

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	defer mf.Close()

	// 	appc.Resume = fh.Filename

	// 	writeFile(&mf, appc.Resume)

	// 	_, errs := ap.appSrv.UpdateApplication(appc)

	// 	if len(errs) > 0 {
	// 		panic(errs)
	// 	}

	// }

}

// ApplicationDelete handle requests on route /user/applications/delete
func (ap *ApplicationHandler) ApplicationDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, errs := ap.appSrv.DeleteApplication(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return

}

func writeFile(mf *multipart.File, fname string) {

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "../../", "ui", "assets", "img", fname)
	image, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	io.Copy(image, *mf)
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/betsegawlemma/restaurant-rest/comment"
// 	"github.com/betsegawlemma/restaurant-rest/entity"
// 	"github.com/julienschmidt/httprouter"
// )

// // AdminCommentHandler handles comment related http requests
// type AdminCommentHandler struct {
// 	commentService comment.CommentService
// }

// // NewAdminCommentHandler returns new AdminCommentHandler object
// func NewAdminCommentHandler(cmntService comment.CommentService) *AdminCommentHandler {
// 	return &AdminCommentHandler{commentService: cmntService}
// }

// // GetComments handles GET /v1/admin/comments request
// func (ach *AdminCommentHandler) GetComments(w http.ResponseWriter,
// 	r *http.Request, _ httprouter.Params) {

// 	comments, errs := ach.commentService.Comments()

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	output, err := json.MarshalIndent(comments, "", "\t\t")

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(output)
// 	return

// }

// // GetSingleComment handles GET /v1/admin/comments/:id request
// func (ach *AdminCommentHandler) GetSingleComment(w http.ResponseWriter,
// 	r *http.Request, ps httprouter.Params) {

// 	id, err := strconv.Atoi(ps.ByName("id"))

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	comment, errs := ach.commentService.Comment(uint(id))

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	output, err := json.MarshalIndent(comment, "", "\t\t")

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(output)
// 	return
// }

// // PostComment handles POST /v1/admin/comments request
// func (ach *AdminCommentHandler) PostComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 	l := r.ContentLength
// 	body := make([]byte, l)
// 	r.Body.Read(body)
// 	comment := &entity.Comment{}

// 	err := json.Unmarshal(body, comment)

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	comment, errs := ach.commentService.StoreComment(comment)

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	p := fmt.Sprintf("/v1/admin/comments/%d", comment.ID)
// 	w.Header().Set("Location", p)
// 	w.WriteHeader(http.StatusCreated)
// 	return
// }

// // PutComment handles PUT /v1/admin/comments/:id request
// func (ach *AdminCommentHandler) PutComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 	id, err := strconv.Atoi(ps.ByName("id"))
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	comment, errs := ach.commentService.Comment(uint(id))

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	l := r.ContentLength

// 	body := make([]byte, l)

// 	r.Body.Read(body)

// 	json.Unmarshal(body, &comment)

// 	comment, errs = ach.commentService.UpdateComment(comment)

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	output, err := json.MarshalIndent(comment, "", "\t\t")

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(output)
// 	return
// }

// // DeleteComment handles DELETE /v1/admin/comments/:id request
// func (ach *AdminCommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 	id, err := strconv.Atoi(ps.ByName("id"))

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	_, errs := ach.commentService.DeleteComment(uint(id))

// 	if len(errs) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNoContent)
// 	return
// }
