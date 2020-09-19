package registration

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
	"golang.org/x/crypto/bcrypt"
)

// path to registration[/registration] template file
const pathTemplateRegistration = "pages/registration/template/register.html"

const createUser = "INSERT INTO users(username, password, timezone) VALUES(?, ?, ?);"

const (
	maxPasswordLen = 40
	maxUsernameLen = 20
)

// TemplateReg contains data for registration[/registration] page template
type TemplateReg struct {
	Warning template.HTML
}

// Page returns HandleFunc for registration[/registration] page
func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for register page
		page, err := tmp.CreateTemplate(pathTemplateRegistration)
		if err != nil {
			errhand.InternalError(err, w)
			return
		}

		switch r.Method {
		case "GET":
			// handling GET requests and response to them is registration page
			templateData := TemplateReg{""}
			err := page.Execute(w, templateData)
			if err != nil {
				errhand.InternalError(err, w)
				return
			}
			return

		case "POST":
			// getting username and passwords from POST request
			username := r.FormValue("username")
			password1 := r.FormValue("password1")
			password2 := r.FormValue("password2")
			timezone := r.FormValue("timezone")

			err = usernameValidator(username)
			if err != nil {
				templateData := TemplateReg{"<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			err = passwordsValidator(password1, password2)
			if err != nil {
				templateData := TemplateReg{"<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			err = timezoneValidator(timezone)
			if err != nil {
				templateData := TemplateReg{"<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			// creating salted hash from password
			hashedPass, err := hashAndSalt(password1)
			if err != nil {
				templateData := TemplateReg{"<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			// writing username and salted hashed password to MySQL database
			// MySQL database does not allow to enter not unique usernames (username is primary key)
			_, err = db.Exec(createUser, username, hashedPass, timezone)
			if err != nil {
				// handling case when username is not unique
				if strings.Contains(err.Error(), "Error 1062") {
					templateData := TemplateReg{"<h2 style=\"color:red\">Username already used</h2>"}
					err := page.Execute(w, templateData)
					if err != nil {
						errhand.InternalError(err, w)
						return
					}
					return
				}
				// handling internal errors related with query to MySQL database
				templateData := TemplateReg{"<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}
}

func timezoneValidator(timezone string) error {
	switch {
	case timezone == "empty":
		return fmt.Errorf("Please set your timezone")

	case timezone == "":
		return fmt.Errorf("Incorrect timezone")
	}

	_, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("Incorrect timezone")

	}
	return nil
}

func usernameValidator(username string) error {
	switch {
	case len(username) == 0:
		return fmt.Errorf("Username cannot be empty")

	case len(username) > maxUsernameLen:
		return fmt.Errorf("Username cannot be longer than " + strconv.Itoa(maxUsernameLen) + " characters")

	case username != strings.ToLower(username):
		return fmt.Errorf("Please use lower case username")
	}

	return nil
}

func passwordsValidator(password1, password2 string) error {
	switch {
	case len(password1) == 0:
		return fmt.Errorf("Password cannot be empty")

	case len(password2) == 0:
		return fmt.Errorf("Password cannot be empty")

	case len(password1) > maxPasswordLen:
		return fmt.Errorf("Password cannot be longer than " + strconv.Itoa(maxPasswordLen) + " characters")

	case len(password2) > maxPasswordLen:
		return fmt.Errorf("Password cannot be longer than " + strconv.Itoa(maxPasswordLen) + " characters")

	case password1 != password2:
		return fmt.Errorf("Passwords doesn't match")
	}

	return nil
}

// hashAndSalt creates salted hash from password
func hashAndSalt(pwd string) (string, error) {
	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}
