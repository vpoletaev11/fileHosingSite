package cookie_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/cookie"
)

const (
	cookieVal = "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu"
	username  = "username"
)

type anyString struct{}

func testHandler(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, username)
	}
}

func TestCreateCookieSuccess(t *testing.T) {
	cookie1, err := cookie.CreateCookie("user")
	assert.NoError(t, err)
	cookie2, err := cookie.CreateCookie("user")
	assert.NoError(t, err)

	if cookie1.Expires.After(time.Now().Add(30*time.Minute + 1*time.Second)) {
		t.Error("cookie.Expires > 30 min. cookie.Expires = " + cookie1.Expires.String())
	}
	if cookie1.Expires.Before(time.Now().Add(30*time.Minute - 1*time.Second)) {
		t.Error("cookie.Expires < 30 min. cookie.Expires = " + cookie1.Expires.String())
	}

	if len(cookie1.Value) != 60 {
		t.Errorf("CreateCookie() creates cookie with invalid cookie.Value (len(cookie.Value) < 60)")
	}

	if cookie1.Value == cookie2.Value {
		t.Errorf("CreateCookie() creates not unique cookie value")
	}

	if cookie1.Name != "session_id" {
		t.Errorf("CreateCookie() creates cookie with invalid cookie.Name (cookie.Name != \"session_id\"). Cookie.Name == " + cookie1.Name)
	}
}

func TestAuthWrapperSuccess(t *testing.T) {
	redisConn := redigomock.NewConn()
	redisConn.Command("GET", cookieVal).Expect(username)
	redisConn.Command("EXPIRE", cookieVal, cookie.CookieLifetime.Seconds())

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: cookieVal,
	}
	r.AddCookie(inHandlerCookie)

	sut := cookie.AuthWrapper(testHandler, nil, redisConn)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, username, bodyString)
}

func TestAuthWrapperEmptyCookieError(t *testing.T) {
	redisConn := redigomock.NewConn()

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	inHandlerCookie := &http.Cookie{}
	r.AddCookie(inHandlerCookie)

	sut := cookie.AuthWrapper(testHandler, nil, redisConn)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "", bodyString)
}

func TestAuthWrapperGettingUsernameError(t *testing.T) {
	redisConn := redigomock.NewConn()
	redisConn.Command("GET", cookieVal).ExpectError(fmt.Errorf("Testing Error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: cookieVal,
	}

	r.AddCookie(inHandlerCookie)

	sut := cookie.AuthWrapper(testHandler, nil, redisConn)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "/login", w.Header().Get("Location"))
	assert.Equal(t, "", bodyString)
}

func TestAuthWrapperExtendingCookieLifetimeError(t *testing.T) {
	redisConn := redigomock.NewConn()
	redisConn.Command("GET", cookieVal).Expect(username)
	redisConn.Command("EXPIRE", cookieVal, cookie.CookieLifetime.Seconds()).ExpectError(fmt.Errorf("Testing error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: cookieVal,
	}
	r.AddCookie(inHandlerCookie)

	sut := cookie.AuthWrapper(testHandler, nil, redisConn)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later.\n", bodyString)
}
