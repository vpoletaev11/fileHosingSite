package tmp

import (
	"html/template"
	"path/filepath"
)

// CreateTemplate creates template from inputted template file path
func CreateTemplate(path string) (*template.Template, error) {
	// creating absolute filepath for template file
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// creating template for current page
	page, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return page, nil
}
