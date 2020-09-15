package logout_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/pages/logout"
	"github.com/vpoletaev11/fileHostingSite/test"
)

// TestSuccess checks workability Page()
func TestSuccess(t *testing.T) {
	dep, _, redisMock := test.NewDep(t)
	redisMock.Command("DEL", redigomock.NewAnyData())

	sut := logout.Page(dep)

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
	test.AssertBodyEqual(t, "<a href=\"/login\">Found</a>.\n\n", w.Body)

	fromHandlerCookie := w.Result().Cookies()
	assert.Equal(t, fromHandlerCookie[0].Name, "session_id")
	assert.Equal(t, fromHandlerCookie[0].MaxAge, -1)
}

// TestDBError checks workability of error handler for database queryer
func TestDBError(t *testing.T) {
	dep, _, redisMock := test.NewDep(t)
	redisMock.Command("DEL", redigomock.NewAnyData()).ExpectError(fmt.Errorf("Testing error"))

	sut := logout.Page(dep)

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
	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
}

// TestNoCookie checks workability of error handler for cookie handler
func TestNoCookie(t *testing.T) {
	dep, _, _ := test.NewDep(t)
	
	sut := logout.Page(dep)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut(w, req)
	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, "/login", w.HeaderMap.Get("Location"))
}
