package login

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/vpoletaev11/fileHostingSite/cookie"
	"golang.org/x/crypto/bcrypt"
)

const (
	// query to MySQL database to SELECT password for user
	selectPass = "SELECT password FROM users WHERE username = ?;"

	// query to MySQL database to add cookie
	insertCookie = "INSERT INTO sessions(username, cookie) VALUES(?, ?);"

	// absolute path to template file
	absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/login.html"
)

// TemplateLog contain field with warning message for login page handler template
type TemplateLog struct {
	Warning template.HTML
}

//Page returns HandleFunc with access to MySQL database for login page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for login page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			// handling GET requests and response to them is login page
			templateData := TemplateLog{Warning: ""}
			page.Execute(w, templateData)
			return
		case "POST":
			// getting username and password from POST request
			username := r.FormValue("username")
			password := r.FormValue("password")

			// handle case when len(username) > 20
			if len(username) > 20 {
				templateData := TemplateLog{"<h2 style=\"color:red\">Username cannot be longer than 20 characters</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(password) > 20
			if len(password) > 20 {
				templateData := TemplateLog{"<h2 style=\"color:red\">Password cannot be longer than 20 characters</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handling case when username is non-lowercase
			if username != strings.ToLower(username) {
				templateData := TemplateLog{"<h2 style=\"color:red\">Please use lower case username</h2>"}
				page.Execute(w, templateData)
				return
			}

			// query to MySQL database to SELECT password for user.
			// This query also checks is username exist
			hashPassDB := ""
			db.QueryRow(selectPass, username).Scan(&hashPassDB)

			// handle case when username doesn't exist
			if hashPassDB == "" {
				templateData := TemplateLog{"<h2 style=\"color:red\">Wrong username or password</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when password for username doesn't match with password from MySQL database
			err := comparePasswords(hashPassDB, password)
			if err != nil {
				templateData := TemplateLog{"<h2 style=\"color:red\">Wrong username or password</h2>"}
				page.Execute(w, templateData)
				return
			}

			// creating cookie
			cookie, err := cookie.CreateCookie(username)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			_, err = db.Exec(insertCookie, username, cookie.Value)
			if err != nil {
				fmt.Fprintln(w, err)
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
