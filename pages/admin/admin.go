package admin

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

type TemplateAdmin struct {
}

//  absolute path to admin[/admin] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/admin/template/admin.html"

// Page returns HandleFunc for admin[/admin] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("index", "Page", "admin", err, w)
			return
		}

		switch r.Method {
		case "GET":
			err = page.Execute(w, TemplateAdmin{})
			if err != nil {
				errhand.InternalError("index", "Page", "admin", err, w)
				return
			}
			return
		}
	}
}
