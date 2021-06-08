package crud

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/date"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/gorm_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/interface_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/models/email"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/models"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
)

type Repository interface {
	Create(p *RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError
	List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
		*[]ResponseUserList,
		response_builder.Pagination,
		*error_utils.ApplicationError,
	)
	ListUserTypes(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
		*[]models.UserType,
		response_builder.Pagination,
		*error_utils.ApplicationError,
	)
}

type RepositoryImpl struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &RepositoryImpl{DB: DB}
}

func (r *RepositoryImpl) getUserTypeById(tx *gorm.DB, p *RequestCreate, d *controller_utils.Caller) (string, error) {
	var userType models.UserType
	err := tx.First(&userType, p.UserTypeId).Error
	if err != nil {
		return userType.Name, gorm_utils.GetGormInternalOrNoRecordsFoundError(
			err,
			"user_type_id not found",
			"error trying to fetch the user type id",
		)
	}
	return userType.Name, nil
}

func (r *RepositoryImpl) emailExists(tx *gorm.DB, p *RequestCreate, d *controller_utils.Caller) (bool, error) {
	var user models.User
	err := tx.Where("email = ?", p.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			fmt.Println("--------- eammil exists ---------", err.Error())
			return false, err
		}
	}
	return true, nil
}

// createValidations validates creation of different user types
func (r *RepositoryImpl) createValidations(tx *gorm.DB, p *RequestCreate, d *controller_utils.Caller) error {
	userType, err := r.getUserTypeById(tx, p, d)
	if err != nil {
		return errors.New("error trying to fetch: " + err.Error())
	}
	if userType == "Organization Member" {
		if d.UserType != "Admin" {
			return errors.New("only admins can create a new organization member")
		}
		emailExists, err := r.emailExists(tx, p, d)
		if err != nil {
			return err
		}
		if emailExists == true {
			return errors.New("email is already taken")
		}
	} else {
		return errors.New("adding members aside from organization member is still to do")
	}
	return err
}

// sendWelcomeEmail sends an email to the email specified in the create parameters
// informing of the new password given, along with the login
func (r *RepositoryImpl) sendWelcomeEmail(
	emailAddress string,
	firstName string,
	lastName string,
	password string,
	organization string,
) error {
	var message string
	var subject string

	subject = "Welcome to MedilegalRecords!"

	message = ""
	message = fmt.Sprintf("Hello <i>%v %v,</i><br/><br/>", firstName, lastName)
	message += fmt.Sprintf("You have been given access as an organization member for <i>%v</i>.<br/>", organization)
	message += fmt.Sprintf("You can login at our website:  %v<br/><br/>", config.BaseUrl)
	message += fmt.Sprintf("Credentials <br/>Email: %v<br/>Password: %v<br/><br/>", emailAddress, password)
	message += "Best,<br/>The MediLegalRecords team"

	// Override email address temporarily
	emailAddress = "dembygenesis@gmail.com"

	err := email.SendMail(emailAddress, subject, message)
	if err != nil {
		fmt.Println("==================== FAILED TO SEND MESSAGE: " + err.Error())
	} else {
		fmt.Println("Sent message")
	}

	return nil
}

