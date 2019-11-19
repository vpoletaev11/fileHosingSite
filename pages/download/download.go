package download

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vpoletaev11/fileHostingSite/database"
	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// absolute path to download[/download] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/download/template/download.html"

const fileInfoDB = "SELECT * FROM files WHERE id = ?;"

// TemplateDownload data for download[/download] page template
type TemplateDownload struct {
	Username string
	FileInfo database.DownloadFileInfo
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

			fi, err := database.FormatedDownloadFileInfo(db, fileInfoDB, fileID)
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}

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
