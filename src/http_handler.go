package main

import (
	"net/http"
	"path"
	"strings"
	"github.com/markbates/pkger"
)

type FileSystemHandler func(string) (http.File, error)

func (fs_handler FileSystemHandler) Open(name string) (http.File, error) {
	file, err := fs_handler(name)
	return file, err
}

type PkgerSingleFileHandler string

func (pkger_handler PkgerSingleFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pkger_path := string(pkger_handler)
	file, err := pkger.Open(pkger_path)
	if err != nil {
		http.Error(w, "pkger error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "pkger error", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, req, path.Base(pkger_path), stat.ModTime(), file)
}

func NewWafDevServer() http.Handler {

	mux := http.NewServeMux()

	static_fs := FileSystemHandler(func(name string) (http.File, error) {
		res_path := path.Join("/assets", strings.TrimPrefix(name, "/"))
		file, err := pkger.Open(res_path)
		return file, err
	})

	api_func := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("null"))
	}

	mux.HandleFunc("/api/dev.json", api_func)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(static_fs)))
	mux.Handle("/", PkgerSingleFileHandler("/assets/index.html"))

	return mux
}
