package admin

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

const deleteOldSessions = "DELETE FROM sessions WHERE expires <= ?;"

const cookieLifetime = 30 * time.Minute

type TemplateAdmin struct {
	Warning template.HTML
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
