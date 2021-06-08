package crud

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
)

type Service interface {
	List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
		*[]ResponseCases,
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

func (s*ServiceImpl) List(d *controller_utils.Caller, p *controller_utils.PaginationDetails) (
	*[]ResponseCases,
	response_builder.Pagination,
	*error_utils.ApplicationError,
) {
	return s.Repository.List(d, p)
}