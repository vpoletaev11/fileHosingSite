package index_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/index"
	"github.com/vpoletaev11/fileHostingSite/test"
)

// TestPageSuccessGET checks workability of GET requests handler in Page()
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

	sqlMock.ExpectQuery("SELECT \\* FROM files ORDER BY uploadDate DESC LIMIT 15;").WithArgs().WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
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

	sut := index.Page(dep)
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
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="label">
        <br><br><br><br><br>
        <p><h1>↓↓↓ NEWLY UPLOADED FILES ↓↓↓</h1></p>
    </div>

    <div class = "newlyUploadedBox">
        <div class = "newlyUploadedContent">
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
            </ul>
        </div>
    </div>
</body>`, bodyString)
}

func TestPageDBError01Get(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT \\* FROM files ORDER BY uploadDate DESC LIMIT 15;").WithArgs().WillReturnError(fmt.Errorf("testing error"))

	sut := index.Page(dep)
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

	sqlMock.ExpectQuery("SELECT \\* FROM files ORDER BY uploadDate DESC LIMIT 15;").WithArgs().WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
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

	sut := index.Page(dep)
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
