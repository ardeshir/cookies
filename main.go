package main

import (
	"fmt"
	"html/templates"
	"log"
	"net/http"
	"time"
)

var (
	tmpl = template.Must(template.ParseGlob("templates/*"))
	port = ":8090"
)

func startHandler(w http.ResposeWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "start.html", nil)
	if err != nil {
		log.Fatal("Templates start.html missing")
	}
}

func middleHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := r.PostFormValue("message")
	cookie := http.Cookie{Name: "message", Value: "message:" + cookieValue, Expires: time.Now().Add(60 * time.Second), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/finish", 301)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue, _ := r.Cookie("message")

	if cookieValue != nil {
		fmt.Fprintln(w, "We got "+string(cookeiValue.Value)+", now try to refesh?!")
		cookie := http.Cookie{Name: "message", Value: "", Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Fprintln(w, "Cookies are yummy, bye bye!")
	}
}

func main() {

	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/view", middleHandler)
	http.HandleFunc("/finish", finishHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
