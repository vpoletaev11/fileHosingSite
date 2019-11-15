package registration

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// query to MySQL database to INSERT username and password
const insertLogPass = "INSERT INTO users(username, password) VALUES(?, ?);"

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/registration/template/register.html"

// TemplateReg contain field with warning message for registration page handler template
type TemplateReg struct {
	Warning template.HTML
}

//Page returns HandleFunc with access to MySQL database for registration page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for register page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}

		switch r.Method {
		case "GET":
			// handling GET requests and response to them is registration page
			templateData := TemplateReg{""}
			page.Execute(w, templateData)
			return

		case "POST":
			// getting username and passwords from POST request
			username := r.FormValue("username")
			password1 := r.FormValue("password1")
			password2 := r.FormValue("password2")

			// handle case when len(username) == 0
			if len(username) == 0 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Username cannot be empty</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(password1) == 0
			if len(password1) == 0 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Password cannot be empty</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(password2) == 0
			if len(password2) == 0 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Password cannot be empty</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(username) > 20
			if len(username) > 20 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Username cannot be longer than 20 characters</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(password1) > 20
			if len(password1) > 20 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Password cannot be longer than 20 characters</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handle case when len(password2) > 20
			if len(password2) > 20 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Password cannot be longer than 20 characters</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handling case when username is non-lowercase
			if username != strings.ToLower(username) {
				templateData := TemplateReg{"<h2 style=\"color:red\">Please use lower case username</h2>"}
				page.Execute(w, templateData)
				return
			}

			// handling case when passwords doesn't match
			if password1 != password2 {
				templateData := TemplateReg{"<h2 style=\"color:red\">Passwords doesn't match</h2>"}
				page.Execute(w, templateData)
				return
			}

			// creating salted hash from password
			hashedPass, err := hashAndSalt(password1)
			if err != nil {
				templateData := TemplateReg{"<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				page.Execute(w, templateData)
				return
			}

			// writing username and salted hashed password to MySQL database
			// MySQL database does not allow to enter not unique usernames (username is primary key)
			_, err = db.Exec(insertLogPass, username, hashedPass)
			if err != nil {
				// handling case when username is not unique
				if strings.Contains(err.Error(), "Error 1062") {
					templateData := TemplateReg{"<h2 style=\"color:red\">Username already used</h2>"}
					page.Execute(w, templateData)
					return
				}
				// handling internal errors related with query to MySQL database
				templateData := TemplateReg{"<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				page.Execute(w, templateData)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}
}

// hashAndSalt create salted hash from password
func hashAndSalt(pwd string) (string, error) {
	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}
