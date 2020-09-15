package login_test

import (
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/login"
	"github.com/vpoletaev11/fileHostingSite/session"
	"github.com/vpoletaev11/fileHostingSite/test"
)

type anyString struct{}

// ()Match() checks is cookie value valid
func (a anyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if len(v.(string)) != 60 {
		return false
	}
	return true
}

type anyTime struct{}

// ()Match() checks is input value are time
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if len(v.(string)) != 19 {
		return false
	}
	return ok
}

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := login.Page(dep)

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
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            
        </form>
    </div>
</body>`, bodyString)
}

// TestPageSuccessPost checks workability of POST requests handler in Page()
func TestPageSuccessPOST(t *testing.T) {
	dep, sqlMock, redisMock := test.NewDep(t)
	row := []string{"password"}
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows(row).AddRow("$2a$10$ITkHbQjRK6AWs.InpysH5em2Lx4jwzmyYOpvFSturS7hRe6oxzUAu"))
	sqlMock.ExpectExec("INSERT INTO sessions").WithArgs("username", anyString{}, anyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	redisMock.Command("SET", redigomock.NewAnyData(), "username", "EX", session.CookieLifetime.Seconds())

	data := url.Values{}
	data.Set("username", "username")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	assert.Equal(t, http.StatusFound, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)

	assert.Equal(t, "/", w.Header().Get("Location"))

	fromHandlerCookie := w.Result().Cookies()
	assert.Equal(t, fromHandlerCookie[0].Name, "session_id")
	assert.Equal(t, len(fromHandlerCookie[0].Value), 60)
}

// TestPageEmptyUsername tests case when username is empty.
func TestPageEmptyUsername(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	data := url.Values{}
	data.Set("username", "")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Username cannot be empty</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageEmptyPassword tests case when password is empty.
func TestPageEmptyPassword(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Password cannot be empty</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageLargerUsername tests case when len(username) > 20.
func TestPageLargerUsername(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	data := url.Values{}
	data.Set("username", "example_larger_than_20_characters")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Username cannot be longer than 20 characters</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageLargerPassword tests case when len(password) > 20.
func TestPageLargerPassword(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "password_larger_than_40_characters____________________")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Password cannot be longer than 40 characters</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageNonLowerCaseUsername tests case when username is non lower-case
func TestPageNonLowerCaseUsername(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	data := url.Values{}
	data.Set("username", "Example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Please use lower case username</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageQuerySelectErr tests case when SELECT query returns error
func TestPageQuerySELECTErr(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("example").WillReturnError(fmt.Errorf("test error"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">INTERNAL ERROR. Please try later</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageSELECTReturnsEmptyPass tests case when SELECT query returns empty password
func TestPageSELECTReturnsEmptyPass(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"password"}
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("example").WillReturnRows(sqlmock.NewRows(row).AddRow(""))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Wrong username or password</h2>
        </form>
    </div>
</body>`, bodyString)
}

func TestPagePassordNotFound(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("example").WillReturnError(fmt.Errorf("sql: no rows in result set"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Wrong username or password</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageComparePasswordsDoesntMatch tests case when comparePasswords() gets not matched password with hashed password and returns error
func TestPageComparePasswordsDoesntMatch(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"password"}
	sqlMock.ExpectQuery("SELECT password FROM users WHERE username =").WithArgs("example").WillReturnRows(sqlmock.NewRows(row).AddRow("broken hash"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password", "example")

	r, err := http.NewRequest("POST", "http://localhost/login", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := login.Page(dep)
	sut(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" href="assets/css/login.css">
<head>
<body bgcolor=#f1ded3>
    <div class="loginForm">
        <form action="" method="post">
            <p>Username: <input required maxlength="20" type="text" name="username"></p>
            <p>Password: <input required maxlength="40" type="password" name="password"></p>
            <input type="submit" value="Login">
            <p><a href="/registration" style="color: #c82020">Not registered?</a></p>
            <h2 style="color:red">Wrong username or password</h2>
        </form>
    </div>
</body>`, bodyString)
}
