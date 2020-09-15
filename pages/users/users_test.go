package users_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/users"
	"github.com/vpoletaev11/fileHostingSite/test"
)

func TestPageSuccessGet(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15").WithArgs().WillReturnRows(sqlmock.NewRows(
		[]string{
			"username",
			"rating",
		}).AddRow(
		"user",
		1000,
	))

	sut := users.Page(dep)
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

func TestPageDBQueryErrorsGet(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)

	sqlMock.ExpectQuery("SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15").WithArgs().WillReturnError(fmt.Errorf("testing error"))

	sut := users.Page(dep)
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
