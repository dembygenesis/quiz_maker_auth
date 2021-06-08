package models

type CaseAccess struct {
	Id        int `gorm:"primary_key"`
	UserRefId int
	CaseRefId int
}
