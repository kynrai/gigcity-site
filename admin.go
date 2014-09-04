package gigcity

import (
	"fmt"
	"net/http"
	"text/template"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Admin landing page
func adminRootHandler(w http.ResponseWriter, r *http.Request) {
	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/admin/overlay.html",
		"static/admin/index.html",
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
			Location: r.FormValue("location"),
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
			"static/admin/overlay.html",
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
