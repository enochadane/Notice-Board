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

	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type CompanyPostHandler struct {
	tmpl	*template.Template
	postSrv model.PostService
}

func NewCompanyPostHandler(T *template.Template, PS model.PostService) *CompanyPostHandler {
	return &CompanyPostHandler{tmpl: T, postSrv: PS}
}

func (cph *CompanyPostHandler) CompanyPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := cph.postSrv.Posts()
	if err != nil {
		panic(err)
	}
	cph.tmpl.ExecuteTemplate(w, "cmp_post.layout", posts)
}

func (cph *CompanyPostHandler) CompanyPostsNew(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		post := entity.Post{}
		post.Title = r.FormValue("title")
		post.Description = r.FormValue("description")
		post.Type = r.Form.Get("category")

		fmt.Println(post.Type)

		mf, fh, err := r.FormFile("postimg")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		post.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		err = cph.postSrv.StorePost(post)
		fmt.Println(post)

		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "admin/cmp-posts", http.StatusSeeOther)

	} else {

		cph.tmpl.ExecuteTemplate(w, "post-job.layout", nil)

	}
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (cph *CompanyPostHandler) CompanyPostsUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		post, err := cph.postSrv.Post(id)

		if err != nil {
			panic(err)
		}

		cph.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", post)

	} else if r.Method == http.MethodPost {

		pst := entity.Post{}
		pst.Id, _ = strconv.Atoi(r.FormValue("id"))
		pst.Title = r.FormValue("name")
		pst.Description = r.FormValue("description")
		pst.Image = r.FormValue("image")
		pst.Type = r.Form.Get("category")

		mf, _, err := r.FormFile("postimg")

		if err != nil {
			panic(err)
		}

		defer mf.Close()

		writeFile(&mf, pst.Image)

		err = cph.postSrv.UpdatePost(pst)

		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/cmp_home", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/cmp_home", http.StatusSeeOther)
	}

}

// AdminCategoriesDelete handle requests on route /admin/categories/delete
func (cph *CompanyPostHandler) CompanyPostsDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		err = cph.postSrv.DeletePost(id)

		if err != nil {
			panic(err)
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
