package upload

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const sendFileInfoToDB = "INSERT INTO files (label, description, owner, category, upload_date) VALUES (?, ?, ?, ?, ?);"

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/upload.html"

// Page returns HandleFunc with access to MySQL database for upload file page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}

		switch r.Method {
		case "GET":
			page.Execute(w, "")
			return
		case "POST":
			filename := r.FormValue("filename")
			description := r.FormValue("description")
			category := r.FormValue("category")
			r.ParseMultipartForm(5 * 1024 * 1025)
			file, handler, err := r.FormFile("uploaded_file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			if len(filename) > 50 {
				fmt.Fprintln(w, "Filename are too long")
				return
			}
			// if filename field in form is empty will be used original filename
			if len(filename) == 0 {
				filename = handler.Filename
			}

			if len(description) > 500 {
				fmt.Fprintln(w, "description are too long")
			}

			_, err = db.Exec(sendFileInfoToDB, filename, description, username, category, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				fmt.Fprintln(w, err)
			}

			id := ""
			db.QueryRow("SELECT LAST_INSERT_ID();").Scan(&id)
			fmt.Println(id)
		}
	}
}
