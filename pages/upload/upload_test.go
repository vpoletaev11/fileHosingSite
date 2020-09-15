package upload_test

import (
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/upload"
	"github.com/vpoletaev11/fileHostingSite/test"
)

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

func TestPageSuccessGET(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := upload.Page(dep)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/upload", nil)
	require.NoError(t, err)

	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageSuccessPOST(t *testing.T) {
	// changing directory because of test are not containing in root folder
	os.Chdir("../../")
	defer os.Chdir("pages/upload")

	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO files").WithArgs(
		"filename",
		11,
		"description",
		"username",
		"other",
		anyTime{},
	).WillReturnResult(sqlmock.NewResult(1, 1))

	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:green">FILE SUCCEEDED UPLOADED</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageErrorFileReceptionPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	// changing directory because of test are not containing in root folder
	os.Chdir("../../")
	defer os.Chdir("pages/upload")

	sqlMock.ExpectExec("INSERT INTO files").WithArgs(
		"filename",
		11,
		"description",
		"username",
		"other",
		anyTime{},
	).WillReturnResult(sqlmock.NewResult(1, 1))

	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
}

func TestPageEmptyFilenameSuccessPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	// changing directory because of test are not containing in root folder
	os.Chdir("../../")
	defer os.Chdir("pages/upload")

	sqlMock.ExpectExec("INSERT INTO files").WithArgs(
		"file",
		11,
		"description",
		"username",
		"other",
		anyTime{},
	).WillReturnResult(sqlmock.NewResult(1, 1))

	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"


--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file";
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:green">FILE SUCCEEDED UPLOADED</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageLargeFilenameErrorPOST(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename_larger_than_50_characters_filename_larger_than_50_characters
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file";
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:red">Filename are too long</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageLargeDescriptionErrorPOST(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_charactersdescription_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_description_larger_than_500_characters_
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file";
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:red">Description are too long</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageWrongCategoryErrorPOST(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

unknown
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file";
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:red">Unknown category</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageDBInsertionErrorPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO files").WithArgs(
		"filename",
		11,
		"description",
		"username",
		"other",
		anyTime{},
	).WillReturnError(fmt.Errorf("testing error"))

	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload file</title>
    <link rel="stylesheet" href="assets/css/upload.css">
<head>
<body bgcolor=#f1ded3>
    <div class="menu">
        <ul class="nav">
            <li><a href="/">Home</a></li>
            <li><a href="/categories">Categories</a></li>
            <li><a href="/popular">Most popular</a></li>
            <li><a href="/users">Users</a></li>
            <li><a href="/logout">Logout</a></li>
        </ul>
    </div>
    <div class="username">Welcome, username</div>

    <div class="uploadFormBox">
        <div class="uploadFormContent">
        <form action="" method="post" enctype="multipart/form-data">
            <p>Filename: <input type="text" maxlength="50" name="filename"></p><br>
            <p>Input description for uploading file:</p>
            <textarea cols="80" rows="15" maxlength="500" name="description"></textarea>
    
            <p>Category: <select name="category">
                <option selected="selected" value="other">other</option>
                <option value="games">games</option>
                <option value="documents">documents</option>
                <option value="projects">projects</option>
                <option value="music">music</option>
                </select></p>
                   
            <p><input required type="file" name="uploaded_file"></input></p>

            <p><input type="submit" value="UPLOAD"></p>
            <h2 style="color:red">INTERNAL ERROR. Please try later</h2>
        </form>
        </div>
    </div>
</body>`, w.Body)
}

func TestPageCreatingFileErrorPOST(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectExec("INSERT INTO files").WithArgs(
		"filename",
		11,
		"description",
		"username",
		"other",
		anyTime{},
	).WillReturnResult(sqlmock.NewResult(1, 1))

	postData :=
		`--xxx
Content-Disposition: form-data; name="filename"

filename
--xxx
Content-Disposition: form-data; name="description"

description
--xxx
Content-Disposition: form-data; name="category"

other
--xxx
Content-Disposition: form-data; name="uploaded_file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	r := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	w := httptest.NewRecorder()

	sut := upload.Page(dep)
	sut(w, r)

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
}
