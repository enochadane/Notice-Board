package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"

	"NoticeBoard/model/repository"
	"NoticeBoard/model/service"
	"NoticeBoard/delivery/http/handler"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "godisgood"
	dbname   = "noticeboard"
)

func main()  {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	dbconn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("DB connection established")

	tmpl := template.Must(template.ParseGlob("../../ui/templates/*"))

	userRepo := repository.NewUserRepositoryImpl(dbconn)
	userSrv := service.NewUserServiceImpl(userRepo)

	usrHandler := handler.NewUserHandler(tmpl, userSrv)

	// companyRepo := repository.NewCompanyRepositoryImpl(dbconn)
	// companySrv := service.NewCompanyServiceImpl(companyRepo)

	// cmpHandler := handler.NewCompanyHandler(tmpl, companySrv)

	// mux := http.NewServeMux()
	
	// Server CSS, JS & Images Statically.
	fs := http.FileServer(http.Dir("../../ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	
	http.HandleFunc("/", usrHandler.Signin)
	http.HandleFunc("/signup", usrHandler.Signup)
	http.HandleFunc("/home", usrHandler.Home)

	http.HandleFunc("/signup_account", usrHandler.CreateAccount)
	http.HandleFunc("/login", usrHandler.Login)
	
	// mux.HandleFunc("/signup_account", cmpHandler.CreateAccount)
	// mux.HandleFunc("/login", cmpHandler.Login)

	http.ListenAndServe(":8080", nil)
}
