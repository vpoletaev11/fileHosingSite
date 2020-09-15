package test

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/session"
)

// NewDep returns dependencies for Page()'s and sqlMock and redisMock interfaces to writting mocks
func NewDep(t *testing.T) (session.Dependency, sqlmock.Sqlmock, *redigomock.Conn) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	redisMock := redigomock.NewConn()

	return session.Dependency{Db: db, Redis: redisMock, Username: "username"}, sqlMock, redisMock
}

// AssertBodyEqual checks if responce body equal expected value
func AssertBodyEqual(t *testing.T, expected string, actual *bytes.Buffer) {
	bodyBytes, err := ioutil.ReadAll(actual)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	assert.Equal(t, expected, bodyString)
}
