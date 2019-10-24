package logout

import (
	"database/sql"
	"fmt"
	"net/http"
)

const deleteCookie = "DELETE FROM fileHostingSite.sessions WHERE cookie=?"

// Page removes user cookie and redirect to login page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, err = db.Exec(deleteCookie, cookie.Value)
		if err != nil {
			fmt.Println(err)
		}

		newCookie := &http.Cookie{
			Name:   "session_id",
			MaxAge: -1,
		}
		http.SetCookie(w, newCookie)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
