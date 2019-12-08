package categories

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Categories</title>
    <link rel="stylesheet" href="/assets/css/categories.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/upload">Upload file</a></li>
            <li><a href="/">Home</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <ul class="categoriesList">
        <li><a href="/categories/other" class="categoryLink">Other</a></li>
        <li><a href="/categories/games" class="categoryLink">Games</a></li>
        <li><a href="/categories/documents" class="categoryLink">Documents</a></li>
        <li><a href="/categories/projects" class="categoryLink">Projects</a></li>
        <li><a href="/categories/music" class="categoryLink">Music</a></li>
    </ul>
</body>`, bodyString)
}

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageAnyCategorySuccessGET(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

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

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 0, 15).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
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

	sut := Page(db, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/upload">Upload file</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
    </div>
</body>`, bodyString)
}
