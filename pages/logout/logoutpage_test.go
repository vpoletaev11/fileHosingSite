package logout

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	sqlMock.ExpectExec("DELETE").WithArgs("test").WillReturnResult(sqlmock.NewResult(1, 1))

	require.NoError(t, err)

	sut := Page(db)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "test",
		Expires: time.Now().Add(30 * time.Minute),
	}

	r.AddCookie(cookie)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)
	assert.Equal(t, http.StatusFound, w.Code)
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "<a href=\"/login\">Found</a>.\n\n", bodyString)
	// todo: test cookie
}

func TestDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	sqlMock.ExpectExec("DELETE").WithArgs("test").WillReturnError(errors.New("DB Error"))

	require.NoError(t, err)

	sut := Page(db)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "test",
		Expires: time.Now().Add(30 * time.Minute),
	}

	r.AddCookie(cookie)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "DB Error\n", bodyString)
}

func TestNoCookie(t *testing.T) {
	// db is not used in this test
	sut := Page(nil)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)
	assert.Equal(t, w.Code, http.StatusFound)
	//todo: check redirect url
}
