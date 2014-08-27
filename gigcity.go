package gigcity

import (
	"fmt"
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Event struct {
	Title    string
	Datetime string
	Details  string
}

func init() {
	// handle static paths
	http.Handle("/static/img/", http.StripPrefix("/static/img/", http.FileServer(http.Dir("static/img"))))

	// hondle application paths
	http.HandleFunc("/admin/events/add", addEventHandler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/", rootHandler)
}

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func eventList(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Events", "default_eventlist", 0, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
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

func addEventHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
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

		key := datastore.NewIncompleteKey(c, "Events", eventList(c))
		_, err := datastore.Put(c, key, &g)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
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

func eventHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Events").Ancestor(eventList(c)).Order("-Datetime").Limit(10)
	events := make([]Event, 0, 10)
	if _, err := q.GetAll(c, &events); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
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
