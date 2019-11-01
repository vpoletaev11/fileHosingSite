package cookie

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	// query to MySQL database to SELECT username from sessions
	//userFromCookie = "SELECT username FROM sessions WHERE cookie=?;"

	// query to MySQL database to select cookie expiring time from session
	getExpiresAndUsername = "SELECT expires, username FROM sessions WHERE cookie=?"

	// query to MySQL database to update cookie expiring time
	updateExpires = "UPDATE sessions SET expires=? WHERE cookie=?;"

	// query to MySQL database to delete session
	deleteSession = "DELETE FROM sessions WHERE cookie=?"
)

// CreateCookie creates cookie for user
func CreateCookie() http.Cookie {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	cookieVal := make([]rune, 60)

	rand.Seed(time.Now().UTC().UnixNano())
	for i := range cookieVal {
		cookieVal[i] = letters[rand.Intn(len(letters))]
	}

	// creating cookie with lifetime
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   string(cookieVal),
		Expires: time.Now().Add(30 * time.Minute),
	}

	return cookie
}

// cookieValidator returns empty username with empty cookie and error if some error with database happends.
// cookieValidator returns username and cookie without error when cookie are:
// 1) came on input,
// 2) doesn't out of date,
// 3) contains in database.
func cookieValidator(db *sql.DB, r *http.Request) (string, http.Cookie, error) {
	// handling case when cookie doesn't came to input
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", http.Cookie{}, nil
	}

	// handling case when cookie aren't contains in database
	var expires time.Time
	username := ""
	err = db.QueryRow(getExpiresAndUsername, cookie.Value).Scan(&expires, &username)
	if err != nil {
		return "", http.Cookie{}, err
	}

	// handling case when cookie is out of date
	if expires.Before(time.Now()) {
		_, err := db.Exec(deleteSession, cookie.Value)
		if err != nil {
			return "", http.Cookie{}, err
		}

		return "", http.Cookie{}, nil
	}

	return username, *cookie, nil
}

// AuthWrapper grants access to pagehandler and extends cookie lifetime if inputted cookie are valid
func AuthWrapper(pageHandler http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking cookie validity
		username, cookie, err := cookieValidator(db, r)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}
		// handling  if cookie invalid
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// extending cookie lifetime
		cookie.Expires = time.Now().Add(30 * time.Minute)
		_, err = db.Exec(updateExpires, time.Now().Add(30*time.Minute).Format("2006-01-02 15:04:05"), cookie.Value)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}
		http.SetCookie(w, &cookie)

		// run page handler
		pageHandler.ServeHTTP(w, r)
	})
}
