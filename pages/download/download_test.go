package download

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
)

func TestPageSuccessGET(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"label",
			"filesizeBytes",
			"description",
			"owner",
			"category",
			"uploadDate",
			"rating",
		}).AddRow(
			1,
			"label",
			1000,
			"description",
			"owner",
			"other",
			time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			100,
		))

	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username").WithArgs("username").WillReturnRows(
		sqlmock.NewRows([]string{
			"timezone",
		}).AddRow(
			"Europe/Moscow",
		))

	sut := Page(db, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/download?id=1", nil)
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
    <title>Download</title>
    <link rel="stylesheet" href="assets/css/download.css">
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

    <div class="fileInfo">
        <div class="filename"><h2>Filename: label</h2></div>
        <div class="filesize"><h2>Filesize: 0.000954 MB</h2></div>
        <div class="description"><h2>Description: description</h2></div>
        <div class="owner"><h2>Owner: owner</h2></div>
        <div class="category"><h2>Category: other</h2></div>
        <div class="uploadDate"><h2>Upload date: 2009-11-17 23:34:58</h2></div>
        <div class="rating"><h2>Rating: 100</h2></div>

        <div class="setRating">
            <form  action="" method="post">
                    Set rating:
                            <select name="rating">
                                <option value="-10">-10</option>
                                <option value="-9">-9</option>
                                <option value="-8">-8</option>
                                <option value="-7">-7</option>
                                <option value="-6">-6</option>
                                <option value="-5">-5</option>
                                <option value="-4">-4</option>
                                <option value="-3">-3</option>
                                <option value="-2">-2</option>
                                <option value="-1">-1</option>
                                <option value="1">1</option>
                                <option value="2">2</option>
                                <option value="3">3</option>
                                <option value="4">4</option>
                                <option value="5">5</option>
                                <option value="6">6</option>
                                <option value="7">7</option>
                                <option value="8">8</option>
                                <option value="9">9</option>
                                <option value="10">10</option>
                            </select>
                <input type="submit" value="VOTE">
            </form>
        </div>

        <div class="download">
            <a href="/files/1" download=><h1>download</h1></a>
        </div>
    </div>
</body>`, bodyString)
}

func TestPageMissingTemplate(t *testing.T) {
	// renaming exists template file
	oldName := "../../" + pathTemplateDownload
	newName := "../../" + pathTemplateDownload + "edit"
	err := os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(nil, "username")
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
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDBFileInfoGatheringErrorGET(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := Page(db, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/download?id=1", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDBFTimezoneGatheringErrorGET(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"label",
			"filesizeBytes",
			"description",
			"owner",
			"category",
			"uploadDate",
			"rating",
		}).AddRow(
			1,
			"label",
			1000,
			"description",
			"owner",
			"other",
			time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			100,
		))

	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username").WithArgs("username").WillReturnError(fmt.Errorf("testing error"))

	sut := Page(db, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/download?id=1", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
