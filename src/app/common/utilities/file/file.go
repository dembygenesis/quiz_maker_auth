package file

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/database"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/http"
	"regexp"
)

// GetMultipartFiles validates an array of files received via HTTP.
func GetMultipartFiles(c *fiber.Ctx, fileKey string, acceptEmpty bool) ([]*multipart.FileHeader, error) {

	form, err := c.MultipartForm()

	if err != nil {
		return nil, err
	}

	files := form.File[fileKey]

	if acceptEmpty == false {
		if len(files) == 0 {
			return nil, errors.New("no files in index: " + fileKey)
		}
	}

	return files, nil
}

// Validate file sizes for an array of files
func ValidateFileSizes(f []*multipart.FileHeader) error {
	for _, file := range f {
		if file.Size > int64(config.MaxFileSize) {
			return errors.New(file.Filename + " exceeds file size of 50MB")
		}
	}

	return nil
}

func GetFileContentType(out *bytes.Reader) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetMultiPartAsBuffer(f *multipart.FileHeader) []byte {
	multipartFile, _ := f.Open()

	size := f.Size
	buffer := make([]byte, size)
	multipartFile.Read(buffer)

	return buffer
}

func GetMultiPartFileType(f *multipart.FileHeader) string {

	fileType := http.DetectContentType(GetMultiPartAsBuffer(f))

	reDb := regexp.MustCompile("image/")
	trimmedFileType := reDb.ReplaceAllString(fileType, "")

	return trimmedFileType
}

// GetParsedJSONFileNameAndTypes - takes a string (JSON) and expects it to conform
// to the []FileNameAndTypes struct
func GetParsedJSONFileNameAndTypes(s string) (*[]FileNameAndTypes, error) {
	var f []FileNameAndTypes

	err := json.Unmarshal([]byte(s), &f)

	return &f, err
}

// GetExtensionOfFileNameAsString fetches the filename's (string), and parses out it's extension.
// Example: file.txt -> ".txt"
func GetExtensionOfFileNameAsString(s string) string {
	r, _ := regexp.Compile(`.[^.]*$`)

	result := r.FindAllString(s, -1)

	if len(result) > 0 {
		return result[0]
	} else {
		return ""
	}
}

// GetFileIdHash returns map[string]string
func GetFileIdHash() (*map[int]string, error) {
	fileHash := make(map[int]string)
	type FileHash struct {
		Id   int    `db:"id"`
		Name string `db:"name"`
	}

	var container []FileHash
	var sql string

	sql = `
		SELECT
		  c.id, 
		  c.name
		FROM
		  category_type ct
		  INNER JOIN category c
			ON 1= 1 
			  AND ct.id = c.category_type_ref_id
		WHERE 1 = 1
		  AND ct.name LIKE '%Document%'
	`
	err := database.DBInstancePublic.Select(&container, sql)
	if err != nil {
		return &fileHash, errors.New("error trying to get the file hash: " + err.Error())
	}

	for _, v := range container {
		fileHash[v.Id] = v.Name
	}

	return &fileHash, nil
}

func GetFileNameHash() (map[string]bool, error) {
	fileHash := make(map[string]bool)
	type FileHash struct {
		Id   int    `db:"id"`
		Name string `db:"name"`
	}

	var container []FileHash
	var sql string

	sql = `
		SELECT
		  c.id, 
		  c.name
		FROM
		  category_type ct
		  INNER JOIN category c
			ON 1= 1 
			  AND ct.id = c.category_type_ref_id
		WHERE 1 = 1
		  AND ct.name LIKE '%Document%'
	`
	err := database.DBInstancePublic.Select(&container, sql)
	if err != nil {
		return fileHash, errors.New("error trying to get the file hash: " + err.Error())
	}

	for _, v := range container {
		fileHash[v.Name] = true
	}

	return fileHash, nil
}
