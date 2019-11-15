package index

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const selectFileInfo = "SELECT * FROM files ORDER BY uploadDate DESC LIMIT 15;"

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/index.html"

// TemplateIndex contains fields with warning message and username for index page handler template
type TemplateIndex struct {
	Warning       template.HTML
	Username      string
	UploadedFiles []fileInfoTable
}

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

type fileInfoTable struct {
	Label       string
	Link        template.HTML
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

// Page returns HandleFunc with access to MySQL database for index page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			rows, err := db.Query(selectFileInfo)
			if err != nil {
				log.Fatal(err)
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
					log.Fatal(err)
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

				fiTable.Link = template.HTML("/download?id=" + strconv.Itoa(fiDB.ID))
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

			page.Execute(w, TemplateIndex{Username: username, UploadedFiles: fiTableCollection})
			return
		}
	}
}
