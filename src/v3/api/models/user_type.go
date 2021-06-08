package models

type UserType struct {
	Id   int    `json:"id" db:"id" gorm:"primary_key"`
	Name string `json:"name" db:"name"`
}
