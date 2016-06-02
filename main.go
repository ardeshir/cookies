package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "cmsdbadmin"
	DBPass  = "cmsdbadmin"
	DBDbase = "CMS"
	port  = ":8090"
)

var database *sql.DB

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}


func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content, page_date, page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		fmt.Fprintln(w, err.Error)
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}

	fmt.Println(pageGUID)

	err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=?",pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		// log.Println("Coudn't get page: " + pageID)
		// log.Println(err.Error)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page " + pageGUID)
		return
	}

	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	tp, _ := template.ParseFiles("templates/start.html")
        tp.Execute(w,nil)
}

func middleHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := r.PostFormValue("message")
	cookie := http.Cookie{Name: "message", Value: "message:" + cookieValue, Expires: time.Now().Add(10 * time.Second), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/finish", 301)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	cookieVal, _ := r.Cookie("message")

	if cookieVal != nil {
		fmt.Fprintln(w, "We got "+string(cookieVal.Value)+", now try to refesh?!")
		cookie := http.Cookie{Name: "message", Value: "", Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Fprintln(w, "Cookies are yummy, bye bye!")
	}
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func main() {

	dbConn := fmt.Sprintf("%s:%s@/%s", DBUser, DBPass, DBDbase)
	// fmt.Println(dbConn)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect to " + DBDbase)
		log.Println(err.Error)
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedirIndex)
	routes.HandleFunc("/home", ServeIndex)
    routes.HandleFunc("/start", startHandler)
	routes.HandleFunc("/view", middleHandler)
	routes.HandleFunc("/finish", finishHandler)
	http.Handle("/", routes)

	log.Fatal(http.ListenAndServe(port, nil))
}
