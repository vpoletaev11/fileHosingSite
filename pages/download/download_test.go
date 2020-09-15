package download_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

	test.AssertBodyEqual(t, `<!doctype html>
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
</body>`, w.Body)
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

	test.AssertBodyEqual(t, "", w.Body)
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

	test.AssertBodyEqual(t, "", w.Body)
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

	test.AssertBodyEqual(t, "", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
}

func TestPageDBFileInfoGatheringErrorGET(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id").WithArgs("1").WillReturnError(fmt.Errorf("testing error"))

	sut := download.Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/download?id=1", nil)
	require.NoError(t, err)

	sut(w, r)

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INCORRECT POST PARAMETER\n", w.Body)
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

	test.AssertBodyEqual(t, "INCORRECT POST PARAMETER\n", w.Body)
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

	test.AssertBodyEqual(t, "INCORRECT POST PARAMETER\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
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

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
}
