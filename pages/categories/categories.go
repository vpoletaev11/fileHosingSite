package categories

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// absolute path to template file
const absPathTemplateCategories = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/categories.html"

const absPathTemplateAnyCategory = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/templates/anyCategory.html"

const selectFileInfo = "SELECT * FROM files where category = ? ORDER BY uploadDate DESC LIMIT 15;"

// TemplateCategories contains fields with warning message and username for login page handler template
type TemplateCategories struct {
	Warning  template.HTML
	Username string
}

// TemplateAnyCategory contains data for any category handler template
type TemplateAnyCategory struct {
	Warning       template.HTML
	Username      string
	Title         string
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

func anyCategoryPage(db *sql.DB, category string) []fileInfoTable {
	rows, err := db.Query(selectFileInfo, category)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var uploadDateTime time.Time

	var fiTableCollection []fileInfoTable
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
	return fiTableCollection
}

// Page returns HandleFunc with access to MySQL database for categories page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplateCategories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal error. Page not found")
			return
		}
		switch r.Method {
		case "GET":
			link := r.URL.Path[len("/categories/"):]
			switch link {
			case "":
				page.Execute(w, TemplateCategories{Username: username})
				return

			case "other":
				pageOther, err := template.ParseFiles(absPathTemplateAnyCategory)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Internal error. Page not found")
					return
				}

				fit := anyCategoryPage(db, "other")
				pageOther.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fit})
				return

			case "games":
				pageOther, err := template.ParseFiles(absPathTemplateAnyCategory)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Internal error. Page not found")
					return
				}

				fit := anyCategoryPage(db, "games")
				pageOther.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fit})
				return

			case "documents":
				pageOther, err := template.ParseFiles(absPathTemplateAnyCategory)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Internal error. Page not found")
					return
				}

				fit := anyCategoryPage(db, "documents")
				pageOther.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fit})
				return

			case "projects":
				pageOther, err := template.ParseFiles(absPathTemplateAnyCategory)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Internal error. Page not found")
					return
				}

				fit := anyCategoryPage(db, "projects")
				pageOther.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fit})
				return

			case "music":
				pageOther, err := template.ParseFiles(absPathTemplateAnyCategory)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Internal error. Page not found")
					return
				}

				fit := anyCategoryPage(db, "music")
				pageOther.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fit})
				return

			default:
				fmt.Fprintln(w, "Page not found")
				return
			}

		}
	}
}
