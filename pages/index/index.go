package index

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/dbformat"
	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// path to index[/index] template file
const pathTemplateIndex = "pages/index/template/index.html"

const selectFileInfo = "SELECT * FROM files ORDER BY uploadDate DESC LIMIT 15;"

// TemplateIndex contains data for index[/index] page template
type TemplateIndex struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []dbformat.FileInfo
}

// Page returns HandleFunc for index[/index] page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := tmp.CreateTemplate(pathTemplateIndex)
		if err != nil {
			errhand.InternalError("index", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			fiCollection, err := dbformat.FormatedFilesInfo(username, db, selectFileInfo)
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
