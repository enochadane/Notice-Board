package handler

import (
	"html/template"

	"github.com/amthesonofGod/Notice-Board/application"
	"github.com/amthesonofGod/Notice-Board/post"
)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl    *template.Template
	appSrv  application.ApplicationService
	postSrv post.PostService
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP application.ApplicationService, PS post.PostService) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, postSrv: PS}
}
