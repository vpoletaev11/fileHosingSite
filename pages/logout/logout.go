package logout

import (
	"database/sql"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// Page returns HandleFunc that removes user cookie and redirect to login page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		redisConn, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			panic(err)
		}
		defer redisConn.Close()

		_, err = redisConn.Do("DEL", cookie.Value)
		if err != nil {
			errhand.InternalError("logout", "Page", "", err, w)
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
