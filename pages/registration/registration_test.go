package registration

import (
	"database/sql/driver"
	"fmt"
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
	"golang.org/x/crypto/bcrypt"
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
	r, err := http.NewRequest(http.MethodGet, "http://localhost/registration", nil)
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
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            
        </form>
    </div>
</body>`, bodyString)
}

// TestPageSuccessPost checks workability of POST requests handler in Page()
func TestPageSuccessPost(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}).WillReturnResult(sqlmock.NewResult(1, 1))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	assert.Equal(t, "/login", w.Header().Get("Location"))
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
	r, err := http.NewRequest(http.MethodGet, "http://localhost/registration", nil)
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

// TestPageEmptyUsername tests case when username is empty.
func TestPageEmptyUsername(t *testing.T) {
	data := url.Values{}
	data.Set("username", "")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Username cannot be empty</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageEmptyPassword1 tests case when password1 is empty.
func TestPageEmptyPassword1(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Password cannot be empty</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageEmptyPassword2 tests case when password2 is empty.
func TestPageEmptyPassword2(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Password cannot be empty</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageLargerUsername tests case when len(username) > 20.
func TestPageLargerUsername(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example_larger_than_20_characters")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Username cannot be longer than 20 characters</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageLargerPassword1 tests case when len(password1) > 20.
func TestPageLargerPassword1(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example_larger_than_20_characters")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Password cannot be longer than 20 characters</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageLargerPassword2 tests case when len(password2) > 20.
func TestPageLargerPassword2(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example_larger_than_20_characters")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Password cannot be longer than 20 characters</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageNonLowerCaseUsername tests case when username is non lower-case
func TestPageNonLowerCaseUsernam(t *testing.T) {
	data := url.Values{}
	data.Set("username", "Example")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Please use lower case username</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageMismatchingPasswords tests case when password1 != password2
func TestPageMismatchingPasswords(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example1")
	data.Add("password2", "example2")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Passwords doesn't match</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageNotUniqueUsername tests case when username already exists in MySQL database
func TestPageNotUniqueUsername(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}).WillReturnError(fmt.Errorf("Error 1062"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">Username already used</h2>
        </form>
    </div>
</body>`, bodyString)
}

// TestPageUsernameInsertionDBInternalError tests case when username insertion in MySQL database unreachable because of internal error.
func TestPageUsernameInsertionDBInternalError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}).WillReturnError(fmt.Errorf("Internal error"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>login</title>
    <link rel="stylesheet" href="assets/css/register.css">
<head>
<body bgcolor=#f1ded3>
    <div class="registerForm">
        <form action="" method="post">
            <p>Create username: <input required type="text" name="username"></p>
            <p>Create password: <input required type="password" name="password1"></p>
            <p>Repeat password: <input required type="password" name="password2"></p>
            <p><input type="submit" value="Register"></p>
            <h2 style="color:red">INTERNAL ERROR. Please try later</h2>
        </form>
    </div>
</body>`, bodyString)
}

func TestHashAndSaltSuccess(t *testing.T) {
	plainPass := "password"
	hashedPass, err := hashAndSalt(plainPass)
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	require.NoError(t, err)
}
