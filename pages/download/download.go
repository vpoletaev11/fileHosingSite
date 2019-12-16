package download

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vpoletaev11/fileHostingSite/dbformat"
	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// path to download[/download] template file
const pathTemplateDownload = "pages/download/template/download.html"

const (
	fileInfoDB = "SELECT * FROM files WHERE id = ?;"

	createFileRating = "INSERT INTO filesRating (fileID, voter, rating) VALUES (?, ?, ?);"

	getFileRating = "SELECT rating FROM filesRating WHERE fileID = ?;"

	updateGlobalFileRating = "UPDATE files SET rating= rating - ? + ?  WHERE id=?;"

	updateFileRating = "UPDATE filesRating SET rating=? WHERE fileID=?;"

	updateUserRating = "UPDATE users SET rating= rating -?  + ?  WHERE username= ?;"

	increaseGlobalFileRating = "UPDATE files SET rating=(rating+?) WHERE id=?;"

	selectOwner = "SELECT owner FROM files WHERE id= ?"

	increaseUserRating = "UPDATE users SET rating= rating + ?  WHERE username= ?;"
)

const (
	maxRating = 10  // maximal rating that user can set
	minRating = -10 // minimal rating that user can set
)

// TemplateDownload data for download[/download] page template
type TemplateDownload struct {
	Username string
	FileInfo dbformat.DownloadFileInfo
}

// Page returns HandleFunc for download[/download] page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := tmp.CreateTemplate(pathTemplateDownload)
		if err != nil {
			errhand.InternalError("download", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			fileID := r.URL.Query().Get("id")

			fi, err := dbformat.FormatedDownloadFileInfo(username, db, fileInfoDB, fileID)
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
			if rating > maxRating {
				fmt.Fprintln(w, "INCORRECT POST PARAMETER")
				return
			}
			if rating < minRating {
				fmt.Fprintln(w, "INCORRECT POST PARAMETER")
				return
			}

			id := r.URL.Query().Get("id")

			alreadyRated, err := setRating(db, id, username, rating)
			if err != nil {
				errhand.InternalError("download", "Page", username, err, w)
				return
			}

			if alreadyRated {
				err := changeRating(db, rating, id, username)
				if err != nil {
					errhand.InternalError("download", "Page", username, err, w)
					return
				}
			}

			http.Redirect(w, r, r.RequestURI, 302)
			return

		}
	}
}

// setRating sets rating for file and file owner
func setRating(db *sql.DB, id, username string, rating int) (alreadyRated bool, err error) {
	_, err = db.Exec(createFileRating, id, username, rating)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return true, nil
		}
		return false, err
	}

	_, err = db.Exec(increaseGlobalFileRating, rating, id)
	if err != nil {
		return false, err
	}

	owner := ""
	err = db.QueryRow(selectOwner, id).Scan(&owner)
	if err != nil {
		return false, err
	}
	_, err = db.Exec(increaseUserRating, rating, owner)
	if err != nil {
		return false, err
	}

	return false, nil
}

// changeRating changes rating for file and file owner
func changeRating(db *sql.DB, rating int, id, username string) error {
	var oldRating int
	err := db.QueryRow(getFileRating, id).Scan(&oldRating)
	if err != nil {
		return err
	}

	if oldRating == rating {
		return nil
	}

	_, err = db.Exec(updateGlobalFileRating, oldRating, rating, id)
	if err != nil {
		return err
	}

	_, err = db.Exec(updateFileRating, rating, id)
	if err != nil {
		return err
	}

	owner := ""
	err = db.QueryRow(selectOwner, id).Scan(&owner)
	if err != nil {
		return err
	}

	_, err = db.Exec(updateUserRating, oldRating, rating, owner)
	if err != nil {
		return err
	}
	return nil
}
