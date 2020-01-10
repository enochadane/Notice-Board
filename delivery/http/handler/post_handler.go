package handler

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
)

// CompanyPostHandler handles post handler admin requests
type CompanyPostHandler struct {
	tmpl       *template.Template
	postSrv    post.PostService
	companySrv company.CompanyService
}

// NewCompanyPostHandler initializes and returns new CompanyPostHandler
func NewCompanyPostHandler(T *template.Template, PS post.PostService, CP company.CompanyService) *CompanyPostHandler {
	return &CompanyPostHandler{tmpl: T, postSrv: PS, companySrv: CP}
}

// CompanyPosts handle requests on route /admin/cmp-posts
func (cph *CompanyPostHandler) CompanyPosts(w http.ResponseWriter, r *http.Request) {

	var cookie, cerr = r.Cookie("session")
	if cerr == nil {
		cookievalue := cookie.Value
		fmt.Println(cookievalue)
	}

	s, serr := cph.companySrv.Session(cookie.Value)
	if len(serr) > 0 {
		panic(serr)
	}

	authorizedPost := []entity.Post{}

	posts, errs := cph.postSrv.Posts()
	if len(errs) > 0 {
		panic(errs)
	}
	for _, post := range posts {
		if s.CompanyID == post.CompanyID {
			authorizedPost = append(authorizedPost, post)
		}
	}

	fmt.Println("All posts")
	fmt.Println(posts)

	fmt.Println("Current Post")
	fmt.Println(authorizedPost)
	cph.tmpl.ExecuteTemplate(w, "cmp_post.layout", authorizedPost)
}

// CompanyPostsNew hanlde requests on route /admin/posts/new
func (cph *CompanyPostHandler) CompanyPostsNew(w http.ResponseWriter, r *http.Request) {

	fmt.Println("companypostsnew function invoked! ")

	if r.Method == http.MethodPost {

		fmt.Println("post method verified! ")

		var cookie, err = r.Cookie("session")
		if err == nil {
			cookievalue := cookie.Value
			fmt.Println(cookievalue)
		}

		s, serr := cph.companySrv.Session(cookie.Value)

		if len(serr) > 0 {
			panic(serr)
		}

		fmt.Println(s.CompanyID)

		cmp, cerr := cph.companySrv.Company(s.CompanyID)

		if len(cerr) > 0 {
			fmt.Println("i am the error")
			panic(cerr)
		}

		fmt.Println(cmp.Name)

		post := &entity.Post{}
		post.CompanyID = s.CompanyID
		post.Owner = cmp.Name
		post.Title = r.FormValue("title")
		post.Description = r.FormValue("description")
		post.Category = r.Form.Get("category")

		fmt.Println(post.Category)

		mf, fh, err := r.FormFile("postimg")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		post.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		_, errs := cph.postSrv.StorePost(post)
		// cph.postSrv.StorePost(post)

		if len(errs) > 0 {
			panic(errs)
		}
		// fmt.Println(entity.Company.ID)
		fmt.Println(post)
		fmt.Println("post added to db")

		http.Redirect(w, r, "/admin/cmp-posts", http.StatusSeeOther)

	} else {

		cph.tmpl.ExecuteTemplate(w, "post-job.layout", nil)

	}
}

// CompanyPostsUpdate handle requests on /admin/posts/update
func (cph *CompanyPostHandler) CompanyPostsUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		post, errs := cph.postSrv.Post(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		cph.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", post)

	} else if r.Method == http.MethodPost {

		pst := entity.Post{}
		// pst.ID, _ = strconv.Atoi(r.FormValue("id"))
		pst.Title = r.FormValue("name")
		pst.Description = r.FormValue("description")
		// pst.Image = r.FormValue("image")
		// pst.Category = r.Form.Get("category")

		// mf, _, err := r.FormFile("postimg")

		// if err != nil {
		// 	panic(err)
		// }

		// defer mf.Close()

		// writeFile(&mf, pst.Image)

		// err = cph.postSrv.UpdatePost(pst)

		// if err != nil {
		// 	panic(err)
		// }

		http.Redirect(w, r, "/cmp_home", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/cmp_home", http.StatusSeeOther)
	}

}

// CompanyPostsDelete handle requests on route /admin/posts/delete
func (cph *CompanyPostHandler) CompanyPostsDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		_, errs := cph.postSrv.DeletePost(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

	}

	http.Redirect(w, r, "/cmp_home", http.StatusSeeOther)
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
