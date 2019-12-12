package cookiescleaner

import (
	"database/sql/driver"
	"io/ioutil"
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

	assert.Equal(t, "deleted 1 cookies at: "+time.Now().Format("2006-01-02 15:04:05")+"\n", bodyString)
}
