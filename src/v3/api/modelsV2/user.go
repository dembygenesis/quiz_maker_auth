package modelsV2

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserTypeID int
	FirstName  string
	Email      string
	Password   string
	LastName   string
	Gender     string
	Address    string
	UserType   UserType
}
