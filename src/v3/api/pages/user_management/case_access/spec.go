package case_access

type RequestCreate2 struct {
	Hello []struct{
		UserRefId int `json:"user_id" validate:"required"`
		CaseRefId int `json:"case_id" validate:"required"`
	} `json:"hello" validate:"required"`
}

// RequestCreate
type RequestCreate struct {
	UserRefId int `json:"user_id" validate:"required"`
	CaseRefId int `json:"case_id" validate:"required"`
}

type RequestDelete struct {
	Id int `json:"id" validate:"required"`
}

type AssigneeDetails struct {
	UserType       string
	OrganizationId int
}
