package download

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/download.html"

const fileInfoDB = "SELECT * FROM files WHERE id = ?;"

// TemplateDownload contains fields with warning message and username for login page handler template
type TemplateDownload struct {
	Warning  template.HTML
	Username string
	FileInfo
}

type FileInfo struct {
	ID          int
	Label       string
	FilesizeMB  string
	Description string
	Owner       string
	Category    string
	UploadDate  string
	Rating      int
}

// Page returns HandleFunc with access to MySQL database for download page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			fileID := r.URL.Query().Get("id")

			fi := FileInfo{}

			var uploadDateTime time.Time
			filesizeBytes := 0
			err := db.QueryRow(fileInfoDB, fileID).Scan(
				&fi.ID,
				&fi.Label,
				&filesizeBytes,
				&fi.Description,
				&fi.Owner,
				&fi.Category,
				&uploadDateTime,
				&fi.Rating,
			)
			if err != nil {
				fmt.Println(err)
			}
			fi.UploadDate = uploadDateTime.Format("2006-01-02 15:04:05")
			//fi.FilesizeMB = filesizeBytes / 1024 / 1024
			fi.FilesizeMB = fmt.Sprintf("%.4f", float64(filesizeBytes)/1024/1024) + " MB"

			page.Execute(w, TemplateDownload{Username: username, FileInfo: fi})
			return
		}
	}
}
