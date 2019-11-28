package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main()  {
	r := mux.NewRouter()

	// r := mux.NewRouter().StrictSlash(true)

	// Server CSS, JS & Images Statically.
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	r.HandleFunc("/", sign_in)
	r.HandleFunc("/signup", sign_up)
	r.HandleFunc("/home", home)
	r.HandleFunc("/post-job", post_job)
	// r.HandleFunc("/schedule", schedule)
	http.ListenAndServe(":8080", r)
}

func sign_in(w http.ResponseWriter, req *http.Request)  {
	tpl.ExecuteTemplate(w, "signin.layout", nil)
}

func sign_up(w http.ResponseWriter, req *http.Request)  {
	tpl.ExecuteTemplate(w, "signup.layout", nil)
}

func post_job(w http.ResponseWriter, req *http.Request)  {
	tpl.ExecuteTemplate(w, "post-job.layout", nil)
}

func home(w http.ResponseWriter, req *http.Request)  {
	tpl.ExecuteTemplate(w, "home.layout", nil)
}