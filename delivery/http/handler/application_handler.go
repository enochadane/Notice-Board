package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"io"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/amthesonofGod/Notice-Board/application"

	"github.com/satori/go.uuid"
)

// ApplicationHandler handles user job application requests
type ApplicationHandler struct {
	tmpl	*template.Template
	appSrv	application.ApplicationService
	postSrv post.PostService
}

// NewApplicationHandler initializes and returns new ApplicationHandler
func NewApplicationHandler(T *template.Template, AP application.ApplicationService, PS post.PostService) *ApplicationHandler {
	return &ApplicationHandler{tmpl: T, appSrv: AP, postSrv: PS}
}