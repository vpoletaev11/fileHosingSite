package download

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil, "username")

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/download", nil)
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
    <title>Categories</title>
    <link rel="stylesheet" href="/assets/css/categories.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/upload">Upload file</a></li>
            <li><a href="/">Home</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <ul class="categoriesList">
        <li><a href="/categories/other" class="categoryLink">Other</a></li>
        <li><a href="/categories/games" class="categoryLink">Games</a></li>
        <li><a href="/categories/documents" class="categoryLink">Documents</a></li>
        <li><a href="/categories/projects" class="categoryLink">Projects</a></li>
        <li><a href="/categories/music" class="categoryLink">Music</a></li>
    </ul>
</body>`, bodyString)
}
