package gigcity

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// LearnEvent contains details about GDG study groups, used when preforming read/write ops to the datastore
type LearnEvent struct {
	// ID is the unique ID (URI) for the study group event
	ID string
	// Title of the Study Group
	Title string
	// Datetime of the study group event
	Datetime string
	// LocID is the location of that the study group meets at
	LocID string
	// Details holds information regarding the event
	Details string
}

func learnList(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "LearnEvent", "default_learneventlist", 0, nil)
}

func learningHandler(w http.ResponseWriter, r *http.Request) {
	// use the request information to determine if this is a new session
	c := appengine.NewContext(r)
	// query the Learn entity in the database
	q := datastore.NewQuery("LearnEvent").Ancestor(learnList(c)).Limit(10)
	// create a slice of LearnEvent with a capacity of 10 items
	learn := make([]LearnEvent, 0, 10)
	// store the results into the learn slice
	if _, err := q.GetAll(c, &learn); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/learn.html",
	))

	if err := page.Execute(w, learn); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func addLearningHandler(w http.ResponseWriter, r *http.Request) {
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
		// set a retrun URL for when authentication succeeds
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	// check the request method
	if r.Method == "POST" {
		var l LearnEvent
		l.Title = r.FormValue("title")
		if l.Title == "" {
			errorHandler(w, r, http.StatusBadRequest, "study group name is required")
			return
		}

		l.Datetime = r.FormValue("date")
		if l.Datetime == "" {
			errorHandler(w, r, http.StatusBadRequest, "study group date and time is requred")
			return
		}

		l.LocID = r.FormValue("location")
		if l.LocID == "" {
			errorHandler(w, r, http.StatusBadRequest, "study group location is required")
			return
		}

		l.Details = r.FormValue("details")
		if l.Details == "" {
			errorHandler(w, r, http.StatusBadRequest, "study group details is required")
			return
		}

		l.ID = getID(l.Title)

		// get the next available index key
		key := datastore.NewIncompleteKey(c, "LearnEvent", learnList(c))
		_, err := datastore.Put(c, key, &l)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		// send the user back to the view page once done
		http.Redirect(w, r, "/learning", http.StatusFound)
	} else {
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/admin/overlay.html",
			"static/admin/add-learn.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
