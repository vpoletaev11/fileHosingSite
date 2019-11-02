package index

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGet(t *testing.T) {
	sut := Page(nil, "")
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>File Hosting</title>
    <link rel="stylesheet" href="assets/css/index.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/upload">Upload file</a></li>
            <li><a href="#">Categories</a></li>
            <li><a href="#">Most popular</a></li>
            <li><a href="#">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>

    </div>

    <div class="label">
        <br><br><br><br><br>
        <p><h1>↓↓↓ NEWLY UPLOADED FILES ↓↓↓</h1></p>
    </div>
</body>`, bodyString)
}

// TestPageMissingTemplate tests case when template file is missing.
// Cannot be runned in parallel.
func TestPageMissingTemplate(t *testing.T) {
	// renaming exists template file
	oldName := absPathTemplate
	newName := absPathTemplate + "edit"
	err := os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/index", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(nil, "")
	sut(w, r)

	// renaming template file to original filename
	defer func() {
		// renaming template file to original filename
		oldName = newName
		newName = oldName[:lenOrigName]
		err = os.Rename(oldName, newName)
		require.NoError(t, err)
	}()

	assert.Equal(t, 500, w.Code)

	// checking error handler works correct
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Internal error. Page not found\n", bodyString)
}
