package login

import (
	"database/sql/driver"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type anyString struct{}

// ()Match() checks is cookie value valid
func (a anyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if !(len(v.(string)) == 60) {
		return false
	}
	return true
}

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/login", nil)
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
    <title>login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required type="text" name="username"></p>
            <p>Password: <input required type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/register" style="color:yellow">Not registered?</a></p>
               
        </form>
    </div>
</body>`, bodyString)
}

// TestPageSuccessPost checks workability of POST requests handler in Page()
func TestPageSuccessPOST(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	row := []string{"password"}
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("example").WillReturnRows(sqlmock.NewRows(row).AddRow("$2a$10$ITkHbQjRK6AWs.InpysH5em2Lx4jwzmyYOpvFSturS7hRe6oxzUAu"))
	sqlMock.ExpectExec("INSERT INTO sessions").WithArgs("example", anyString{}).WillReturnResult(sqlmock.NewResult(1, 1))

	apiURL := "http://localhost/login"

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode())) // URL-encoded payload
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	assert.Equal(t, http.StatusFound, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)

	fromHandlerCookie := w.Result().Cookies()
	assert.Equal(t, fromHandlerCookie[0].Name, "session_id")
	assert.Equal(t, len(fromHandlerCookie[0].Value), 60)
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
	r, err := http.NewRequest(http.MethodGet, "http://localhost/login", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(nil)
	sut(w, r)

	// renaming template file to original filename
	oldName = newName
	newName = oldName[:lenOrigName]
	// todo defer
	err = os.Rename(oldName, newName)
	require.NoError(t, err)

	// checking error handler works correct
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Internal error. Page not found\n", bodyString)
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
