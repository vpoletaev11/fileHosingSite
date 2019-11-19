package index

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/database"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

//  absolute path to index[/index] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/index/template/index.html"

const selectFileInfo = "SELECT * FROM files ORDER BY uploadDate DESC LIMIT 15;"

// TemplateIndex contains data for index[/index] page template
type TemplateIndex struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []database.FileInfo
}

// Page returns HandleFunc for index[/index] page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("index", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			fiCollection, err := database.FormatedFilesInfo(db, selectFileInfo)
			if err != nil {
				errhand.InternalError("index", "Page", username, err, w)
				return
			}

			err = page.Execute(w, TemplateIndex{Username: username, UploadedFiles: fiCollection})
			if err != nil {
				errhand.InternalError("index", "Page", username, err, w)
				return
			}
			return
		}
	}
}
