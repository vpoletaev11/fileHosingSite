package tmp

import (
	"html/template"
	"os"
	"path/filepath"
)

// CreateTemplate creates template from inputted template file path
func CreateTemplate(path string) (*template.Template, error) {
	// if working directory != root - directory level will be lowered (case when used tests)
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join("../", path)
		} else {
			break
		}
	}
	// creating template from template file path
	page, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return page, nil
}
