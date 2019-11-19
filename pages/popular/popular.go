package popular

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/database"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// absolute path to popular[/popular] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/popular/template/popular.html"

const selectFileInfo = "SELECT * FROM files ORDER BY rating DESC LIMIT 15;"

// TemplatePopular contains data for popular[/popular] page template
type TemplatePopular struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []database.FileInfo
}

// Page returns HandleFunc for popular[/popular] page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("popular", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			fiCollection, err := database.FormatedFilesInfo(db, selectFileInfo)
			if err != nil {
				errhand.InternalError("popular", "Page", username, err, w)
				return
			}

			err = page.Execute(w, TemplatePopular{Username: username, UploadedFiles: fiCollection})
			if err != nil {
				errhand.InternalError("popular", "Page", username, err, w)
				return
			}
			return

		}
	}
}
