package gigcity

import (
	"html/template"
	"net/http"
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