// createValidations inserts a new user record
func (r *RepositoryImpl) createUser(tx *gorm.DB, p *RequestCreate, d *controller_utils.Caller) error {
	organizationRefId := gorm_utils.ToNullInt64(d.OrganizationId)
	birthday, err := date.StrToDate(p.Birthday)
	if err != nil {
		return errors.New("error converting date")
	}
	password, err := string_utils.Encrypt(p.Password)
	if err != nil {
		return errors.New("error encrypting password: " + err.Error())
	}
	user := models.User{
		FirstName:         p.FirstName,
		LastName:          p.LastName,
		Email:             p.Email,
		MobileNumber:      p.MobileNumber,
		Password:          password,
		UserTypeId:        p.UserTypeId,
		CreatedBy:         d.UserId,
		UpdatedBy:         d.UserId,
		OrganizationRefId: &organizationRefId,
		Address:           p.Address,
		Birthday:          birthday,
		Gender:            p.Gender,
	}
	omits := []string{"CreatedDate", "LastUpdated", "IsActive", "ResetKey"}
	err = tx.Omit(omits...).Create(&user).Error
	if err != nil {
		return err
	}

	var organization models.Organization
	err = tx.Where("id = ?", d.OrganizationId).First(&organization).Error

	// Send password email
	go func() {

		if err != nil {
			fmt.Println("================ errors trying to get organization", err.Error())
			return
		}
		// Probably send email, password, firstname, and lastname
		err = r.sendWelcomeEmail(p.Email, p.FirstName, p.LastName, p.Password, organization.Name)
		if err != nil {
			fmt.Println("================ Error sending email ================")
		}
	}()

	// Send off shit
	return nil
}

// Create inserts a new user if successful validation requirements are met
func (r *RepositoryImpl) Create(p *RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError {
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err = r.createValidations(tx, p, d)
		if err != nil {
			return err
		}

		err = r.createUser(tx, p, d)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return nil
}

// getUsersHasUserTypeFilter checks if there is a user_type parameter
// passed.
func (r *RepositoryImpl) getUsersHasUserTypeFilter(
	p *controller_utils.PaginationDetails,
) (bool, error) {
	if p.Search["user_type"] == nil {
		return false, nil
	}
	val, err := interface_utils.GetJSONValueIfString(p.Search["user_type"])
	if err != nil {
		return false, errors.New("error trying to parse user_type")
	}
	if val == "Organization Member" {
		return true, nil
	}
	return false, nil
}

// getUserTypes - fetches the user types
func (r *RepositoryImpl) getUserTypes(
	tx *gorm.DB,
) (*[]models.UserType, response_builder.Pagination, error) {
	var err error
	var response []models.UserType
	var pagination response_builder.Pagination
	err = tx.Find(&response).Error
	if err != nil {
		return &response, pagination, err
	}
	return &response, pagination, err
}

// getUsers fetches the users
func (r *RepositoryImpl) getUsers(
	tx *gorm.DB,
	d *controller_utils.Caller,
	p *controller_utils.PaginationDetails,
) (*[]ResponseUserList, response_builder.Pagination, error) {
	var q *gorm.DB
	var err error
	var response []ResponseUserList
	var pagination response_builder.Pagination

	selects := []string{
		"id",
		"firstname",
		"lastname",
	}
	q = tx.Where("organization_ref_id = ?", d.OrganizationId)

	if p.Search["user_id"] != nil {
		if reflect.TypeOf(p.Search["user_id"]).String() == "string" {
			userId, err := strconv.Atoi(p.Search["user_id"].(string))
			if err != nil {
				return &response, pagination, err
			}
			q = q.Where("id = ?", userId)
		}
	}

	hasUserTypeFilter, err := r.getUsersHasUserTypeFilter(p)
	if err != nil {

		return &response, pagination, err
	}
	if hasUserTypeFilter == true {

		userType := p.Search["user_type"].(string)
		// Add Where
		q = q.Where("user_type_id = (SELECT id FROM user_type ut WHERE ut.name = ?)", userType)
	} else {

	}

	q = q.Model(models.User{})
	q = q.Select(selects)
	pagination, err = gorm_utils.GetGormPaginatedQuery(q, &response, p.Page, p.Rows, 1000)
	return &response, pagination, err
}

// List displays all the users
func (r *RepositoryImpl) List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]ResponseUserList,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	var users *[]ResponseUserList
	var pagination response_builder.Pagination
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		users, pagination, err = r.getUsers(tx, d, p)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, pagination, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return users, pagination, nil
}


// ListUserTypes displays all the users
func (r *RepositoryImpl) ListUserTypes(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]models.UserType,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	var userTypes *[]models.UserType
	var pagination response_builder.Pagination
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		userTypes, pagination, err = r.getUserTypes(tx)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, pagination, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return userTypes, pagination, nil
}
