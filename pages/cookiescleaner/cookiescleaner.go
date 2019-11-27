package cookiescleaner

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vpoletaev11/fileHostingSite/errhand"
	"golang.org/x/crypto/bcrypt"
)

const (
	deleteOldSessions = "DELETE FROM sessions WHERE expires <= ?;"

	selectPass = "SELECT password FROM users WHERE username = ?;"
)

const cookieLifetime = 30 * time.Minute

type reqKey struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Page returns HandleFunc for cookiescleaner[/cookiescleaner] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			key := reqKey{}
			err := json.NewDecoder(r.Body).Decode(&key)
			if key.Login != "admin" {
				fmt.Fprintln(w, "Wrong username or password")
				return
			}

			password := ""
			err = db.QueryRow(selectPass, key.Login).Scan(&password)
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}
			err = comparePasswords(password, key.Password)
			if err != nil {
				fmt.Fprintln(w, "Wrong username or password")
				return
			}

			res, err := db.Exec(deleteOldSessions, time.Now().Add(-cookieLifetime).Format("2006-01-02 15:04:05"))
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}

			cookiesDeleted, err := res.RowsAffected()
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}
			fmt.Fprintln(w, "deleted", cookiesDeleted, "cookies")
		}
	}
}

// ComparePasswords compare hashed password with plain.
// In non-matching case CopmarePassword returns error
func comparePasswords(hashedPwd, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	plainPwdByte := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdByte)

	return err
}
