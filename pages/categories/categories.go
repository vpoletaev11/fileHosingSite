package categories

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/vpoletaev11/fileHostingSite/database"
	"github.com/vpoletaev11/fileHostingSite/errhand"
)

const (
	// absolute path to categories[/categories/] template file
	absPathTemplateCategories = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/categories/template/categories.html"

	// absolute path to any category[/categories/*any category*] template file
	absPathTemplateAnyCategory = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/categories/template/anyCategory.html"
)

const (
	selectFileInfo = "SELECT * FROM files WHERE category = ? ORDER BY uploadDate DESC LIMIT ?, ?;"

	countRows = "SELECT COUNT(*) FROM files WHERE category = ?;"
)

// TemplateCategories contains data for categories[/categories/] page template
type TemplateCategories struct {
	Warning  template.HTML
	Username string
}

// TemplateAnyCategory contains data for any category[/categories/*any category*] template
type TemplateAnyCategory struct {
	Warning       template.HTML
	Username      string
	Title         string
	LinkList      []numLink
	UploadedFiles []database.FileInfo
}

// numLink contains relations of page number and page link
type numLink struct {
	NumPage int
	Link    string
}

// anyCategoryPageHandler handling any category[/categories/*any category*] page
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
		errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
		return
	}

	rowsCount := 0
	err = db.QueryRow(countRows, category).Scan(&rowsCount)
	if err != nil {
		errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
		return
	}
	if rowsCount == 0 {
		err = page.Execute(w, TemplateAnyCategory{Username: username})
		if err != nil {
			errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
			return
		}
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

	fiCollection, err := database.FormatedFilesInfo(db, selectFileInfo, category, (numPage-1)*15, numPage*15)
	if err != nil {
		errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
		return
	}

	if pagesCount == 1 {
		err := page.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fiCollection})
		if err != nil {
			errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
			return
		}
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
	err = page.Execute(w, TemplateAnyCategory{Username: username, UploadedFiles: fiCollection, LinkList: numsLinks})
	if err != nil {
		errhand.InternalError("categories", "anyCategoryPageHandler", username, err, w)
		return
	}
	return
}

// Page returns HandleFunc for categories[/categories/] and any category[/categories/*any category*] pages
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := template.ParseFiles(absPathTemplateCategories)
		if err != nil {
			errhand.InternalError("categories", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			if r.URL.Path[len("/categories/"):] == "" {
				err := page.Execute(w, TemplateCategories{Username: username})
				if err != nil {
					errhand.InternalError("categories", "Page", username, err, w)
				}
				return
			}
			anyCategoryPageHandler(db, username, w, r)
		}
	}
}
