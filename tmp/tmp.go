package tmp

import (
	"html/template"
	"os"
	"path/filepath"
)

// CreateTemplate creates template from inputted template file path
func CreateTemplate(path string) (*template.Template, error) {
	//
	for i := 0; i < 10; i++ {
		// getting current working directory
		abspath, err := filepath.Abs("")
		if err != nil {
			return nil, err
		}

		switch {
		// if working directory != root (when runned tests)
		case filepath.Base(abspath) != "fileHostingSite":
			err := os.Chdir("..")
			if err != nil {
				return nil, err
			}

		// if working directory == root (usual running of program)
		case filepath.Base(abspath) == "fileHostingSite":
			abspath = filepath.Join(abspath, path)
			page, err := template.ParseFiles(abspath)
			if err != nil {
				return nil, err
			}
			return page, nil

		default:
			break
		}
	}
	return nil, nil
}
