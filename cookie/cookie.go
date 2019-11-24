package cookie

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	getExpiresAndUsername = "SELECT expires, username FROM sessions WHERE cookie=?;"

	updateExpires = "UPDATE sessions SET expires=? WHERE cookie=?;"

	deleteSession = "DELETE FROM sessions WHERE cookie=?"
)

const (
	cookieLifetime = 30 * time.Minute
)

type page func(db *sql.DB, username string) http.HandlerFunc

type adminPage func(db *sql.DB) http.HandlerFunc

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
		Path:    "/",
		Value:   string(cookieVal),
		Expires: time.Now().Add(cookieLifetime),
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
func AuthWrapper(pageHandler page, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking cookie validity
		username, cookie, err := cookieValidator(db, r)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}

		// handling case when cookie invalid
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// extending cookie lifetime
		cookie.Expires = time.Now().Add(cookieLifetime)
		cookie.Path = "/"
		_, err = db.Exec(updateExpires, cookie.Expires.Format("2006-01-02 15:04:05"), cookie.Value)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}
		http.SetCookie(w, &cookie)

		pageHandler := pageHandler(db, username)
		// run page handler
		pageHandler.ServeHTTP(w, r)
	})
}

// AdminAuthWrapper grants access to admin page and extends cookie lifetime if inputted cookie are valid
func AdminAuthWrapper(pageHandler adminPage, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking cookie validity
		username, cookie, err := cookieValidator(db, r)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}

		// handling case when cookie invalid
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if username != "admin" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// extending cookie lifetime
		cookie.Expires = time.Now().Add(30 * time.Minute)
		cookie.Path = "/"
		_, err = db.Exec(updateExpires, cookie.Expires.Format("2006-01-02 15:04:05"), cookie.Value)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "INTERNAL ERROR. Please try later.")
			return
		}
		http.SetCookie(w, &cookie)

		pageHandler := pageHandler(db)
		// run page handler
		pageHandler.ServeHTTP(w, r)
	})
}
