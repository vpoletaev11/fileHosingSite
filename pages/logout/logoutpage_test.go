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

// TestSuccess checks workability Page()
func TestSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("DELETE").WithArgs("test").WillReturnResult(sqlmock.NewResult(1, 1))

	sut := Page(db)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	inHandlerCookie := &http.Cookie{
		Name:    "session_id",
		Value:   "test",
		Expires: time.Now().Add(30 * time.Minute),
	}

	r.AddCookie(inHandlerCookie)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)
	assert.Equal(t, http.StatusFound, w.Code)
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "<a href=\"/login\">Found</a>.\n\n", bodyString)

	fromHandlerCookie := w.Result().Cookies()
	assert.Equal(t, fromHandlerCookie[0].Name, "session_id")
	assert.Equal(t, fromHandlerCookie[0].MaxAge, -1)
}

// TestDBError checks workability of error handler for database queryer
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

// TestNoCookie checks workability of error handler for cookie handler
func TestNoCookie(t *testing.T) {
	// db is not used in this test
	sut := Page(nil)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, r)
	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, "/login", w.HeaderMap.Get("Location"))
}
