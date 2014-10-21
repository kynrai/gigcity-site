package gigcity

import (
	"html/template"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Location contains details on locations for GDG Events
type Location struct {
	// ID is the unique ID for each location
	ID string
	// Name is the name of the location
	Name string
	// Address is the street address
	Address string
	// Details is any additonal details for the location, like how to find
	// the group
	Details string
}

// Fetches the next key out of the datastore for the Locations entity
func locationList(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Locations", "default_locationlist", 0, nil)
}

// Handles requests for /admin/location
func locationHandler(w http.ResponseWriter, r *http.Request) {
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

	q := datastore.NewQuery("Locations").Ancestor(locationList(c))
	var locations []Location
	if _, err := q.GetAll(c, &locations); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/admin/overlay.html",
		"static/admin/location.html",
	))

	if err := page.Execute(w, locations); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func addLocationHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == "GET" {
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/admin/overlay.html",
			"static/admin/add-location.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError,
				err.Error())
			return
		}
	} else {
		var loc Location
		loc.Name = r.FormValue("name")
		if loc.Name == "" {
			errorHandler(w, r, http.StatusBadRequest, "location name is required")
			return
		}

		loc.Address = r.FormValue("address")
		if loc.Address == "" {
			errorHandler(w, r, http.StatusBadRequest, "location address is required")
			return
		}

		loc.Details = r.FormValue("details")
		loc.ID = locID(loc.Name)

		key := datastore.NewIncompleteKey(c, "Locations", locationList(c))
		_, err := datastore.Put(c, key, &loc)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		http.Redirect(w, r, "/admin/location", http.StatusFound)
	}
}

func locID(n string) string {
	return strings.ToLower(strings.Replace(n, " ", "-", -1))
}
