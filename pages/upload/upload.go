package upload

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// send fileinfo of uploaded file into MySQL database
	sendFileInfoToDB = "INSERT INTO files (label, filesizeBytes, description, owner, category, upload_date) VALUES (?, ?, ?, ?, ?, ?);"

	// delete fileinfo of uploaded file into MySQL database
	deleteFileInfoFromDB = "DELETE FROM files WHERE id = ?"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/upload.html"

// TemplateUpload contains fields with warning message and username for login page handler template
type TemplateUpload struct {
	Warning  template.HTML
	Username string
}

// Page returns HandleFunc with access to MySQL database for upload file page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "INTERNAL ERROR. Page not found")
			return
		}

		switch r.Method {
		case "GET":
			page.Execute(w, TemplateUpload{Username: username})
			return
		case "POST":
			filename := r.FormValue("filename")
			description := r.FormValue("description")
			category := r.FormValue("category")

			// getting file from upload form
			r.ParseMultipartForm(5 * 1024 * 1025)
			file, handler, err := r.FormFile("uploaded_file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			// handling of case when filesize more than 1GB
			if handler.Size > 1000000000 {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">Filesize cannot bo more than 1GB</h2>"})
			}

			// handling of case when in form field filename len(filename) > 50
			if len(filename) > 50 {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">Filename are too long</h2>"})
				return
			}

			// if filename field in form is empty will be used original filename
			if len(filename) == 0 {
				filename = handler.Filename
			}

			// handling of case when in form field description len(description) > 500
			if len(description) > 500 {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">Description are too long</h2>"})
				return
			}

			// sending information about uploaded file to MySQL server
			_, err = db.Exec(sendFileInfoToDB, filename, handler.Size, description, username, category, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"})
				return
			}

			// getting id of uploaded file from MySQL server
			id := ""
			err = db.QueryRow("SELECT LAST_INSERT_ID();").Scan(&id)
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"})

				_, err := db.Exec(deleteFileInfoFromDB, id)
				if err != nil {
					log.Println(err)
				}
				return
			}

			// creating file on disk with name == id
			f, err := os.Create("files/" + id)
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"})

				_, err := db.Exec(deleteFileInfoFromDB, id)
				if err != nil {
					log.Println(err)
				}
				return
			}

			// writting data to file on disk from uploaded file
			_, err = io.Copy(f, file)
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>"})

				_, err := db.Exec(deleteFileInfoFromDB, id)
				if err != nil {
					log.Println(err)
				}
				return
			}

			page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:green\">FILE SUCCEEDED UPLOADED</h2>"})
			return
		}
	}
}
