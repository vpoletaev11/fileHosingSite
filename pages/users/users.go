package users

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const selectUsers = "SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15;"

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/users/template/users.html"

// TemplateUsers contains fields with warning message and username for users page handler template
type TemplateUsers struct {
	Warning  template.HTML
	Username string
	UserList []UserInfo
}

type UserInfo struct {
	Username string
	Rating   int
}

// Page returns HandleFunc with access to MySQL database for index page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			rows, err := db.Query(selectUsers)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			usersInfo := []UserInfo{}
			for rows.Next() {
				ui := UserInfo{}

				err := rows.Scan(
					&ui.Username,
					&ui.Rating,
				)
				if err != nil {
					log.Fatal(err)
				}
				usersInfo = append(usersInfo, ui)
			}

			page.Execute(w, TemplateUsers{Username: username, UserList: usersInfo})
			return
		}
	}
}
