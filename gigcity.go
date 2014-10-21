package gigcity

import (
	"html/template"
	"log"
	"net/http"
)

// Under appengine our code runs as a package, not a binary.  Due to this
// define the routes during package initilization.  Normally this wourd happen
// with in main()
func init() {
	// hondle application paths
	http.HandleFunc("/admin/location", locationHandler)
	http.HandleFunc("/admin/location/add", addLocationHandler)
	http.HandleFunc("/admin/events/add", addEventHandler)
	http.HandleFunc("/admin", adminRootHandler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/", rootHandler)
}

// Handle errors here, this allows us to control the format of the output rather
// than using http.Error() defaults
func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	log.Println(err)
	switch status {
	case http.StatusNotFound:
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/404.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.StatusInternalServerError:
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
