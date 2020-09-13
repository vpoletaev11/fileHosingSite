package categories

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/test"
)

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/", nil)
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

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageAnyCategorySuccessGET(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 0, 15).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
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


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
    </div>
</body>`, bodyString)
}

func TestPageAnyCategoryFewPagesInPageBarSuccess(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(rowsInPage * 3))

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 0, 15).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
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


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
        <a href=/categories/other?p&#61;1>1</a>
        
        <a href=/categories/other?p&#61;2>2</a>
        
        <a href=/categories/other?p&#61;3>3</a>
        
    </div>
</body>`, bodyString)
}

func TestPageAnyCategoryAlotPagesInPageBarSuccess(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(rowsInPage * 30))

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 0, 15).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
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


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
        <a href=/categories/other?p&#61;1>1</a>
        
        <a href=/categories/other?p&#61;2>2</a>
        
        <a href=/categories/other?p&#61;3>3</a>
        
        <a href=/categories/other?p&#61;4>4</a>
        
        <a href=/categories/other?p&#61;5>5</a>
        
        <a href=/categories/other?p&#61;6>6</a>
        
        <a href=/categories/other?p&#61;7>7</a>
        
        <a href=/categories/other?p&#61;8>8</a>
        
        <a href=/categories/other?p&#61;9>9</a>
        
        <a href=/categories/other?p&#61;10>10</a>
        
        <a href=/categories/other?p&#61;11>11</a>
        
        <a href=/categories/other?p&#61;12>12</a>
        
        <a href=/categories/other?p&#61;13>13</a>
        
        <a href=/categories/other?p&#61;14>14</a>
        
        <a href=/categories/other?p&#61;15>15</a>
        
        <a href=/categories/other?p&#61;16>16</a>
        
        <a href=/categories/other?p&#61;17>17</a>
        
        <a href=/categories/other?p&#61;18>18</a>
        
        <a href=/categories/other?p&#61;19>19</a>
        
        <a href=/categories/other?p&#61;20>20</a>
        
        <a href=/categories/other?p&#61;21>21</a>
        
        <a href=/categories/other?p&#61;22>22</a>
        
        <a href=/categories/other?p&#61;23>23</a>
        
        <a href=/categories/other?p&#61;24>24</a>
        
        <a href=/categories/other?p&#61;25>25</a>
        
        <a href=/categories/other?p&#61;30>30</a>
        
    </div>
</body>`, bodyString)
}

func TestPageAnyCategoryAlotPagesInPageBarDefaultCaseSuccess(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(rowsInPage * 30))

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 15*rowsInPage, 16*rowsInPage).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other?p=16", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
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


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
        <a href=/categories/other?p&#61;1>1</a>
        
        <a href=/categories/other?p&#61;11>11</a>
        
        <a href=/categories/other?p&#61;12>12</a>
        
        <a href=/categories/other?p&#61;13>13</a>
        
        <a href=/categories/other?p&#61;14>14</a>
        
        <a href=/categories/other?p&#61;15>15</a>
        
        <a href=/categories/other?p&#61;16>16</a>
        
        <a href=/categories/other?p&#61;17>17</a>
        
        <a href=/categories/other?p&#61;18>18</a>
        
        <a href=/categories/other?p&#61;19>19</a>
        
        <a href=/categories/other?p&#61;20>20</a>
        
        <a href=/categories/other?p&#61;21>21</a>
        
        <a href=/categories/other?p&#61;22>22</a>
        
        <a href=/categories/other?p&#61;23>23</a>
        
        <a href=/categories/other?p&#61;24>24</a>
        
        <a href=/categories/other?p&#61;25>25</a>
        
        <a href=/categories/other?p&#61;26>26</a>
        
        <a href=/categories/other?p&#61;27>27</a>
        
        <a href=/categories/other?p&#61;28>28</a>
        
        <a href=/categories/other?p&#61;29>29</a>
        
        <a href=/categories/other?p&#61;30>30</a>
        
    </div>
</body>`, bodyString)
}

