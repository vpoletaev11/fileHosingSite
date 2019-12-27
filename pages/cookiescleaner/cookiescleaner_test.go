package cookiescleaner

import (
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestPageSuccessPOST(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT password, timezone FROM users WHERE username").WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{
			"password",
			"timezone",
		}).AddRow(
			"$2a$10$7d2k15x6gBqEt7hA4HJ4JuIF9dNZVcDwExP/xnBWWXsoeO88b8CY6",
			"Europe/Moscow",
		))
	sqlMock.ExpectExec("DELETE FROM sessions WHERE expires").WithArgs(anyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	data := `{"username":"admin","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "deleted 1 cookies at: "+time.Now().In(location).Format("2006-01-02 15:04:05")+"\n", bodyString)
}

func TestPageSuccessJsonError(t *testing.T) {
	data := "skidaddle skidoodle"
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageWrongUsername(t *testing.T) {
	data := `{"username":"wrongUsername","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "Wrong username or password\n", bodyString)
}

func TestPagePassAndTimezoneQueryError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT password, timezone FROM users WHERE username").WithArgs("admin").WillReturnError(fmt.Errorf("testing error"))

	data := `{"username":"admin","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageComparePasswordsError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT password, timezone FROM users WHERE username").WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{
			"password",
			"timezone",
		}).AddRow(
			"wrongPass",
			"Europe/Moscow",
		))

	data := `{"username":"admin","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "Wrong username or password\n", bodyString)
}

func TestPageDeleteCookieError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT password, timezone FROM users WHERE username").WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{
			"password",
			"timezone",
		}).AddRow(
			"$2a$10$7d2k15x6gBqEt7hA4HJ4JuIF9dNZVcDwExP/xnBWWXsoeO88b8CY6",
			"Europe/Moscow",
		))
	sqlMock.ExpectExec("DELETE FROM sessions WHERE expires").WithArgs(anyTime{}).WillReturnError(fmt.Errorf("testing error"))

	data := `{"username":"admin","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

func TestPageDeletedAtError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT password, timezone FROM users WHERE username").WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{
			"password",
			"timezone",
		}).AddRow(
			"$2a$10$7d2k15x6gBqEt7hA4HJ4JuIF9dNZVcDwExP/xnBWWXsoeO88b8CY6",
			"WrongTimezone",
		))
	sqlMock.ExpectExec("DELETE FROM sessions WHERE expires").WithArgs(anyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	data := `{"username":"admin","password":"password"}`
	r, err := http.NewRequest("POST", "http://localhost/cookiescleaner", strings.NewReader(data))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}
