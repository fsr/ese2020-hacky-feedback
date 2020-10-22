package main

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"os"
	"log"
	"strings"
)

type Feedback struct {
      Str []string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("could not parse template: %s", err)
	}

	text, err := ioutil.ReadFile("feedback.txt")
	if err != nil {
		log.Fatalf("could not read feedback file: %s", err)
	}
	f := Feedback {}
	for _, item := range(strings.Split(string(text), "\x00")) {
		if strings.TrimSpace(item) != "" {
			f.Str = append(f.Str, item)
		}
	}

	err = t.Execute(w, f)
	if err != nil {
		log.Printf("template error: %v", err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    f, err := os.OpenFile("feedback.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte(r.FormValue("feedback"))); err != nil {
        log.Fatal(err)
    }

    if _, err := f.Write([]byte("\x00")); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
    http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("static"))))
	http.HandleFunc("/post", handlePost)
	http.ListenAndServe(":7777", nil)
}
