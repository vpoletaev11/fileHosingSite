package upload

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

func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/upload", nil)
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
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="goback">
        <a href="/"><h2>Go back</h2></a>
    </div>

    <div class="uploadForm">
        <form action="" method="post">
            <p>Filename: <input required type="text" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input type="file"></input></p>

            <input type="submit" value="UPLOAD">
        </form>

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
	r, err := http.NewRequest(http.MethodGet, "http://localhost/upload", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(nil)
	sut(w, r)

	assert.Equal(t, 500, w.Code)

	// renaming template file to original filename
	defer func() {
		// renaming template file to original filename
		oldName = newName
		newName = oldName[:lenOrigName]
		err = os.Rename(oldName, newName)
		require.NoError(t, err)
	}()

	// checking error handler works correct
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Internal error. Page not found\n", bodyString)
}
