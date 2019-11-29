package admin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/admin", nil)
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
    <title>File Hosting - Admin</title>
    <link rel="stylesheet" href="assets/css/admin.css">
<head>
<body bgcolor=#f1ded3>
    <div class="cookieCleaner">
        <form action="admin" method="post">
            <input type="submit" value="Delete old cookies"> 
        </form>
    </div>
</body>`, bodyString)
}

// TestPageSuccessPost checks workability of POST requests handler in Page()
func TestPageSuccessPOST(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("DELETE FROM sessions WHERE expires").WithArgs(time.Now().Add(-cookieLifetime).Format("2006-01-02 15:04:05")).WillReturnResult(sqlmock.NewResult(1, 1))

	sut := Page(db)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(""))
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
    <title>File Hosting - Admin</title>
    <link rel="stylesheet" href="assets/css/admin.css">
<head>
<body bgcolor=#f1ded3>
    <div class="cookieCleaner">
        <form action="admin" method="post">
            <input type="submit" value="Delete old cookies"> Deleted 1 old cookies
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
	r, err := http.NewRequest(http.MethodGet, "http://localhost/admin", nil)
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
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

// TestPageDeleteCookiesDBError checks workability of error handler of cookie deleter in Page()
func TestPageDeleteCookiesDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("DELETE FROM sessions WHERE expires").WithArgs(time.Now().Add(-cookieLifetime).Format("2006-01-02 15:04:05")).WillReturnError(fmt.Errorf("testing error"))

	sut := Page(db)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "http://localhost/admin", strings.NewReader(""))
	require.NoError(t, err)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
