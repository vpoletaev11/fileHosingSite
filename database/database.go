package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// FileInfo contains processed file info from fileInfoDB{}
type FileInfo struct {
	Label        string
	DownloadLink string
	FilesizeMb   string
	Description  string
	Owner        string
	Category     string
	UploadDate   string
	Rating       int

	LabelComment         string
	FilesizeBytesComment string
	DescriptionComment   string
}

// FormatedFilesInfo returns array of formatted file information
func FormatedFilesInfo(rows *sql.Rows) ([]FileInfo, error) {
	var fiTableCollection []FileInfo
	fiTable := new(FileInfo)

	id := 0
	var uploadDateTime time.Time
	for rows.Next() {
		err := rows.Scan(
			&id,
			&fiTable.LabelComment,
			&fiTable.FilesizeBytesComment,
			&fiTable.DescriptionComment,
			&fiTable.Owner,
			&fiTable.Category,
			&uploadDateTime,
			&fiTable.Rating,
		)
		if err != nil {
			return []FileInfo{}, err
		}
		fiTable.UploadDate = uploadDateTime.Format("2006-01-02 15:04:05")

		if len(fiTable.LabelComment) > 20 {
			fiTable.Label = fiTable.LabelComment[:20] + "..."
		} else {
			fiTable.Label = fiTable.LabelComment
		}

		if len(fiTable.DescriptionComment) > 25 {
			fiTable.Description = fiTable.DescriptionComment[:25] + "..."
		} else {
			fiTable.Description = fiTable.DescriptionComment
		}

		fsBytes, err := strconv.Atoi(fiTable.FilesizeBytesComment)
		if err != nil {
			return []FileInfo{}, err
		}
		fiTable.FilesizeMb = fmt.Sprintf("%.4f", float64(fsBytes)/1024/1024) + " MB"
		fiTable.DownloadLink = "/download?id=" + strconv.Itoa(id)
		fiTable.FilesizeBytesComment = fiTable.FilesizeBytesComment + " Bytes"

		fiTableCollection = append(fiTableCollection, *fiTable)
	}
	return fiTableCollection, nil
}
