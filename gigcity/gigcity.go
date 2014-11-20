package gigcity

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
)

// Under appengine our code runs as a package, not a binary.  Due to this
// define the routes during package initilization.  Normally this wourd happen
// with in main()
func init() {
	// hondle application paths
	m := pat.New()
	m.Get("/admin/location/add", http.HandlerFunc(addLocationHandler))
	m.Post("/admin/location/add", http.HandlerFunc(addLocationHandler))
	m.Get("/admin/location", http.HandlerFunc(locationHandler))
	m.Get("/admin/events/add", http.HandlerFunc(addEventHandler))
	m.Post("/admin/events/add", http.HandlerFunc(addEventHandler))
	m.Get("/admin", http.HandlerFunc(adminRootHandler))
	m.Get("/coc", http.HandlerFunc(cocHandler))
	m.Get("/events/:event", http.HandlerFunc(getEventHandler))
	m.Get("/events", http.HandlerFunc(eventHandler))
	m.Get("/about", http.HandlerFunc(aboutHandler))
	m.Get("/", http.HandlerFunc(rootHandler))
	http.Handle("/", m)
}

// Handle messages that should be written out to the log.  lvl is the level of the message
// and msg contains the message body
func logHandler(lvl, msg string) {
	// switch on message log level
	// TODO work out a good way to see if this is running in a dev
	// environment
	switch lvl {
	case "INFO":
		log.Print("[INFO]: " + logTime() + " " + msg)
	case "WARN":
		log.Print("[WARNING]: " + logTime() + " " + msg)
	case "ERROR":
		log.Print("[ERROR]: " + logTime() + " " + msg)
	case "FATAL":
		log.Fatal("[FATEL]: " + logTime() + " " + msg)
	}
}

// returns the current time in RFC3339 format
func logTime() string {
	return time.Now().Format(time.RFC3339)
}

// Handle errors here, this allows us to control the format of the output rather
// than using http.Error() defaults
func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		logHandler("ERROR", fmt.Sprintf("client %s tried to request %v", r.RemoteAddr, r.URL.Path))
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/404.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.StatusInternalServerError:
		logHandler("ERROR", fmt.Sprintf("an internal server error occured when %s requested %s with error:\n%s", r.RemoteAddr, r.URL.Path, err))
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/500.html",
		))

		if err := page.Execute(w, nil); err != nil {
			// IF for some reason the tempalets for 500 errors fails, fallback
			// on http.Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Handles requests to '/' as well as any unmatched routes to the server
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// If the request is not for the root of the app, then it is a 404
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/index.html",
	))

	if err := page.Execute(w, nil); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// Handles requests to /about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/about.html",
	))

	if err := page.Execute(w, nil); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// Handles requests to /coc
func cocHandler(w http.ResponseWriter, r *http.Request) {
	type Organizer struct {
		Name, Role, Email, GooglePlus, IRC string
	}

	var organizers []Organizer
	organizers = append(organizers, Organizer{Name: "Adam Jimerson", Role: "Lead Organizer", Email: "vendion@gmail.com", GooglePlus: "https://google.com/+AdamJimerson", IRC: "vendion"})

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/coc.html",
	))

	if err := page.Execute(w, organizers); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
