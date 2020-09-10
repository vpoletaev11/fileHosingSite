package cookie

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	CookieLifetime = 30 * time.Minute
)

type page func(db *sql.DB, username string) http.HandlerFunc

// CreateCookie creates cookie for user
func CreateCookie(username string) (http.Cookie, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return http.Cookie{}, err
	}

	defer conn.Close()

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
		Expires: time.Now().Add(CookieLifetime),
	}

	_, err = conn.Do("SET", cookie.Value, username, "EX", CookieLifetime.Seconds())
	if err != nil {
		return http.Cookie{}, err
	}

	return cookie, nil
}

// cookieValidator returns username and cookie when cookie are:
// 1) came on input,
// 2) doesn't out of date,
// 3) contains in redis.
func cookieValidator(redisConn redis.Conn, r *http.Request) (string, http.Cookie) {
	// handling case when cookie doesn't came to input
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", http.Cookie{}
	}

	username, err := redis.String(redisConn.Do("GET", cookie.Value))
	if err != nil {
		return "", http.Cookie{}
	}

	return username, *cookie
}

// AuthWrapper grants access to pagehandler and extends cookie lifetime if inputted cookie are valid
func AuthWrapper(pageHandler page, db *sql.DB, redisConn redis.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking cookie validity
		username, cookie := cookieValidator(redisConn, r)

		// handling case when cookie invalid
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// extending cookie lifetime
		cookie.Expires = time.Now().Add(CookieLifetime)
		cookie.Path = "/"
		_, err := redisConn.Do("EXPIRE", cookie.Value, CookieLifetime.Seconds())
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
