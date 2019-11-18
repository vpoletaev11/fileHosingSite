package download

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// absolute path to download[/download] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/download/template/download.html"

const fileInfoDB = "SELECT * FROM files WHERE id = ?;"

// TemplateDownload data for download[/download] page template
type TemplateDownload struct {
	Username string
	FileInfo
}

// FileInfo contains processed file info getted from MySQL database
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

// Page returns HandleFunc for download[/download] page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("download", "Page", username, err, w)
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
				errhand.InternalError("download", "Page", username, err, w)
				return
			}
			fi.DownloadLink = "/files/" + strconv.Itoa(id)
			fi.UploadDate = uploadDateTime.Format("2006-01-02 15:04:05")
			fi.FilesizeMB = fmt.Sprintf("%.6f", float64(filesizeBytes)/1024/1024) + " MB"

			err = page.Execute(w, TemplateDownload{Username: username, FileInfo: fi})
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}
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

					_, err = db.Exec("UPDATE users SET rating= rating -?  + ?  WHERE username= ?;", oldRating, rating, username)
					if err != nil {
						log.Fatal(err)
					}

					http.Redirect(w, r, r.RequestURI, 302)
					return
				}

				errhand.InternalError("download", "Page", username, err, w)
				return
			}
			_, err = db.Exec("UPDATE files SET rating=(rating+?) WHERE id=?;", rating, id)
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}

			username := ""
			err = db.QueryRow("SELECT owner FROM files WHERE id= ?", id).Scan(&username)
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}
			_, err = db.Exec("UPDATE users SET rating= rating + ?  WHERE username= ?;", rating, username)
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}

			http.Redirect(w, r, r.RequestURI, 302)
			return
		}
	}
}
