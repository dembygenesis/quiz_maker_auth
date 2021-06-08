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
	Gender         string `json:"gender" required:"true"`
}

type ResponseUserList struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstname" gorm:"column:firstname"`
	LastName  string `json:"lastname" gorm:"column:lastname"`
}
