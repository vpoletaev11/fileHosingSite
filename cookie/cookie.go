package cookie

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	userFromCookie = "SELECT username FROM fileHostingSite.sessions WHERE cookie=?;"

	updateCookie = "UPDATE fileHostingSite.sessions SET cookie=? WHERE cookie=?;"

	deleteCookie = "DELETE FROM fileHostingSite.sessions WHERE cookie=?"
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

// cookieValidator handling cases when cookie:
// 1) don't came on input,
// 2) out of date,
// 3) aren't contains in database.
func cookieValidator(db *sql.DB, r *http.Request) (string, http.Cookie) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	if cookie.Expires.After(time.Now()) {
		_, err = db.Exec(deleteCookie, cookie.Value)
		if err != nil {
			fmt.Println(err)
		}
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	username := ""
	err = db.QueryRow(userFromCookie, cookie.Value).Scan(&username)

	if err != nil {
		// fmt.Println(err)
		emptyCookie := &http.Cookie{}
		return "", *emptyCookie
	}

	return username, *cookie
}

// AuthWrapper afsasdf
func AuthWrapper(pageHandler http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, oldCookie := cookieValidator(db, r)
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		cookie, err := CreateCookie(username)

		db.Exec(updateCookie, cookie.Value, oldCookie.Value)

		if err != nil {
			fmt.Println(w, err)
		}

		http.SetCookie(w, &cookie)
		// run page handler
		pageHandler.ServeHTTP(w, r)
	})
}
