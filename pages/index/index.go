package index

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/index.html"

// TemplateIndex contains fields with warning message and username for login page handler template
type TemplateIndex struct {
	Warning  template.HTML
	Username string
}

// Page returns HandleFunc with access to MySQL database for index page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			page.Execute(w, TemplateIndex{Username: username})
			return
		}
	}
}
