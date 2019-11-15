package download

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/download/template/download.html"

const fileInfoDB = "SELECT * FROM files WHERE id = ?;"

// TemplateDownload contains fields with warning message and username for login page handler template
type TemplateDownload struct {
	Username string
	FileInfo
}

type FileInfo struct {
	DownloadLink string
	Label        string
	FilesizeMB   string
	Description  string
	Owner        string
	Category     string
	UploadDate   string
	Rating       int
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
			id := 0
			filesizeBytes := 0
			err := db.QueryRow(fileInfoDB, fileID).Scan(
				&id,
				&fi.Label,
				&filesizeBytes,
				&fi.Description,
				&fi.Owner,
				&fi.Category,
				&uploadDateTime,
				&fi.Rating,
			)
			if err != nil {
				fmt.Fprintln(w, "Internal error. Page not found")
				return
			}
			fi.DownloadLink = "/files/" + strconv.Itoa(id)
			fi.UploadDate = uploadDateTime.Format("2006-01-02 15:04:05")
			fi.FilesizeMB = fmt.Sprintf("%.6f", float64(filesizeBytes)/1024/1024) + " MB"

			page.Execute(w, TemplateDownload{Username: username, FileInfo: fi})
			return

		case "POST":
			ratingStr := r.FormValue("rating")
			rating, err := strconv.Atoi(ratingStr)
			if err != nil {
				fmt.Fprintln(w, "INCORRECT POST PARAMETER")
				return
			}
			if rating > 10 {
				fmt.Fprintln(w, "INCORRECT POST PARAMETER")
				return
			}
			if rating < -10 {
				fmt.Fprintln(w, "INCORRECT POST PARAMETER")
				return
			}

			id := r.RequestURI[len("/download?id="):]
			_, err = db.Exec("INSERT INTO filesRating (fileID, voter, rating) VALUES (?, ?, ?);", id, username, rating)
			if err != nil {
				if strings.Contains(err.Error(), "Error 1062") {
					var oldRating int
					err := db.QueryRow("SELECT rating FROM filesRating WHERE fileID = ?;", id).Scan(&oldRating)
					fmt.Println(oldRating, rating)
					if err != nil {
						fmt.Println(err)
						return
					}
					if oldRating == rating {
						http.Redirect(w, r, r.RequestURI, 302)
						return
					}

					_, err = db.Exec("UPDATE files SET rating= rating - ? + ?  WHERE id=?;", oldRating, rating, id)
					if err != nil {
						fmt.Println(err)
						return
					}

					_, err = db.Exec("UPDATE filesRating SET rating=? WHERE fileID=?;", rating, id)
					if err != nil {
						fmt.Println(err)
						return
					}
					http.Redirect(w, r, r.RequestURI, 302)
					return
				}
				fmt.Println(err)
				return
			}
			_, err = db.Exec("UPDATE files SET rating=(rating+?) WHERE id=?;", rating, id)
			if err != nil {
				fmt.Println(err)
				return
			}

			http.Redirect(w, r, r.RequestURI, 302)
			return
		}
	}
}
