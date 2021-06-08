package quiz

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/gorm_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/interface_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/models"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/modelsV2"
	"gorm.io/gorm"
	"net/http"
)

type Repository interface {
	List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
		*[]ResponseCases,
		response_builder.Pagination,
		*error_utils.ApplicationError,
	)
	FetchQuiz(id int) (
		*modelsV2.Quiz,
		response_builder.Pagination,
		*error_utils.ApplicationError,
	)
	AnswerQuiz(p *[]RequestAnswerQuiz) (
		*error_utils.ApplicationError,
	)
}

type RepositoryImpl struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &RepositoryImpl{DB: DB}
}

// hasUserManagementAccessFilter returns a bool if the user has access
func (r *RepositoryImpl) hasUserManagementAccessFilter(p *controller_utils.PaginationDetails) (bool, error) {
	// If params for including case access for a specific user are specified
	if p.Search["management_access"] != nil && p.Search["user_id"] != nil {
		managementAccess, err := interface_utils.GetJSONValueIfString(p.Search["management_access"])
		if err != nil {
			return false, errors.New("error when trying to parse the management_access: " + err.Error())
		}

		userId, err := interface_utils.GetJSONValueIfInt(p.Search["user_id"])
		if err != nil {
			return false, errors.New("error when trying to parse the user_id: " + err.Error())
		}

		if userId != 0 && managementAccess == "1" {
			return true, nil
		}
	}
	return false, nil
}

// getCases fetches the cases (now at v2, using Gorm)
func (r *RepositoryImpl) getCases(
	tx *gorm.DB,
	d *controller_utils.Caller,
	p *controller_utils.PaginationDetails,
) (*[]ResponseCases, response_builder.Pagination, error) {
	var q *gorm.DB
	var err error
	var response *[]ResponseCases
	var pagination response_builder.Pagination
	selects := []string{
		"case.id",
		"case.patient_first_name",
		"case.patient_last_name",
		"UNIX_TIMESTAMP(case.created_date) AS created_date",
		"cat.name AS treatment_status",
	}
	q = tx.Where("organization_ref_id = ?", d.OrganizationId)
	q = q.Model(&models.Case{}).
		Joins("JOIN category cat ON case.treatment_status_ref_id = cat.id")
	// Filter: User management
	hasUserManagementAccessFilter, err := r.hasUserManagementAccessFilter(p)
	if err != nil {
		return response, pagination, err
	}
	if hasUserManagementAccessFilter == true {
		fmt.Println("============= geometric benefits =============")

		userId, _ := interface_utils.GetJSONValueIfInt(p.Search["user_id"])
		selects = append(selects, "IF(ca.id IS NULL, 0, 1) AS has_access")
		selects = append(selects, "IF(ca.id IS NULL, 0, ca.id) AS case_access_id")
		q = q.Joins(`
			LEFT JOIN case_access ca
				ON case.id = ca.case_ref_id
					AND ca.user_ref_id = ?
		`, userId)
	} else {
		fmt.Println("============= no geometric benefits =============")
	}
	q = q.Select(selects)
	pagination, err = gorm_utils.GetGormPaginatedQuery(q, &response, p.Page, p.Rows, 1000)
	return response, pagination, err
}

// Fetch
func (r *RepositoryImpl) FetchQuiz(id int) (
	*modelsV2.Quiz,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	var pagination response_builder.Pagination
	var err error
	var quiz *modelsV2.Quiz
	err = r.DB.Where("id = ?", id).Preload("QuizQuestions.QuizChoices").Find(&quiz).Error
	if err != nil {
		return nil, pagination, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return quiz, pagination, nil
}

// List displays all the users
func (r *RepositoryImpl) List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]ResponseCases,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	var cases *[]ResponseCases
	var pagination response_builder.Pagination
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		cases, pagination, err = r.getCases(tx, d, p)
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
	return cases, pagination, nil
}

// validateQuizAnswers ensures the questions asked matches the database records
func (r *RepositoryImpl) validateQuizAnswers(tx *gorm.DB, p *[]RequestAnswerQuiz) error {
	var c int
	var i []int
	l := len(*p)
	if l == 0 {
		return errors.New("answers provided are empty")
	}
	for _, v := range *p {
		i = append(i, v.QuizQuestionId)
	}
	err := tx.Model(&modelsV2.QuizQuestion{}).
		Where("id IN ?", i).
		Select("count(ID)").
		Scan(&c).Error
	if err != nil {
		return err
	}
	if c != l {
		return errors.New("quiz questions provided do not match the database records")
	}
	return nil
}

// evaluateQuizAnswers checks the number of correct answers
func (r *RepositoryImpl) evaluateQuizAnswers(tx *gorm.DB, p *[]RequestAnswerQuiz) error {
	for _, v := range *p {
		var r modelsV2.QuizQuestion
		err := tx.Where("id = ?", v.QuizQuestionId).Find(&r).Error
		if err != nil {
			return err
		}
		if r.Answer == v.Answer {
			fmt.Println("Correct!")
		} else {
			fmt.Println("Incorrect!")
		}
	}
	return nil
}

func (r *RepositoryImpl) AnswerQuiz(p *[]RequestAnswerQuiz) *error_utils.ApplicationError {
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err = r.validateQuizAnswers(tx, p)
		if err != nil {
			return err
		}
		return r.evaluateQuizAnswers(tx, p)
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