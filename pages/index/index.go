package index

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

// Page returns HandleFunc with access to MySQL database for index page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles("templates/index.html")
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Internal error. Page not found")
		}
		switch r.Method {
		case "GET":
			page.Execute(w, "")
			return
		}
	}
}
