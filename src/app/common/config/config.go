package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

/**
This holds all the static variables
*/

var (
	IsDev         int
	BaseUrl       string
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
)

var InsertFailed = "INSERT_FAILED"
var FetchSuccess = "FETCH_SUCCESS"
var InsertSuccess = "INSERT_SUCCESS"
var UpdateSuccess = "UPDATE_SUCCESS"
var DeleteSuccess = "DELETE_SUCCESS"
var DeleteFailed = "DELETE_FAILED"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	IsDev, err = strconv.Atoi(os.Getenv("IS_DEV"))
	BaseUrl = os.Getenv("BASE_URL")
	EmailHost = os.Getenv("EMAIL_HOST")
	EmailPort = os.Getenv("EMAIL_PORT")
	EmailUsername = os.Getenv("EMAIL_USERNAME")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
}

// Notification type
var ValidNotificationTypes = map[string]bool{
	"Recently Created":         true,
	"Document Type Added":      true,
	"Document Type Deleted":    true,
	"Document Type Updated":    true,
	"Case Status Updated":      true,
	"Treatment Status Updated": true,
}

// Email
var From = "MedilegalRecords <support@medilegalrecords.com>"

// Resolve
var CaseResolveTypes = []string{"Resolve", "Re-open"}
var CaseResolveCategoryType = "Resolution Status"

// Files
var MaxFileSize = 50000000 // 50MB
var BucketPath = "medilegalrecords.com"
var DocumentTypeOther = "Other"

// ValidDocuments are all the registered document types
var ValidDocuments = map[string]bool{
	"Document Type":                          true,
	"Document Type HIPAA Release":            true,
	"Document Type Signed Lien Letter":       true,
	"Document Type Liability":                true,
	"Document Type Letter Of Request":        true,
	"Document Type Policy Coverage":          true,
	"Document Type Change of Representation": true,
}

// SingularAttachedDocuments are document types where the user can only upload 1 at a time
var SingularAttachedDocuments = map[string]bool{
	"Document Type Signed Lien Letter": true,
	"Document Type Policy Coverage":    true,
	"Document Type HIPAA Release":      true,
}
var RestrictedDocumentType = "Document Type"

// Policy
var PolicyTypeCategory = "Document Type Policy Coverage"
var PolicyTypeLowerLimitExclusion = "MedPay"

var CaseStatusDoneTreating = "Done Treating"
