package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

var index_template = template.Must(template.ParseFiles("index.html"))

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/", previewHandler).Methods("POST")
	router.HandleFunc("/style.css", styleHandler)
	http.Handle("/", router)

	fmt.Println("Starting server on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homeHandler(w http.ResponseWriter, router *http.Request) {
	err := index_template.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func previewHandler(w http.ResponseWriter, router *http.Request) {
	router.ParseForm()
	markdwn := []byte(router.FormValue("markdown"))
	html := blackfriday.MarkdownBasic(markdwn)
	w.Write(html)
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("style.css")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	contentType := "text/css"
	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(content))
}
