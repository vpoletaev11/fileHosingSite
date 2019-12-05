package dbformat

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLocalTimeSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("example").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	tm, err := userLocalTime(db, time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), "example")
	require.NoError(t, err)

	location, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	assert.Equal(t, time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).In(location), tm)
}

func TestUserLocalUserLocationError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("example").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Unknown/location"))

	tm, err := userLocalTime(db, time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), "example")

	if err != nil {
	} else {
		t.Errorf("Unknown timezone not recognized as error")
	}
	assert.Equal(t, time.Time{}, tm)
}

func TestDownloadFileInfoSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{"id", "label", "filesizeBytes", "description", "owner", "category", "uploadDate", "rating"}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(1, "label", 1024, "description", "owner", "other", time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), 1000))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")
	require.NoError(t, err)

	assert.Equal(t, DownloadFileInfo{DownloadLink: "/files/1", Label: "label", FilesizeMB: "0.000977 MB", Description: "description", Owner: "owner", Category: "other", UploadDate: "2009-11-17 23:34:58", Rating: 1000}, fileInfo)
}

func TestDownloadFileInfoDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnError(fmt.Errorf("testing Error"))
	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, DownloadFileInfo{DownloadLink: "", Label: "", FilesizeMB: "", Description: "", Owner: "", Category: "", UploadDate: "", Rating: 0}, fileInfo)
	assert.Equal(t, fmt.Errorf("testing Error"), err)
}

func TestDownloadFileInfoUserLocalTime(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{"id", "label", "filesizeBytes", "description", "owner", "category", "uploadDate", "rating"}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(1, "label", 1024, "description", "owner", "other", time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), 1000))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnError(fmt.Errorf("testing Error"))

	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, DownloadFileInfo{DownloadLink: "", Label: "", FilesizeMB: "", Description: "", Owner: "", Category: "", UploadDate: "", Rating: 0}, fileInfo)
	assert.Equal(t, fmt.Errorf("testing Error"), err)
}
