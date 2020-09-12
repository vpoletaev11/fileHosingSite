package test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/require"
	"github.com/vpoletaev11/fileHostingSite/session"
)

// NewDep returns dependencies for pages and sqlmock interface to writting mocks
func NewDep(t *testing.T) (session.Dependency, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	redisConn := redigomock.NewConn()

	return session.Dependency{Db: db, Redis: redisConn, Username: "username"}, sqlMock
}
