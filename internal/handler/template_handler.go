package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

const templateDirFmt = "./static/%s.html"

type TemplateHandler struct{}

type templateData struct {
	DocumentID string
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) ServeEditor(res http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(fmt.Sprintf(templateDirFmt, "editor"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("error please try again"))
		return
	}

	documentID := req.PathValue("id")

	tmpl.Execute(res, templateData{
		DocumentID: documentID,
	})
}

func (h *TemplateHandler) ServeViewer(res http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(fmt.Sprintf(templateDirFmt, "viewer"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("error please try again"))
		return
	}

	documentID := req.PathValue("id")

	tmpl.Execute(res, templateData{
		DocumentID: documentID,
	})
}
