// To send post request on this page use ./deleteCookieBackend.sh. [You may probably need to reconfigure script]
// Script should be executable. To make this use: chmod +x deleteCookieBackend.sh

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

	selectPassandTimezone = "SELECT password, timezone FROM users WHERE username = ?;"
)

const cookieLifetime = 30 * time.Minute

type reqKey struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Page returns HandleFunc for cookiescleaner[/cookiescleaner] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			key := reqKey{}
			err := json.NewDecoder(r.Body).Decode(&key)
			if key.Username != "admin" {
				fmt.Fprintln(w, "Wrong username or password")
				return
			}

			password := ""
			timezone := ""
			err = db.QueryRow(selectPassandTimezone, key.Username).Scan(&password, &timezone)
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}

			err = comparePasswords(password, key.Password)
			if err != nil {
				fmt.Fprintln(w, "Wrong username or password")
				return
			}

			cookiesDeleted, err := deleteCookies(db)
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}
			deletedAt, err := deletedAt(db, key.Username, timezone)
			if err != nil {
				errhand.InternalError("cookiescleaner", "Page", "admin", err, w)
				return
			}

			fmt.Fprintln(w, "deleted", cookiesDeleted, "cookies at:", deletedAt)
		}
	}
}

// deleteCookies deletes expired cookies and returns deleted cookies count
func deleteCookies(db *sql.DB) (cookiesDeleted int64, err error) {
	res, err := db.Exec(deleteOldSessions, time.Now().Add(-cookieLifetime).Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}

	cookiesDeleted, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return cookiesDeleted, nil
}

// deletedAt returns localized time for user when cookies was deleted
func deletedAt(db *sql.DB, username, timezone string) (deletedAt string, err error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	return time.Now().In(location).Format("2006-01-02 15:04:05"), nil
}

// ComparePasswords compare hashed password with plain.
// In non-matching case CopmarePassword returns error
func comparePasswords(hashedPwd, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	plainPwdByte := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdByte)

	return err
}
