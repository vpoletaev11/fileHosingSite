package test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaeljusto/redigomock"
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
