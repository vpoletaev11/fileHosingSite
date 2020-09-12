package popular

import (
	"html/template"
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/dbformat"
	"github.com/vpoletaev11/fileHostingSite/session"
	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// path to popular[/popular] template file
const pathTemplatePopular = "pages/popular/template/popular.html"

const selectFileInfo = "SELECT * FROM files WHERE rating >0 ORDER BY rating DESC LIMIT 15;"

// TemplatePopular contains data for popular[/popular] page template
type TemplatePopular struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []dbformat.FileInfo
}

// Page returns HandleFunc for popular[/popular] page
func Page(dep session.Dependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := tmp.CreateTemplate(pathTemplatePopular)
		if err != nil {
			errhand.InternalError("popular", "Page", dep.Username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			fiCollection, err := dbformat.FormatedFilesInfo(dep.Username, dep.Db, selectFileInfo)
			if err != nil {
				errhand.InternalError("popular", "Page", dep.Username, err, w)
				return
			}

			err = page.Execute(w, TemplatePopular{Username: dep.Username, UploadedFiles: fiCollection})
			if err != nil {
				errhand.InternalError("popular", "Page", dep.Username, err, w)
				return
			}
			return

		}
	}
}
