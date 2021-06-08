package crud

type RequestCreate struct {
	FirstName      string `validate:"required,min=3,max=32" json:"firstname" name:"firstname"`
	LastName       string `validate:"required,min=3,max=32" json:"lastname" required:"true"`
	Email          string `json:"email" validate:"email" required:"true"`
	Password       string `json:"password" required:"true"`
	UserTypeId     int    `json:"user_type_id" validate:"required"`
	MobileNumber   string `json:"mobile_number" validate:"required"`
	OrganizationId int    `json:"organization" required:"false"`
	Address        string `json:"address" validate:"required"`
	Birthday       string `json:"birthday" validate:"required,date_format" required:"false"`
	Gender         string `json:"gender" Wrequired:"true"`
}

type ResponseCases struct {
	Id               int    `json:"id" db:"id"`
	PatientFirstName string `json:"patient_first_name" db:"patient_first_name"`
	PatientLastName  string `json:"patient_last_name" db:"patient_last_name"`
	TreatmentStatus  string `json:"treatment_status" db:"treatment_status"`
	CreatedDate      int    `json:"created_date" db:"created_date"`
	HasAccess        int    `json:"has_access" db:"has_access"`
	CaseAccessId     int    `json:"case_access_id" db:"case_access_id"`
}