func TestPageAnyCategoryAlotPagesInPagesBarNumPage1Success(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(rowsInPage * 30))

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 10*rowsInPage, 11*rowsInPage).WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other?p=11", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// html text uses spaces instead of tabs
	assert.Equal(t, `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>other</title>
    <link rel="stylesheet" href="/assets/css/anyCategory.css">
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


    <div class = "newlyUploadedBox">
                    <table border="1" width="100%" cellpadding="5">
                        <tr>
                            <th>Filename</th>
                            <th>Filesize</th>
                            <th>Description</th>
                            <th>Owner</th>
                            <th>Upload date</th>
                            <th>Rating</th>
                        </tr>
                        
                        <tr>
                            <td width="15%" title=label><a href=/download?id&#61;1>label</a></td>
                            <td width="10%" title=1024&#32;Bytes>0.0010 MB</td>
                            <td width="25%" title=description>description</td>
                            <td width="15%">owner</td>
                            <td width="15%">2009-11-17 23:34:58</td>
                            <td width="10%">1000</td>
                        </tr>
                        
                    </table>
        </div>
    </div>

    <div class="pagesNums">
        
        <a href=/categories/other?p&#61;1>1</a>
        
        <a href=/categories/other?p&#61;6>6</a>
        
        <a href=/categories/other?p&#61;7>7</a>
        
        <a href=/categories/other?p&#61;8>8</a>
        
        <a href=/categories/other?p&#61;9>9</a>
        
        <a href=/categories/other?p&#61;10>10</a>
        
        <a href=/categories/other?p&#61;11>11</a>
        
        <a href=/categories/other?p&#61;12>12</a>
        
        <a href=/categories/other?p&#61;13>13</a>
        
        <a href=/categories/other?p&#61;14>14</a>
        
        <a href=/categories/other?p&#61;15>15</a>
        
        <a href=/categories/other?p&#61;16>16</a>
        
        <a href=/categories/other?p&#61;17>17</a>
        
        <a href=/categories/other?p&#61;18>18</a>
        
        <a href=/categories/other?p&#61;19>19</a>
        
        <a href=/categories/other?p&#61;20>20</a>
        
        <a href=/categories/other?p&#61;21>21</a>
        
        <a href=/categories/other?p&#61;22>22</a>
        
        <a href=/categories/other?p&#61;23>23</a>
        
        <a href=/categories/other?p&#61;24>24</a>
        
        <a href=/categories/other?p&#61;25>25</a>
        
        <a href=/categories/other?p&#61;26>26</a>
        
        <a href=/categories/other?p&#61;30>30</a>
        
    </div>
</body>`, bodyString)
}

// TestPageMissingTemplate tests case when template file is missing.
// Cannot be runned in parallel.
func TestPageMissingTemplate(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	// renaming exists template file
	oldName := "../../" + pathTemplateAnyCategory
	newName := "../../" + pathTemplateAnyCategory + "edit"
	err := os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(dep)
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

func TestPageAnyCategoryWrongCategory(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/wrongCategory", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	t.Log(bodyString)
	assert.Equal(t, "ERROR: Incorrect category\n", bodyString)
}

func TestPageAnyCategoryPagesCountError(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnError(fmt.Errorf("testing error"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageAnyCategoryWrongPage(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other?p=wrongPage", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "ERROR: Incorrect get request\n", bodyString)
}

func TestPageAnyCategoryWrongPageLowerThanZero(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other?p=-1", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "ERROR: Incorrect get request\n", bodyString)
}

func TestPageAnyCategoryNumPageBiggerThanPagesCount(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other?p=100", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "ERROR: Incorrect get request\n", bodyString)
}

func TestPageAnyCategorySuccessFileInfoGatheringError(t *testing.T) {
	dep, sqlMock, _ := test.NewDep(t)
	row := []string{"count"}
	sqlMock.ExpectQuery("SELECT COUNT").WithArgs("other").WillReturnRows(sqlmock.NewRows(row).AddRow(1))

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE category =").WithArgs("other", 0, 15).WillReturnError(fmt.Errorf("testing error"))

	sut := Page(dep)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/categories/other", nil)
	require.NoError(t, err)

	sut(w, r)
	err = sqlMock.ExpectationsWereMet()
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
