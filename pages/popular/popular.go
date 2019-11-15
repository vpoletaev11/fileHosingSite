package popular

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/popular/template/popular.html"

// TemplatePopular contains data for popular handler template
type TemplatePopular struct {
	Warning  template.HTML
	Username string
}

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
			page.Execute(w, TemplatePopular{Username: username})
			return

		}
	}
}
