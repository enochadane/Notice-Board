package handler

import (
	"fmt"
	"github.com/amthesonofGod/Notice-Board/form"
	"github.com/amthesonofGod/Notice-Board/rtoken"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	// "encoding/json"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
)

// CompanyPostHandler handles post handler admin requests
type CompanyPostHandler struct {
	tmpl       		*template.Template
	postSrv   		post.PostService
	companySrv 		company.CompanyService
	sessionService	company.SessionServiceCamp	
	campSess		*entity.CompanySession
	csrfSignKey		[]byte
}

// NewCompanyPostHandler initializes and returns new CompanyPostHandler
func NewCompanyPostHandler(T *template.Template, PS post.PostService, CP company.CompanyService, csKey []byte) *CompanyPostHandler {
	return &CompanyPostHandler{tmpl: T, postSrv: PS, companySrv: CP, csrfSignKey:csKey}
}

// CompanyPosts handle requests on route /admin/posts
func (cph *CompanyPostHandler) CompanyPosts(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(cph.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	handler := CompanyHandler{loggedInUserCamp: currentCompUser}

	authorizedPost := []entity.Post{}

	posts, errs := cph.postSrv.Posts()
	if len(errs) > 0 {
		panic(errs)
	}
	for _, post := range posts {
		if handler.loggedInUserCamp.ID == post.CompanyID {
			authorizedPost = append(authorizedPost, post)
		}
	}

	fmt.Println("Current Post")
	fmt.Println(authorizedPost)

	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Posts   []entity.Post
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Posts:   authorizedPost,
		CSRF:    token,
	}

	cph.tmpl.ExecuteTemplate(w, "cmp_post.layout", tmplData)
}

// CompanyPostsNew hanlde requests on route /admin/new-post
func (cph *CompanyPostHandler) CompanyPostsNew(w http.ResponseWriter, r *http.Request) {

	fmt.Println("companypostsnew function invoked! ")

	token, err := rtoken.CSRFToken(cph.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		newPostForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		cph.tmpl.ExecuteTemplate(w, "post-job.layout", newPostForm)
	}

	if r.Method == http.MethodPost {

		handler := CompanyHandler{loggedInUserCamp: currentCompUser}

		compID := handler.loggedInUserCamp.ID

		cmp, cerr := cph.companySrv.Company(compID)

		fmt.Println(compID)

		if len(cerr) > 0 {
			fmt.Println("i am the error")
			panic(cerr)
		}

		fmt.Println(cmp.Name)

		post := &entity.Post{}
		post.CompanyID = compID
		post.Owner = cmp.Name
		post.Title = r.FormValue("title")
		post.Description = r.FormValue("description")
		post.Category = r.Form.Get("category")
		//post.Image = ""

		fmt.Println(post.Category)

		mf, fh, err := r.FormFile("postimg")
		if err == nil {
			post.Image = fh.Filename
			writeFile(&mf, fh.Filename)
		}
		if mf != nil {
			defer mf.Close()
		}



		// Validate the form contents
		newPostForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newPostForm.Required("title", "description", "category")
		newPostForm.MinLength("description", 10)
		newPostForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !newPostForm.Valid() {
			cph.tmpl.ExecuteTemplate(w, "post-job.layout", newPostForm)
			return
		}

		_, errs := cph.postSrv.StorePost(post)
		// cph.postSrv.StorePost(post)

		if len(errs) > 0 {
			panic(errs)
		}
		// fmt.Println(entity.Company.ID)
		fmt.Println(post)
		fmt.Println("post added to db")

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)

	}

}

// CompanyPostUpdate handle requests on /admin/posts/update
func (cph *CompanyPostHandler) CompanyPostUpdate(w http.ResponseWriter, r *http.Request) {

	token, err := rtoken.CSRFToken(cph.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if r.Method == http.MethodGet {

		fmt.Println("get method invoked")

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		post, errs := cph.postSrv.Post(uint(id))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		values := url.Values{}
		values.Add("id", idRaw)
		//values.Add("owner", post.Owner)
		//values.Add("companyid", string(post.CompanyID))
		values.Add("title", post.Title)
		values.Add("description", post.Description)
		values.Add("image", post.Image)
		upPostForm := struct {
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
		cph.tmpl.ExecuteTemplate(w, "post_update.layout", upPostForm)

	} else if r.Method == http.MethodPost {
		fmt.Println("post method invoked")
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		updatePostForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		updatePostForm.Required("title", "description")
		updatePostForm.MinLength("description", 10)
		updatePostForm.CSRF = token

		pstID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		post, errs := cph.postSrv.Post(uint(pstID))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		pst := &entity.Post{}
		pst.ID = uint(pstID)
		pst.Title = r.FormValue("title")
		pst.Description = r.FormValue("description")
		pst.Image = r.FormValue("image")
		pst.Category = post.Category
		pst.CompanyID = post.CompanyID
		pst.Owner = post.Owner

		fmt.Println(pst.ID)
		fmt.Println(pst.CompanyID)
		fmt.Println(pst.Owner)

		mf, fh, err := r.FormFile("postimg")

		if err == nil {
			pst.Image = fh.Filename
			writeFile(&mf, pst.Image)
		}

		if mf != nil {
			defer mf.Close()
		}

		_, errs = cph.postSrv.UpdatePost(pst)

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return

	} else {
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	}

}

// CompanyPostDelete handle requests on route /cmp/posts/delete
func (cph *CompanyPostHandler) CompanyPostDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		_, errs := cph.postSrv.DeletePost(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
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
