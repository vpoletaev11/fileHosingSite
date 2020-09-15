package download_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/download"
	"github.com/vpoletaev11/fileHostingSite/test"
)

func TestPageSuccessGET(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
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

	sut := download.Page(dep)

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

func TestPageSettingRatingSuccessPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"owner",
		}).AddRow(
			"owner",
		))
	sqlMock.ExpectExec("UPDATE users SET rating").WithArgs(10, "owner").WillReturnResult(sqlmock.NewResult(1, 1))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)
}

func TestPageUpdatingRatingSuccessPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"0",
		))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(0, 10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE filesRating SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"owner",
		}).AddRow(
			"owner",
		))
	sqlMock.ExpectExec("UPDATE users SET rating").WithArgs(0, 10, "owner").WillReturnResult(sqlmock.NewResult(1, 1))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)
}

func TestPageUpdatingRatingSameRatingSuccessPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"10",
		))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)
}

func TestPageUpdatingRatingsDBError01POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageUpdatingRatingsDBError02POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"0",
		))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(0, 10, "1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageUpdatingRatingsDBError03POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"0",
		))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(0, 10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE filesRating SET rating").WithArgs(10, "1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageUpdatingRatingsDBError04POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"0",
		))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(0, 10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE filesRating SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageUpdatingRatingsDBError05POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("Error 1062"))
	sqlMock.ExpectQuery("SELECT rating FROM filesRating WHERE fileID").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"rating",
		}).AddRow(
			"0",
		))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(0, 10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE filesRating SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"owner",
		}).AddRow(
			"owner",
		))
	sqlMock.ExpectExec("UPDATE users SET rating").WithArgs(0, 10, "owner").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDBFileInfoGatheringErrorGET(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

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
	dep, sqlMock, _ := test.NewDep(t)
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

	sut := download.Page(dep)

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

func TestPageIncorrectPOSTParameter01(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "wrongParameter")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INCORRECT POST PARAMETER\n", bodyString)
}

func TestPageIncorrectPOSTParameter02(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "11")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INCORRECT POST PARAMETER\n", bodyString)
}

func TestPageIncorrectPOSTParameter03(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "-11")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INCORRECT POST PARAMETER\n", bodyString)
}

func TestPageSetRatingError01POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageSetRatingError02POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(10, "1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageSetRatingError03POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageSetRatingError04POST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO filesRating").WithArgs("1", "username", 10).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec("UPDATE files SET rating").WithArgs(10, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectQuery("SELECT owner FROM files WHERE id").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{
			"owner",
		}).AddRow(
			"owner",
		))
	sqlMock.ExpectExec("UPDATE users SET rating").WithArgs(10, "owner").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("rating", "10")
	r, err := http.NewRequest(http.MethodPost, "http://localhost/download?id=1", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
