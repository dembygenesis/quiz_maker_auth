package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/file"
	"mime/multipart"
	"net/url"
	"os"
	"sync"
	"time"
)

var l = sync.Mutex{}

var (
	s3Session *s3.S3
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load .env file")
	}

	GetS3Instance()
}

func GetS3Instance() {
	accessId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	awsRegion := os.Getenv("AWS_REGION")

	fmt.Println("awsRegion", awsRegion)

	s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			accessId,
			secretKey,
			"",
		),
	})))
}

// Delete file to AWS
func DeleteObject(filename string, bucketPath string) error {

	o := s3.DeleteObjectInput{
		Bucket: aws.String(bucketPath),
		Key:    aws.String(filename),
	}

	_, err := s3Session.DeleteObject(&o)

	if err != nil {
		fmt.Println("error delete object", err)
		return err
	} else {
		fmt.Println("No error deleting the s3 item")
	}

	return nil
}

// Delete file to AWS (has additional parameters for errors and go routines)
func DeleteObjectGoRoutine(filename string, bucketPath string, errors *[]error, wg *sync.WaitGroup) {
	defer wg.Done()

	o := s3.DeleteObjectInput{
		Bucket: aws.String(bucketPath),
		Key:    aws.String(filename),
	}

	_, err := s3Session.DeleteObject(&o)

	if err != nil {
		fmt.Println("error delete object", err)

		l.Lock()
		*errors = append(*errors, err)
		l.Unlock()

	} else {
		fmt.Println("No error deleting the s3 item")
	}

}

// Upload file to AWS
func UploadObjectMultiPart(filename string, f *multipart.FileHeader, bucketPath string) error {
	_, err := s3Session.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(file.GetMultiPartAsBuffer(f)),
		Bucket:      aws.String(bucketPath),
		ContentType: aws.String(file.GetMultiPartFileType(f)),
		Key:         aws.String(filename),
		// ACL: aws.String("public-read"),
	})

	if err != nil {
		fmt.Println("error uploading to s3", err)
		return err
	} else {
		fmt.Println("No error uploading to s3!")
	}

	return nil
}

// GetS3PreSignedURL - Return a link to either view or download an S3 file
func GetS3PreSignedURL(s string, fileName string) (string, error) {
	bucket := os.Getenv("AWS_BUCKET")

	var preSignedUrl string

	o := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s),
	}

	// If file name is provided, execute DOWNLOAD
	if fileName != "" {
		o.ResponseContentDisposition = aws.String(`attachment; filename=` + fileName + ``)
		o.ResponseContentType = aws.String("application/octet-stream")
	}

	req, _ := s3Session.GetObjectRequest(o)

	preSignedUrl, err := req.Presign(15 * time.Minute)

	if err != nil {
		return preSignedUrl, err
	}

	return preSignedUrl, nil
}

// CopyObject - takes a key of the object to move, and it's new destination path
// Also tweaked to add errors to a map pointer, and wg.Done - because "Go".
func CopyObject(
	source string,
	dest string,
	bucketPath string,
	errors *[]error,
	wg *sync.WaitGroup,
) error {
	_, err := s3Session.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucketPath),
		CopySource: aws.String(url.PathEscape(bucketPath + "/" + source)),
		Key:        aws.String(dest),
	})

	if err != nil {
		fmt.Println("error copying to s3", err)

		l.Lock()
		*errors = append(*errors, err)
		l.Unlock()

		return err
	} else {
		fmt.Println("No error copying to s3!")
	}

	if wg != nil {
		defer wg.Done()
	}

	return nil
}
