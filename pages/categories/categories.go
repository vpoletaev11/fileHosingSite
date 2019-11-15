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

const selectFileInfo = "SELECT * FROM files WHERE category = ? ORDER BY uploadDate DESC LIMIT ?, ?;"

const countRows = "SELECT COUNT(*) FROM files WHERE category = ?;"

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
	LinkList      []numLink
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
	Label        string
	DownloadLink string
	FilesizeMb   string
	Description  string
	Owner        string
	Category     string
	UploadDate   string
	Rating       int
	//
	LabelComment       string
	FilesizeMbComment  string
	DescriptionComment string
}

type numLink struct {
	NumPage int
	Link    string
}

func anyCategoryPageHandler(db *sql.DB, username string, w http.ResponseWriter, r *http.Request) {
	link := r.URL.Path[len("/categories/"):]
	category := ""
	switch link {

	case "other":
		category = "other"

	case "games":
		category = "games"

	case "documents":
		category = "documents"

	case "projects":
		category = "projects"

	case "music":

		category = "music"

	default:
		http.Redirect(w, r, "/categories/", http.StatusFound)
	}
	//
	//
	//

	page, err := template.ParseFiles(absPathTemplateAnyCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal error. Page not found")
		return
	}

	rowsCount := 0
	db.QueryRow(countRows, category).Scan(&rowsCount)
	if rowsCount == 0 {
		page.Execute(w, TemplateAnyCategory{Username: username})
		return
	}

	numPageStr := r.URL.Query().Get("p")
	numPage := 0
	if numPageStr == "" {
		numPage = 1
	} else {
		numPage, err = strconv.Atoi(numPageStr)
		if err != nil {
			http.Redirect(w, r, "/categories/"+category, http.StatusFound)
			return
		}
	}
	if numPage == 0 {
		numPage++
	}

	pagesCount := rowsCount / 15
	if rowsCount%15 != 0 {
		pagesCount++
	}
	if pagesCount == 0 {
		pagesCount++
	}

	if numPage > pagesCount {
		fmt.Fprintln(w, "ERROR: Incorrect get request")
		return
	}

	rows, err := db.Query(selectFileInfo, category, (numPage-1)*15, numPage*15)
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

		fiTable.DownloadLink = "/download?id=" + strconv.Itoa(fiDB.ID)
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
	if pagesCount == 1 {
		page.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fiTableCollection})
		return
	}

	var numsLinks []numLink
	if pagesCount > 25 {
		link := "/categories/" + category + "?p=" + "1"
		numsLinks = append(numsLinks, numLink{NumPage: 1, Link: link})
		if numPage < 10 {
			for i := 2; i <= 25; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i, Link: link})
			}

		} else if numPage > pagesCount-15 {
			for i := numPage - 5; i <= pagesCount-2; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i + 1, Link: link})
			}

		} else {
			for i := numPage - 5; i < numPage+20; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i + 1, Link: link})
			}
		}
		link = "/categories/" + category + "?p=" + strconv.Itoa(pagesCount)
		numsLinks = append(numsLinks, numLink{NumPage: pagesCount, Link: link})

	} else {
		for i := 0; i != pagesCount; i++ {
			pageNum := strconv.Itoa(i + 1)
			link := "/categories/" + category + "?p=" + pageNum
			numsLinks = append(numsLinks, numLink{NumPage: i + 1, Link: link})
		}
	}
	page.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fiTableCollection, LinkList: numsLinks})
	return
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
			if r.URL.Path[len("/categories/"):] == "" {

				page.Execute(w, TemplateCategories{Username: username})
				return
			}
			anyCategoryPageHandler(db, username, w, r)
		}
	}
}
