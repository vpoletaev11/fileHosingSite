package upload

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// Page returns HandleFunc with access to MySQL database for upload file page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles("templates/upload.html")
		if err != nil {
			fmt.Fprintln(w, err)
		}

		switch r.Method {
		case "GET":
			page.Execute(w, "")
		}

	}
}
