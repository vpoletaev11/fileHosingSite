package login

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSuccess checks workability Page()
func TestPageSuccess(t *testing.T) {
	// db, sqlmock, err := sqlmock.New()
	// require.NoError(t, err)
	// row := []string{"password"}
	// sqlmock.ExpectQuery("SELECT password FROM users WHERE username = ?;").WithArgs("example").WillReturnRows(sqlmock.NewRows(row).AddRow("$2a$10$ITkHbQjRK6AWs.InpysH5em2Lx4jwzmyYOpvFSturS7hRe6oxzUAu"))
	// sqlmock.ExpectExec("INSERT INTO sessions(username, cookie) VALUES(?, ?);").WithArgs("example", "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu").WillReturnResult()

	sut := Page(nil)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/login", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)

	// testing case when template file is missing:
	// renaming exists template file
	oldName := absPathTemplate
	newName := absPathTemplate + "edit"
	err = os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w = httptest.NewRecorder()

	// running of the page handler with un-exists template file
	sut(w, r)

	// renaming template file to original filename
	oldName = newName
	newName = oldName[:lenOrigName]
	err = os.Rename(oldName, newName)
	require.NoError(t, err)

	// checking error handler works correct
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Internal error. Page not found\n", bodyString)

	// testing GET requests handler:
	// creating expected template
	expectedPage := httptest.NewRecorder()
	expectPageTemplate, err := template.ParseFiles(absPathTemplate)
	require.NoError(t, err)
	expectPageTemplate.Execute(expectedPage, TemplateLog{""})

	// creating actual template
	w = httptest.NewRecorder()
	sut(w, r)

	// comparing expected and actual templates
	assert.Equal(t, expectedPage.Body, w.Body)

}

// tests for comparePasswords():
func TestComparePasswordSuccess(t *testing.T) {
	plainPass := "example"
	hashedPass := "$2a$10$ITkHbQjRK6AWs.InpysH5em2Lx4jwzmyYOpvFSturS7hRe6oxzUAu"

	err := comparePasswords(hashedPass, plainPass)
	if err != nil {
		t.Error(err)
	}
}

func TestComparePasswordError(t *testing.T) {
	plainPass := "example_changed"
	hashedPass := "$2a$10$ITkHbQjRK6AWs.InpysH5em2Lx4jwzmyYOpvFSturS7hRe6oxzUAu"

	err := comparePasswords(hashedPass, plainPass)
	if err != nil {
		return
	}
	t.Error("comparePassword doesn't return error when password and hash doesn't match")
}
