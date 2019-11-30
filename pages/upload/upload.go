package upload

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vpoletaev11/fileHostingSite/tmp"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

// path to upload[/upload] template file
const pathTemplateUpload = "pages/upload/template/upload.html"

const (
	sendFileInfoToDB = "INSERT INTO files (label, filesizeBytes, description, owner, category, uploadDate) VALUES (?, ?, ?, ?, ?, ?);"

	deleteFileInfoFromDB = "DELETE FROM files WHERE id = ?"
)

const (
	maxFilenameLen    = 50
	maxDescriptionLen = 500
	maxFilesize       = 1024 * 1024 * 1024
)

// TemplateUpload contains data for login[/login] page template
type TemplateUpload struct {
	Warning  template.HTML
	Username string
}

// PassThru contains reader and total writted on disk bytes
type PassThru struct {
	io.Reader
	total int64 // Total # of bytes transferred
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy().
// This is used while copying to check is the uploaded file size larger than 1GB.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)

	if pt.total == 1024*1024*1024 {
		return 0, fmt.Errorf("File more than 1GB")
	}

	return n, err
}

// Page returns HandleFunc for upload[/upload] file page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := tmp.CreateTemplate(pathTemplateUpload)
		if err != nil {
			errhand.InternalError("upload", "Page", username, err, w)
			return
		}

		switch r.Method {
		case "GET":
			err := page.Execute(w, TemplateUpload{Username: username})
			if err != nil {
				errhand.InternalError("upload", "Page", username, err, w)
				return
			}
			return
		case "POST":
			filename := r.FormValue("filename")
			description := r.FormValue("description")
			category := r.FormValue("category")

			// getting file from upload form
			file, header, err := r.FormFile("uploaded_file")
			if err != nil {
				errhand.InternalError("upload", "Page", username, err, w)
				return
			}
			defer file.Close()

			// if filename field in form is empty will be used original filename
			if len(filename) == 0 {
				filename = header.Filename
			}

			err = fileInfoValidator(header.Size, filename, description, category)
			if err != nil {
				err := page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">" + template.HTML(err.Error()) + "</h2>", Username: username})
				if err != nil {
					errhand.InternalError("upload", "Page", username, err, w)
					return
				}
				return
			}

			// todo: timezone utc
			// sending information about uploaded file to MySQL server
			res, err := db.Exec(sendFileInfoToDB, filename, header.Size, description, username, category, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>", Username: username})
				return
			}

			// getting id of uploaded file from exec
			idInt, err := res.LastInsertId()
			if err != nil {
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">INTERNAL ERROR. Please try later</h2>", Username: username})
				return
			}
			id := strconv.FormatInt(idInt, 10)

			// writting data to file on disk from uploaded file
			f, err := os.Create("files/" + id)
			if err != nil {
				errhand.InternalError("upload", "Page", username, err, w)
				return
			}

			_, err = io.Copy(f, &PassThru{Reader: file})
			if err != nil {
				os.Remove(f.Name())
				page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:red\">Filesize more than 1GB</h2>", Username: username})
				return
			}

			err = page.Execute(w, TemplateUpload{Warning: "<h2 style=\"color:green\">FILE SUCCEEDED UPLOADED</h2>", Username: username})
			if err != nil {
				errhand.InternalError("upload", "Page", username, err, w)
				return
			}
			return
		}
	}
}

func fileInfoValidator(filesize int64, filename, description, category string) error {
	// handling of case when filesize more than 1GB
	switch {
	case filesize > maxFilesize:
		return fmt.Errorf("Filesize cannot be more than 1GB")

	case len(filename) > maxFilenameLen:
		return fmt.Errorf("Filename are too long")

	case len(description) > maxDescriptionLen:
		return fmt.Errorf("Description are too long")
	}

	switch category {
	case "other", "games", "documents", "projects", "music":
	default:
		return fmt.Errorf("Unknown category")
	}

	return nil
}
