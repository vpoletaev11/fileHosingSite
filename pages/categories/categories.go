package categories

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/vpoletaev11/fileHostingSite/dbformat"
	"github.com/vpoletaev11/fileHostingSite/errhand"
	"github.com/vpoletaev11/fileHostingSite/session"
	"github.com/vpoletaev11/fileHostingSite/tmp"
)

const (
	//  path to categories[/categories/] template file
	pathTemplateCategories = "pages/categories/template/categories.html"

	//  path to any category[/categories/*any category*] template file
	pathTemplateAnyCategory = "pages/categories/template/anyCategory.html"
)

const (
	selectFileInfo = "SELECT * FROM files WHERE category = ? ORDER BY uploadDate DESC LIMIT ?, ?;"

	countRows = "SELECT COUNT(*) FROM files WHERE category = ?;"
)

const (
	rowsInPage       = 15 // how many rows of file info will be displayed on page
	maxLinksInNavBar = 25 // how many links will be displayed on navigation bar
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
	UploadedFiles []dbformat.FileInfo
}

// numLink contains relations of page number and page link
type numLink struct {
	NumPage int
	Link    string
}

// anyCategoryPageHandler handling any category[/categories/*any category*] page
func anyCategoryPageHandler(dep session.Dependency, w http.ResponseWriter, r *http.Request) {
	page, err := tmp.CreateTemplate(pathTemplateAnyCategory)
	if err != nil {
		errhand.InternalError(err, w)
		return
	}

	// getting category
	link := r.URL.Path[len("/categories/"):]
	switch link {
	case "other", "games", "documents", "projects", "music":
	default:
		fmt.Fprintln(w, "ERROR: Incorrect category")
		return
	}
	category := link

	// getting count of pages
	pagesCount, err := pagesCount(dep.Db, category)
	if err != nil {
		errhand.InternalError(err, w)
		return
	}

	// getting number of current page
	numPage, err := numPage(r)
	if err != nil {
		fmt.Fprintln(w, "ERROR: Incorrect get request")
		return
	}

	if numPage > pagesCount {
		fmt.Fprintln(w, "ERROR: Incorrect get request")
		return
	}

	// getting files info for current page
	fiCollection, err := dbformat.FormatedFilesInfo(dep.Username, dep.Db, selectFileInfo, category, (numPage-1)*rowsInPage, numPage*rowsInPage)
	if err != nil {
		errhand.InternalError(err, w)
		return
	}

	if pagesCount == 1 {
		err := page.Execute(w, TemplateAnyCategory{Username: dep.Username, UploadedFiles: fiCollection, Title: r.URL.Path[len("/categories/"):]})
		if err != nil {
			errhand.InternalError(err, w)
			return
		}
		return
	}

	// creating navigation bar if count of pages > 1
	numsLinks := navigationBar(pagesCount, numPage, category)
	err = page.Execute(w, TemplateAnyCategory{Username: dep.Username, UploadedFiles: fiCollection, LinkList: numsLinks, Title: r.URL.Path[len("/categories/"):]})
	if err != nil {
		errhand.InternalError(err, w)
		return
	}
	return
}

// Page returns HandleFunc for categories[/categories/] and any category[/categories/*any category*] pages
func Page(dep session.Dependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for categories page
		page, err := tmp.CreateTemplate(pathTemplateCategories)
		if err != nil {
			errhand.InternalError(err, w)
			return
		}
		switch r.Method {
		case "GET":
			if r.URL.Path[len("/categories/"):] == "" {
				err := page.Execute(w, TemplateCategories{Username: dep.Username})
				if err != nil {
					errhand.InternalError(err, w)
				}
				return
			}
			anyCategoryPageHandler(dep, w, r)
			return
		}
	}
}

// pagesCount returns pages count calculated from count MySQL database file info rows
func pagesCount(db *sql.DB, category string) (int, error) {
	rowsCount := 0
	err := db.QueryRow(countRows, category).Scan(&rowsCount)
	if err != nil {
		return 0, err
	}
	pagesCount := (rowsCount-1)/rowsInPage + 1

	return pagesCount, nil
}

// numPage gets number of page from GET request
func numPage(r *http.Request) (int, error) {
	numPageStr := r.URL.Query().Get("p")
	if numPageStr == "" {
		return 1, nil
	}
	numPage := 0
	numPage, err := strconv.Atoi(numPageStr)
	if err != nil {
		return 0, err
	}
	if numPage <= 0 {
		return 0, fmt.Errorf("Incorrect page number")
	}

	return numPage, nil
}

// navigationBar returns array with relations of page number and page link, where page number == page link
func navigationBar(pagesCount, numPage int, category string) []numLink {
	var numsLinks []numLink
	if pagesCount > maxLinksInNavBar {
		// add the first link (literally 1)
		link := "/categories/" + category + "?p=" + "1"
		numsLinks = append(numsLinks, numLink{NumPage: 1, Link: link})

		switch {
		case numPage < 10:
			for i := 2; i <= maxLinksInNavBar; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i, Link: link})
			}
		case numPage >= pagesCount-15:
			for i := numPage - 5; i <= pagesCount-1; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i, Link: link})
			}
		default:
			for i := numPage - 5; i <= numPage+15; i++ {
				pageNum := strconv.Itoa(i)
				link := "/categories/" + category + "?p=" + pageNum
				numsLinks = append(numsLinks, numLink{NumPage: i, Link: link})
			}

		}

		// add the last link == len(pagesCount)
		link = "/categories/" + category + "?p=" + strconv.Itoa(pagesCount)
		numsLinks = append(numsLinks, numLink{NumPage: pagesCount, Link: link})

	} else {
		for i := 1; i <= pagesCount; i++ {
			pageNum := strconv.Itoa(i)
			link := "/categories/" + category + "?p=" + pageNum
			numsLinks = append(numsLinks, numLink{NumPage: i, Link: link})
		}
	}
	return numsLinks
}
