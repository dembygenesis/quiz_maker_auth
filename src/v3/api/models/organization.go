package models

type Organization struct {
	Id   int `gorm:"primary_key"`
	Name string
}
