package main

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"os"
	"log"
	"strings"
)


func handleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("could not parse template: %s", err)
	}

	text, err := ioutil.ReadFile("feedback.txt")
	if err != nil {
		return
	}	
	t.Execute(w, strings.Split(string(text), "\n"))
 
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    f, err := os.OpenFile("feedback.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte(r.FormValue("feedbacktext"))); err != nil {
        log.Fatal(err)
    }

    if _, err := f.Write([]byte("\n")); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/post", handlePost)
	http.ListenAndServe(":7777", nil)
}
