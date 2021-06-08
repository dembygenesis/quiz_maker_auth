package models

import (
	"database/sql"
	"time"
)

// User schema
type User struct {
	Id                int    `gorm:"primary_key"`
	FirstName         string `gorm:"column:firstname"`
	LastName          string `gorm:"column:lastname"`
	Email             string `gorm:"column:email"`
	MobileNumber      string
	Password          string
	UserTypeId        int
	CreatedBy         int
	CreatedDate       *sql.NullTime
	LastUpdated       *sql.NullTime
	UpdatedBy         int
	IsActive          *sql.NullInt64
	OrganizationRefId *sql.NullInt64
	LawFirmRefId      *sql.NullInt64
	Address           string
	Birthday          time.Time
	Gender            string
	ResetKey          string
}
