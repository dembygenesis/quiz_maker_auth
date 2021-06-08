package crud

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/models"
)

type Service interface {
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

type ServiceImpl struct {
	Repository Repository
}

func NewService(rep Repository) Service {
	return &ServiceImpl{Repository: rep}
}

func (s*ServiceImpl) Create(p *RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError {
	err := controller_utils.ValidateRequestParams(p)
	if err != nil {
		return err
	}
	return s.Repository.Create(p, d)
}

func (s*ServiceImpl) List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]ResponseUserList,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	return s.Repository.List(d, p)
}

func (s*ServiceImpl) ListUserTypes(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]models.UserType,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	return s.Repository.ListUserTypes(d, p)
}

