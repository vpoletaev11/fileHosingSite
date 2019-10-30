package upload

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/upload.html"

// Page returns HandleFunc with access to MySQL database for upload file page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}

		switch r.Method {
		case "GET":
			page.Execute(w, "")
			return
		}

	}
}
