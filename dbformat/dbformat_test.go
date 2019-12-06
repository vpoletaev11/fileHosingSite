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

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")
	require.NoError(t, err)

	assert.Equal(t, DownloadFileInfo{
		DownloadLink: "/files/1",
		Label:        "label",
		FilesizeMB:   "0.000977 MB",
		Description:  "description",
		Owner:        "owner",
		Category:     "other",
		UploadDate:   "2009-11-17 23:34:58",
		Rating:       1000,
	}, fileInfo)
}

func TestDownloadFileInfoDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnError(fmt.Errorf("testing Error"))
	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, DownloadFileInfo{
		DownloadLink: "",
		Label:        "",
		FilesizeMB:   "",
		Description:  "",
		Owner:        "",
		Category:     "",
		UploadDate:   "",
		Rating:       0,
	}, fileInfo)
	assert.Equal(t, fmt.Errorf("testing Error"), err)
}

func TestDownloadFileInfoUserLocalTimeError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnError(fmt.Errorf("testing Error"))

	fileInfo, err := FormatedDownloadFileInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, DownloadFileInfo{
		DownloadLink: "",
		Label:        "",
		FilesizeMB:   "",
		Description:  "",
		Owner:        "",
		Category:     "",
		UploadDate:   "",
		Rating:       0,
	}, fileInfo)
	assert.Equal(t, fmt.Errorf("testing Error"), err)
}

func TestFormatedFilesInfoSuccess(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")
	require.NoError(t, err)

	assert.Equal(t, []FileInfo{FileInfo{
		Label:                "label",
		DownloadLink:         "/download?id=1",
		FilesizeMb:           "0.0010 MB",
		Description:          "description",
		Owner:                "owner",
		Category:             "other",
		UploadDate:           "2009-11-17 23:34:58",
		Rating:               1000,
		LabelComment:         "label",
		FilesizeBytesComment: "1024 Bytes",
		DescriptionComment:   "description",
	}}, fileInfo)
}

func TestFormatedFilesInfoSuccessLongFilename(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label_longer_than_20_characters",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")
	require.NoError(t, err)

	assert.Equal(t, []FileInfo{FileInfo{
		Label:                "label_longer_than_20...",
		DownloadLink:         "/download?id=1",
		FilesizeMb:           "0.0010 MB",
		Description:          "description",
		Owner:                "owner",
		Category:             "other",
		UploadDate:           "2009-11-17 23:34:58",
		Rating:               1000,
		LabelComment:         "label_longer_than_20_characters",
		FilesizeBytesComment: "1024 Bytes",
		DescriptionComment:   "description",
	}}, fileInfo)
}

func TestFormatedFilesInfoSuccessLongDescription(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description_longer_than_35_characters",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"timezone"}).AddRow("Europe/Moscow"))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")
	require.NoError(t, err)

	assert.Equal(t, []FileInfo{FileInfo{
		Label:                "label",
		DownloadLink:         "/download?id=1",
		FilesizeMb:           "0.0010 MB",
		Description:          "description_longer_than_35_characte...",
		Owner:                "owner",
		Category:             "other",
		UploadDate:           "2009-11-17 23:34:58",
		Rating:               1000,
		LabelComment:         "label",
		FilesizeBytesComment: "1024 Bytes",
		DescriptionComment:   "description_longer_than_35_characters",
	}}, fileInfo)
}

func TestFormatedFilesInfoDBError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnError(fmt.Errorf("testing Error"))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, fmt.Errorf("testing Error"), err)
	assert.Equal(t, []FileInfo{}, fileInfo)
}

func TestFormatedFilesInfoRowsNextError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	if err != nil {
	} else {
		t.Errorf("rows.Scan don't recognize incorrect variables types")
	}
	assert.Equal(t, []FileInfo{}, fileInfo)
}

func TestFileInfoUserLocalTimeError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	fileInfoRows := []string{
		"id",
		"label",
		"filesizeBytes",
		"description",
		"owner",
		"category",
		"uploadDate",
		"rating",
	}

	sqlMock.ExpectQuery("SELECT \\* FROM files WHERE id =").WithArgs("1").WillReturnRows(sqlmock.NewRows(fileInfoRows).AddRow(
		1,
		"label",
		1024,
		"description",
		"owner",
		"other",
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		1000,
	))
	sqlMock.ExpectQuery("SELECT timezone FROM users WHERE username =").WithArgs("username").WillReturnError(fmt.Errorf("testing Error"))

	fileInfo, err := FormatedFilesInfo("username", db, "SELECT * FROM files WHERE id = ?;", "1")

	assert.Equal(t, []FileInfo{}, fileInfo)
	assert.Equal(t, fmt.Errorf("testing Error"), err)
}
