package popular

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/test"
)

func TestPageSuccessGet(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE rating >0 ORDER BY rating DESC LIMIT 15;").WithArgs().WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)
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
    <title>Popular files</title>
    <link rel="stylesheet" href="assets/css/popular.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/upload">Upload file</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/">Home</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>
    
    <div class="label">
        <br><br><br><br><br>
        <p><h1>↓↓↓ MOST POPULAR FILES ↓↓↓</h1></p>
    </div>

    <div class = "newlyUploadedBox">
        <table border="1" width="100%" cellpadding="5">
                <tr>
                    <th>Filename</th>
                    <th>Filesize</th>
                    <th>Description</th>
                    <th>Owner</th>
                    <th>Category</th>
                    <th>Upload date</th>
                    <th>Rating</th>
                </tr>
                
                <tr>
                    <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                    <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                    <td width="25%" title=description>description</td>
                    <td width="15%">owner</td>
                    <td width="10%"><a href=/categories/other>other</a></td>
                    <td width="15%">2009-11-17 23:34:58</td>
                    <td width="10%">1000</td>
                </tr>
                
        </table>
    </div>
</body>`, bodyString)
}

// TestPageMissingTemplate tests case when template file is missing.
// Cannot be runned in parallel.
func TestPageMissingTemplate(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	// renaming exists template file
	oldName := "../../" + pathTemplatePopular
	newName := "../../" + pathTemplatePopular + "edit"
	err := os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/index", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(dep)
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
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDBError01Get(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE rating >0 ORDER BY rating DESC LIMIT 15;").WithArgs().WillReturnError(fmt.Errorf("testing error"))

	sut := Page(dep)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDBError02Get(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE rating >0 ORDER BY rating DESC LIMIT 15;").WithArgs().WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnError(fmt.Errorf("testing error"))

	sut := Page(dep)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
