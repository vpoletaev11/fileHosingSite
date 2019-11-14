package download

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/download.html"

// TemplateDownload contains fields with warning message and username for login page handler template
type TemplateDownload struct {
	Warning  template.HTML
	Username string
}

// Page returns HandleFunc with access to MySQL database for download page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			page.Execute(w, TemplateDownload{Username: username})
			return
		}
	}
}
