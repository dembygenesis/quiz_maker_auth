package models

type Case struct {
	Id                int `gorm:"primary_key"`
	PatientFirstName  string
	PatientLastName   string
	CreatedDate       int
	OrganizationRefId int
}
