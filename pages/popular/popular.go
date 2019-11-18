package popular

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// absolute path to popular[/popular] template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/popular/template/popular.html"

const selectFileInfo = "SELECT * FROM files ORDER BY rating DESC LIMIT 15;"

// TemplatePopular contains data for popular[/popular] page template
type TemplatePopular struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []fileInfoTable
}

// fileInfoDB contains file info getted from MySQL database
type fileInfoDB struct {
	ID            int
	Label         string
	FilesizeBytes int
	Description   string
	Owner         string
	Category      string
	UploadDate    string
	Rating        int
}

// fileInfoTable contains processed file info from fileInfoDB{}
type fileInfoTable struct {
	Label       string
	Link        string
	FilesizeMb  string
	Description string
	Owner       string
	Category    string
	UploadDate  string
	Rating      int
	//
	LabelComment       string
	FilesizeMbComment  string
	DescriptionComment string
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
			rows, err := db.Query(selectFileInfo)
			if err != nil {
				errhand.InternalError("popular", "Page", username, err, w)
				return
			}
			defer rows.Close()

			var fiTableCollection []fileInfoTable
			var uploadDateTime time.Time
			for rows.Next() {
				fiDB := new(fileInfoDB)
				fiTable := new(fileInfoTable)

				err := rows.Scan(
					&fiDB.ID,
					&fiDB.Label,
					&fiDB.FilesizeBytes,
					&fiDB.Description,
					&fiDB.Owner,
					&fiDB.Category,
					&uploadDateTime,
					&fiDB.Rating,
				)
				if err != nil {
					errhand.InternalError("popular", "Page", username, err, w)
					return
				}
				fiDB.UploadDate = uploadDateTime.Format("2006-01-02 15:04:05")

				if len(fiDB.Label) > 20 {
					fiTable.Label = fiDB.Label[:20] + "..."
				} else {
					fiTable.Label = fiDB.Label
				}

				if len(fiDB.Description) > 30 {
					fiTable.Description = fiDB.Description[:25] + "..."
				} else {
					fiTable.Description = fiDB.Description
				}

				fiTable.Link = "/download?id=" + strconv.Itoa(fiDB.ID)
				fiTable.FilesizeMb = fmt.Sprintf("%.4f", float64(fiDB.FilesizeBytes)/1024/1024) + " MB"

				fiTable.Owner = fiDB.Owner
				fiTable.Category = fiDB.Category
				fiTable.UploadDate = fiDB.UploadDate
				fiTable.Rating = fiDB.Rating

				fiTable.LabelComment = fiDB.Label
				fiTable.DescriptionComment = fiDB.Description
				fiTable.FilesizeMbComment = strconv.Itoa(fiDB.FilesizeBytes) + " Bytes"

				fiTableCollection = append(fiTableCollection, *fiTable)
			}
			err = page.Execute(w, TemplatePopular{Username: username, UploadedFiles: fiTableCollection})
			if err != nil {
				errhand.InternalError("popular", "Page", username, err, w)
				return
			}
			return

		}
	}
}
