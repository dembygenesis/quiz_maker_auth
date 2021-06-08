package case_access

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
)

type Service interface {
	Create(p *[]RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError
	Delete(p *[]RequestDelete, d *controller_utils.Caller) *error_utils.ApplicationError
}

type ServiceImpl struct {
	Repository Repository
}

func NewService(rep Repository) Service {
	return &ServiceImpl{Repository: rep}
}

func (s *ServiceImpl) Create(p *[]RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError {
	err := controller_utils.ValidateRequestParams(p)
	if err != nil {
		return err
	}
	return s.Repository.Create(p, d)
}

func (s *ServiceImpl) Delete(p *[]RequestDelete, d *controller_utils.Caller) *error_utils.ApplicationError {
	err := controller_utils.ValidateRequestParams(p)
	if err != nil {
		return err
	}
	return s.Repository.Delete(p, d)
}