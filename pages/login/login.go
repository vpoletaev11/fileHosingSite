package login

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	"github.com/vpoletaev11/fileHostingSite/cookie"
	"github.com/vpoletaev11/fileHostingSite/errhand"
	"golang.org/x/crypto/bcrypt"
)

// absolute path to login[/login] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/login/template/login.html"

const (
	selectPass = "SELECT password FROM users WHERE username = ?;"

	insertCookie = "INSERT INTO sessions (username, cookie, expires) VALUES(?, ?, ?);"
)

// TemplateLog contain data login[/login] page template
type TemplateLog struct {
	Warning template.HTML
}

//Page returns HandleFunc for login[/login] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for login page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("login", "Page", "", err, w)
			return
		}
		switch r.Method {
		case "GET":
			// handling GET requests and response to them is login page
			templateData := TemplateLog{Warning: ""}
			err := page.Execute(w, templateData)
			if err != nil {
				errhand.InternalError("login", "Page", "", err, w)
				return
			}
			return
		case "POST":
			// getting username and password from POST request
			username := r.FormValue("username")
			password := r.FormValue("password")

			//handle case when len(username) == 0
			if username == "" {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Username cannot be empty</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError("login", "Page", "", err, w)
					return
				}
				return
			}

			// handle case when len(password) == 0
			if password == "" {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Password cannot be empty</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError("login", "Page", "", err, w)
					return
				}
				return
			}

			// handle case when len(username) > 20
			if len(username) > 20 {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Username cannot be longer than 20 characters</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError("login", "Page", "", err, w)
					return
				}
				return
			}

			// handle case when len(password) > 20
			if len(password) > 20 {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Password cannot be longer than 20 characters</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError("login", "Page", "", err, w)
					return
				}
				return
			}

			// handling case when username is non-lowercase
			if username != strings.ToLower(username) {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Please use lower case username</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError("login", "Page", "", err, w)
					return
				}
				return
			}

			// query to MySQL database to SELECT password for user.
			// This query also checks is username exist
			hashPassDB := ""
			err := db.QueryRow(selectPass, username).Scan(&hashPassDB)
			if err != nil {
				w.WriteHeader(500)
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when username doesn't exist
			if hashPassDB == "" {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Wrong username or password</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when password for username doesn't match with password from MySQL database
			err = comparePasswords(hashPassDB, password)
			if err != nil {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Wrong username or password</h2>"}
				page.Execute(w, templateData)
				return
			}

			// creating cookie
			cookie := cookie.CreateCookie()

			_, err = db.Exec(insertCookie, username, cookie.Value, cookie.Expires.Format("2006-01-02 15:04:05"))
			if err != nil {
				w.WriteHeader(500)
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				page.Execute(w, templateData)
			}

			// sending cookie
			http.SetCookie(w, &cookie)
			// redirecting to index page
			http.Redirect(w, r, "/", http.StatusFound)
			return
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
