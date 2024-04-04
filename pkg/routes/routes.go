package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

const (
	endpointHome   = "/"
	endpointAws    = "/aws"
	endpointUpload = "/upload"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func SetRoutes(router *mux.Router) {
	staticDir := "../static/files"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.Handle(endpointHome, &templateHandler{filename: "index.html"}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc(endpointAws, AwsUploadHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
}

func AwsUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}

	resultadoCh := make(chan map[string]interface{})
	errCh := make(chan error)
	var wg sync.WaitGroup

	if r.Method == http.MethodPost {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.ParseMultipartForm(10 << 20)
			file, handler, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()
			fmt.Println(handler.Filename)
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"message": "File uploaded successfully"}
			resultadoCh <- response
		}()

		go func() {
			wg.Wait()
			close(resultadoCh)
			close(errCh)
		}()

		select {
		case res := <-resultadoCh:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	}
}
