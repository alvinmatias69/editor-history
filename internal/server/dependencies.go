package server

import (
	"net/http"
)

type documentHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Finish(http.ResponseWriter, *http.Request)
	GetForWrite(http.ResponseWriter, *http.Request)
	GetForRead(http.ResponseWriter, *http.Request)
	AddOperations(http.ResponseWriter, *http.Request)
	OverrideEditor(http.ResponseWriter, *http.Request)
}

type templateHandler interface {
	ServeEditor(http.ResponseWriter, *http.Request)
	ServeViewer(http.ResponseWriter, *http.Request)
}
