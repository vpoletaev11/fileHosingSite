package cookiescleaner

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

const deleteOldSessions = "DELETE FROM sessions WHERE expires <= ?;"

const cookieLifetime = 30 * time.Minute

// Page returns HandleFunc for cookiescleaner[/cookiescleaner] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
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
			fmt.Println(cookiesDeleted)
		}
	}
}
