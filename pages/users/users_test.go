package users

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/test"
)

func TestPageSuccessGet(t *testing.T) {
	dep, sqlMock := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15").WithArgs().WillReturnRows(sqlmock.NewRows(
		[]string{
			"username",
			"rating",
		}).AddRow(
		"user",
		1000,
	))

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
    <title>Users</title>
    <link rel="stylesheet" href="assets/css/users.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/upload">Upload file</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/">Home</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="label">
        <br><br><br><br><br>
        <p><h1>↓↓↓ TOP UPLOADERS ↓↓↓</h1></p>
    </div>

    <div class = "userList">
        <table border="1" width="100%" cellpadding="5">
            <tr>
                <th>Username</th>
                <th>User rating</th>
            </tr>
            
            <tr>
                <td width="70%">user</td>
                <td width="30%">1000</td>
            </tr>
            
        </table>
    </div>
</body>`, bodyString)
}

func TestPageMissingTemplate(t *testing.T) {
	dep, _ := test.NewDep(t)
	// renaming exists template file
	oldName := "../../" + pathTemplateUsers
	newName := "../../" + pathTemplateUsers + "edit"
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

func TestPageDBQueryErrorsGet(t *testing.T) {
	dep, sqlMock := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15").WithArgs().WillReturnError(fmt.Errorf("testing error"))

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
