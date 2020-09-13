package logout

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/test"
)

// TestSuccess checks workability Page()
func TestSuccess(t *testing.T) {
	dep, _, redisMock := test.NewDep(t)
	redisMock.Command("DEL", redigomock.NewAnyData())

	sut := Page(dep)

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
	dep, _, redisMock := test.NewDep(t)
	redisMock.Command("DEL", redigomock.NewAnyData()).ExpectError(fmt.Errorf("Testing error"))

	sut := Page(dep)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "test",
		Expires: time.Now().Add(30 * time.Minute),
	}

	req.AddCookie(cookie)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

// TestNoCookie checks workability of error handler for cookie handler
func TestNoCookie(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	// db is not used in this test
	sut := Page(dep)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, req)
	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, "/login", w.HeaderMap.Get("Location"))
}
