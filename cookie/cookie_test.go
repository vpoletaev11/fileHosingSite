package cookie

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type anyString struct{}

// ()Match() checks is cookie value valid
func (a anyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if !(len(v.(string)) == 60) {
		return false
	}
	return true
}

type anyTime struct{}

// ()Match() checks is time valid
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if len(v.(string)) != 19 {
		return false
	}
	return true
}

func TestCreateCookie(t *testing.T) {
	cookie1 := CreateCookie()
	cookie2 := CreateCookie()

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

func TestCookieValidatorSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(30*time.Minute), "example"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	username, cookie, err := cookieValidator(db, r)
	require.NoError(t, err)

	assert.Equal(t, "example", username)
	assert.Equal(t, "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu", cookie.Value)
}

func TestCookieValidatorEmptyCookie(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)

	inHandlerCookie := &http.Cookie{}

	r.AddCookie(inHandlerCookie)

	username, cookie, err := cookieValidator(nil, r)
	require.NoError(t, err)

	assert.Equal(t, "", username)
	assert.Equal(t, http.Cookie{}, cookie)
}

func TestCookieValidatorErrGetExpiresAndUsername(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WillReturnError(fmt.Errorf("testing error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	username, cookie, err := cookieValidator(db, r)

	assert.Equal(t, "", username)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, "testing error", err.Error())
}

func TestCookieValidatorCookieExpired(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(-30*time.Minute), "example"))
	sqlMock.ExpectExec("DELETE FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnResult(sqlmock.NewResult(1, 1))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	username, cookie, err := cookieValidator(db, r)
	require.NoError(t, err)

	assert.Equal(t, "", username)
	assert.Equal(t, http.Cookie{}, cookie)
}

func TestCookieValidatorCookieExpiredErrorDB(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(-30*time.Minute), "example"))
	sqlMock.ExpectExec("DELETE FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnError(fmt.Errorf("testing error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	username, cookie, err := cookieValidator(db, r)

	assert.Equal(t, "", username)
	assert.Equal(t, http.Cookie{}, cookie)
	assert.Equal(t, "testing error", err.Error())
}

func testHandler(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, username)
	}
}

func TestAuthWrapperSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(30*time.Minute), "example"))
	sqlMock.ExpectExec("UPDATE sessions SET expires=").WithArgs(anyTime{}, anyString{}).WillReturnResult(sqlmock.NewResult(1, 1))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	sut := AuthWrapper(testHandler, db)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "example", bodyString)
}

func TestAuthWrapperValidatorError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WillReturnError(fmt.Errorf("testing error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	sut := AuthWrapper(testHandler, db)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later.\n", bodyString)
}

func TestAuthWrapperEmptyUsename(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(30*time.Minute), ""))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	sut := AuthWrapper(testHandler, db)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "/login", w.Header().Get("Location"))
	assert.Equal(t, "", bodyString)
}

func TestAuthWrapperExtendingCookieLifetimeDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	rows := []string{"expires", "username"}
	sqlMock.ExpectQuery("SELECT expires, username FROM sessions WHERE cookie=").WithArgs(anyString{}).WillReturnRows(sqlmock.NewRows(rows).AddRow(time.Now().Add(30*time.Minute), "example"))
	sqlMock.ExpectExec("UPDATE sessions SET expires=").WithArgs(anyTime{}, anyString{}).WillReturnError(fmt.Errorf("testing error"))

	r, err := http.NewRequest(http.MethodPost, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	inHandlerCookie := &http.Cookie{
		Name:  "session_id",
		Value: "D8SgghMYJQSo9PXuH7wihJlrRFP18RKBzITHDwXou8VGqaVHW1Yi9KWyIrUu",
	}

	r.AddCookie(inHandlerCookie)

	sut := AuthWrapper(testHandler, db)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "INTERNAL ERROR. Please try later.\n", bodyString)
}
