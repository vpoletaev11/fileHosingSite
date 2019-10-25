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
	userFromCookie = "SELECT username FROM sessions WHERE cookie=?;"

	// query to MySQL database to rewrite cookie value
	updateCookie = "UPDATE sessions SET cookie=? WHERE cookie=?;"

	// query to MySQL database to delete session
	deleteSession = "DELETE FROM sessions WHERE cookie=?"
)

// CreateCookie creates cookie for user
func CreateCookie(username string) (http.Cookie, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	cookieVal := make([]rune, 60)

	rand.Seed(time.Now().UTC().UnixNano())
	for i := range cookieVal {
		cookieVal[i] = letters[rand.Intn(len(letters))]
	}

	// creating cookie with lifetime
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   string(cookieVal),
		Expires: expiration,
	}

	return cookie, nil
}

// cookieValidator returns username and cookie when cookie:
// 1) came on input,
// 2) doesn't out of date,
// 3) contains in database.
func cookieValidator(db *sql.DB, r *http.Request) (string, http.Cookie) {
	// handling case when cookie doesn't came to input
	cookie, err := r.Cookie("session_id")
	if err != nil {
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	// handling case when cookie is out of date
	if cookie.Expires.After(time.Now()) {
		_, err = db.Exec(deleteSession, cookie.Value)
		if err != nil {
			fmt.Println(err)
		}
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	// handling case when cookie aren't contains in database
	username := ""
	err = db.QueryRow(userFromCookie, cookie.Value).Scan(&username)
	if err != nil {
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	return username, *cookie
}

// AuthWrapper grants access to pagehandler and extends cookie lifetime if inputed cookie are valid
func AuthWrapper(pageHandler http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking cookie validity
		username, cookie := cookieValidator(db, r)
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// extending cookie lifetime
		cookie.Expires = time.Now().Add(30 * time.Minute)
		http.SetCookie(w, &cookie)

		// run page handler
		pageHandler.ServeHTTP(w, r)
	})
}
