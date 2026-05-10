package server

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

type Server struct {
	documentHandler documentHandler
	templateHandler templateHandler
}

func New(documentHandler documentHandler, templateHandler templateHandler) *Server {
	return &Server{
		documentHandler: documentHandler,
		templateHandler: templateHandler,
	}
}

func (s *Server) Start(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("POST /api/document", s.documentHandler.Create)
	mux.HandleFunc("POST /api/document/{id}/finish", s.documentHandler.Finish)
	mux.HandleFunc("GET /api/document/{id}/read", s.documentHandler.GetForRead)
	mux.HandleFunc("GET /api/document/{id}/write", s.documentHandler.GetForWrite)
	mux.HandleFunc("POST /api/document/{id}/operations", s.documentHandler.AddOperations)
	mux.HandleFunc("POST /api/document/{id}/override", s.documentHandler.OverrideEditor)

	mux.HandleFunc("GET /document/{id}/edit", s.templateHandler.ServeEditor)
	mux.HandleFunc("GET /document/{id}/view", s.templateHandler.ServeViewer)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /", fs)

	log.Printf("Starting server on port: %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), logger(mux))
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("error dumping http request")
		}
		log.Printf("\n====request====\n%s\n===============", bytes)

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		bytes, err = httputil.DumpResponse(rec.Result(), true)
		if err != nil {
			log.Printf("error dumping http resposne")
		}
		log.Printf("\n====response====\n%s\n================", bytes)

		maps.Copy(w.Header(), rec.Result().Header)
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}
