package logout

import (
	"database/sql"
	"fmt"
	"net/http"
)

// query to MySQL database to DELETE session
const deleteSession = "DELETE FROM sessions WHERE cookie=?"

// Page removes user cookie and redirect to login page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, err = db.Exec(deleteSession, cookie.Value)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newCookie := &http.Cookie{
			Name:   "session_id",
			MaxAge: -1,
		}
		http.SetCookie(w, newCookie)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
