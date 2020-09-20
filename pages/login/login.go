package login

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/vpoletaev11/fileHostingSite/session"
	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
	"golang.org/x/crypto/bcrypt"
)

// path to login[/login] template file
const pathTemplateLogin = "pages/login/template/login.html"

const (
	selectPass = "SELECT password FROM users WHERE username = ?;"
)

const (
	maxPasswordLen = 40
	maxUsernameLen = 20
)

// TemplateLog contain data login[/login] page template
type TemplateLog struct {
	Warning template.HTML
}

//Page returns HandleFunc for login[/login] page
func Page(dep session.Dependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for login page
		page, err := tmp.CreateTemplate(pathTemplateLogin)
		if err != nil {
			errhand.InternalError(err, w)
			return
		}
		switch r.Method {
		case "GET":
			// handling GET requests and response to them is login page
			templateData := TemplateLog{Warning: ""}
			err := page.Execute(w, templateData)
			if err != nil {
				errhand.InternalError(err, w)
				return
			}
			return
		case "POST":
			// getting username and password from POST request
			dep.Username = r.FormValue("username")
			password := r.FormValue("password")

			err := usernameValidator(dep.Username)
			if err != nil {
				templateData := TemplateLog{"<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			err = passwordValidator(password)
			if err != nil {
				templateData := TemplateLog{"<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			// query to MySQL database to SELECT password for user.
			// This query also checks is username exist
			hashPassDB := ""
			err = dep.Db.QueryRow(selectPass, dep.Username).Scan(&hashPassDB)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					templateData := TemplateLog{"<h2 style=\"color:red\">Wrong username or password</h2>"}
					err := page.Execute(w, templateData)
					if err != nil {
						errhand.InternalError(err, w)
						return
					}
					return
				}
				errhand.InternalError(err, w)
				return
			}

			// handle case when username doesn't exist
			if hashPassDB == "" {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Wrong username or password</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			// handle case when password for username doesn't match with password from MySQL database
			err = comparePasswords(hashPassDB, password)
			if err != nil {
				templateData := TemplateLog{Warning: "<h2 style=\"color:red\">Wrong username or password</h2>"}
				err := page.Execute(w, templateData)
				if err != nil {
					errhand.InternalError(err, w)
					return
				}
				return
			}

			// creating cookie
			cookie, err := session.CreateCookie(dep)
			if err != nil {
				errhand.InternalError(err, w)
				return
			}
			// sending cookie
			http.SetCookie(w, &cookie)
			// redirecting to index page
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
}

func passwordValidator(password string) error {
	switch {
	case password == "":
		return fmt.Errorf("Password cannot be empty")

	case len(password) > maxPasswordLen:
		return fmt.Errorf("Password cannot be longer than " + strconv.Itoa(maxPasswordLen) + " characters")
	}

	return nil
}

func usernameValidator(username string) error {
	switch {
	case username == "":
		return fmt.Errorf("Username cannot be empty")

	case len(username) > maxUsernameLen:
		return fmt.Errorf("Username cannot be longer than " + strconv.Itoa(maxUsernameLen) + " characters")

	case username != strings.ToLower(username):
		return fmt.Errorf("Please use lower case username")
	}

	return nil
}

// ComparePasswords compare hashed password with plain.
// In non-matching case CopmarePassword returns error
func comparePasswords(hashedPwd, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	plainPwdByte := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdByte)

	return err
}
