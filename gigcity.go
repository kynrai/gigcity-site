package gigcity

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Event contains details about GDG events, used when preforming read/write ops
// to the datastore
type Event struct {
	// Title of the event
	Title string
	// Datetime of the event, sent from the browser in YYYY-MM-DDTHH:MM
	// format
	Datetime string
	// The details about the event
	Details string
}

// Under appengine our code runs as a package, not a binary.  Due to this
// define the routes during package initilization.  Normally this wourd happen
// with in main()
func init() {
	// hondle application paths
	http.HandleFunc("/admin/events/add", addEventHandler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/", rootHandler)
}

// Handle errors here, this allows us to control the format of the output rather
// than using http.Error() defaults
func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
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

// Fetches the next index key out of the datastore for the Events entity
func eventList(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Events", "default_eventlist", 0, nil)
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

// Admin page to add new event information to the datastore
func addEventHandler(w http.ResponseWriter, r *http.Request) {
	// use the request information to determine if this is a new session
	c := appengine.NewContext(r)
	// get user information if one is logged in
	u := user.Current(c)
	if u == nil {
		// the person that made the request is anonymous, redirect them to the login
		// page
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			// was unable to get a login URL, so die with a 500 error
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		// set a return URL for when authentication succeeds
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	// check the request method
	if r.Method == "POST" {
		// handle post requests
		g := Event{
			Title:    r.FormValue("title"),
			Datetime: r.FormValue("date"),
			Details:  r.FormValue("details"),
		}

		// get the next available index key
		key := datastore.NewIncompleteKey(c, "Events", eventList(c))
		// write the data to the datastore
		_, err := datastore.Put(c, key, &g)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		// send the user back to the view page once done
		http.Redirect(w, r, "/events", http.StatusFound)
	} else if r.Method == "GET" {
		// handle get requests
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/admin/add-event.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		fmt.Fprint(w, r.Method)
	}
}

// Handles requests to /events
func eventHandler(w http.ResponseWriter, r *http.Request) {
	// use the request information to determine if this is a new session
	c := appengine.NewContext(r)
	// query the Events entity in the datastore
	q := datastore.NewQuery("Events").Ancestor(eventList(c)).Order("-Datetime").Limit(10)
	// create a slice of Event with a capacity of 10 items
	events := make([]Event, 0, 10)
	// store the results into the events slice
	if _, err := q.GetAll(c, &events); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// loop through events, changing events.Datetime from YYYY-MM-DDTHH:MM
	// to YYYY-MM-DD HH:MM
	for key, event := range events {
		t, err := time.Parse("2006-01-02T15:04", event.Datetime)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		events[key].Datetime = t.Format("2006-01-02 03:04 PM")
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/events.html",
	))

	if err := page.Execute(w, events); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
