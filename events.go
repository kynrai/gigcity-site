package gigcity

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Event contains details about GDG events, used when preforming read/write ops
// to the datastore
type Event struct {
	// ID is the unique ID (URI) for the event
	ID string
	// Title of the event
	Title string
	// Datetime of the event, sent from the browser in YYYY-MM-DDTHH:MM
	// format
	Datetime string
	// LocID is the location the event is being held
	LocID string
	// GooglePlus is the URL to the Google+ event page
	GooglePlus string
	// The details about the event
	Details string
}

// Fetches the next index key out of the datastore for the Events entity
func eventList(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Events", "default_eventlist", 0, nil)
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
		var g Event
		g.Title = r.FormValue("title")
		if g.Title == "" {
			errorHandler(w, r, http.StatusBadRequest, "event title is required")
			return
		}

		g.Datetime = r.FormValue("date")
		if g.Datetime == "" {
			errorHandler(w, r, http.StatusBadRequest, "event date and time is required")
			return
		}

		g.LocID = r.FormValue("location")
		if g.LocID == "" {
			errorHandler(w, r, http.StatusBadRequest, "event location is required")
			return
		}

		g.GooglePlus = r.FormValue("gplus")
		if g.GooglePlus == "" {
			errorHandler(w, r, http.StatusBadRequest, "Google+ event page is required")
			return
		}

		g.Details = r.FormValue("details")
		if g.Details == "" {
			errorHandler(w, r, http.StatusBadRequest, "Event details is required")
			return
		}
		g.ID = getID(g.Title)

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

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	type Content struct {
		EventDetails Event
		LocDetails   Location
	}

	var context Content
	c := appengine.NewContext(r)
	eventID := r.URL.Query().Get(":event")
	if eventID == "" {
		errorHandler(w, r, http.StatusInternalServerError, "no event ID found in URL")
		return
	}

	q := datastore.NewQuery("Events").Filter("ID =", eventID)
	t := q.Run(c)
	for {
		var e Event
		_, err := t.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching event details failed: %v", err)
			break
		}

		context.EventDetails = e
	}

	q = datastore.NewQuery("Locations").Filter("ID =", context.EventDetails.LocID)
	t = q.Run(c)
	for {
		var l Location
		_, err := t.Next(&l)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching location details failed: %v", err)
			break
		}

		context.LocDetails = l
	}

	log.Printf("%#v\n", context)

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/view-event.html",
	))

	if err := page.Execute(w, context); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
