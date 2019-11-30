package admin

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

const deleteOldSessions = "DELETE FROM sessions WHERE expires <= ?;"

const cookieLifetime = 30 * time.Minute

// TemplateAdmin contains data for admin[/admin] page template
type TemplateAdmin struct {
	Warning template.HTML
}

//  path to admin[/admin] template file
const pathTemplateAdmin = "pages/admin/template/admin.html"

// Page returns HandleFunc for admin[/admin] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := tmp.CreateTemplate(pathTemplateAdmin)
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

		case "POST":
			res, err := db.Exec(deleteOldSessions, time.Now().Add(-cookieLifetime).Format("2006-01-02 15:04:05"))
			if err != nil {
				errhand.InternalError("index", "Page", "admin", err, w)
				return
			}

			cookiesDeleted, err := res.RowsAffected()
			if err != nil {
				errhand.InternalError("index", "Page", "admin", err, w)
				return
			}

			err = page.Execute(w, TemplateAdmin{Warning: "Deleted " + template.HTML(strconv.Itoa(int(cookiesDeleted))) + " old cookies"})
			if err != nil {
				errhand.InternalError("index", "Page", "admin", err, w)
				return
			}
			return
		}
	}
}
