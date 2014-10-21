package gigcity

import (
	"html/template"
	"net/http"
	"strings"
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

func getID(n string) string {
	return strings.ToLower(strings.Replace(n, " ", "-", -1))
}
