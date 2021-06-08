package case_access

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/gorm_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/models"
	"gorm.io/gorm"
	"net/http"
)

type Repository interface {
	Create(p *[]RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError
	Delete(p *[]RequestDelete, d *controller_utils.Caller) *error_utils.ApplicationError
}

type RepositoryImpl struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &RepositoryImpl{DB: DB}
}

// getCaseAccessByUserIdAndCaseId fetches a case access entry via "user_id" and "case_id"
func (r *RepositoryImpl) getCaseAccessByUserIdAndCaseId(tx *gorm.DB, userId int, caseId int, d *controller_utils.Caller) (*models.CaseAccess, error) {
	var caseAccess models.CaseAccess
	err := tx.Where("user_ref_id = ?", userId).
		Where("case_ref_id = ?", caseId).
		First(&caseAccess).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// It's fine
		} else {
			return nil, err
		}
	}
	return &caseAccess, nil
}

// getAssigneeDetails returns "user_id" details: "user_type" and "organization_id"
func (r *RepositoryImpl) getAssigneeDetails(tx *gorm.DB, userId int, d *controller_utils.Caller) (*AssigneeDetails, error) {
	var assigneeDetails AssigneeDetails
	selects := []string{
		"ut.name AS user_type",
		"u.organization_ref_id AS organization_id",
	}
	err := tx.Debug().Table("user u").
		Select(selects).
		Joins("JOIN user_type ut ON u.user_type_id = ut.id").
		Where("u.id = ?", userId).
		Where("u.is_active = ?", 1).Scan(&assigneeDetails).Error
	if err != nil {
		return &assigneeDetails, errors.New(fmt.Sprintf("error fetching assignee details: %v", err.Error()))
	}
	if assigneeDetails.OrganizationId == 0 {
		return &assigneeDetails, errors.New(fmt.Sprintf("no assignee details found for user_id: %v", userId))
	}
	return &assigneeDetails, err
}

// validateAdminAccessViaCaseIdAndUserId ensure caller can access the "case_id" and "user_id"
func (r *RepositoryImpl) validateAdminAccessViaCaseIdAndUserId(tx *gorm.DB, userId int, caseId int, d *controller_utils.Caller) error {
	if d.UserType != "Admin" {
		return errors.New("only administrators can assign new members")
	}

	assigneeDetails, err := r.getAssigneeDetails(tx, userId, d)
	if err != nil{
		return err
	}
	if assigneeDetails.UserType != "Organization Member" {
		return errors.New("user_type must be Organization Member")
	}
	if d.OrganizationId != assigneeDetails.OrganizationId {
		return errors.New("assignee belongs to another organization that you don't have access to")
	}
	return err
}

// validateCreate ensures create parameters are correct
func (r *RepositoryImpl) validateCreate(tx *gorm.DB, p *[]RequestCreate, d *controller_utils.Caller) error {
	for _, v := range *p {
		caseAccess, err := r.getCaseAccessByUserIdAndCaseId(tx, v.UserRefId, v.CaseRefId, d)
		if err != nil{
			return err
		}
		if caseAccess.Id != 0 {
			return errors.New("case_access entry already exists")
		}

		err = r.validateAdminAccessViaCaseIdAndUserId(tx, v.UserRefId, v.CaseRefId, d)
		if err != nil {
			return err
		}
	}
	return nil
}

// createAccess performs the insert operation given a "case_id" and "user_id"
func (r *RepositoryImpl) createAccess(tx *gorm.DB, p *[]RequestCreate, d *controller_utils.Caller) error {
	var caseAccesses []models.CaseAccess
	for _, v := range *p {
		caseAccesses = append(caseAccesses, models.CaseAccess{
			UserRefId: v.UserRefId,
			CaseRefId: v.CaseRefId,
		})
	}
	return tx.Omit("Id").Create(&caseAccesses).Error
}

// validateDelete ensures access conditions and variables are correct
func (r *RepositoryImpl) validateDelete(tx *gorm.DB, p *[]RequestDelete, d *controller_utils.Caller) error {
	for _, v := range *p {
		var caseAccess models.CaseAccess
		err := tx.Where("id = ?", v.Id).First(&caseAccess).Error
		if err != nil {
			return gorm_utils.GetGormInternalOrNoRecordsFoundError(err,
				fmt.Sprintf("no case_access entry found given the id: %v", v.Id),
				err.Error(),
			)
		}
		err = r.validateAdminAccessViaCaseIdAndUserId(tx,
			caseAccess.UserRefId,
			caseAccess.CaseRefId,
			d,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteAccess performs the delete operation given a case access id
func (r *RepositoryImpl) deleteAccess(tx *gorm.DB, p *[]RequestDelete, d *controller_utils.Caller) error {
	var ids []int
	for _, v := range *p {
		ids = append(ids, v.Id)
	}
	return tx.Delete(models.CaseAccess{}, ids).Error
}

// Create inserts a case access entry
func (r *RepositoryImpl) Create(p *[]RequestCreate, d *controller_utils.Caller) *error_utils.ApplicationError {
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err = r.validateCreate(tx, p, d)
		if err != nil {
			return err
		}

		err = r.createAccess(tx, p, d)
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

// Delete deletes a case access entry
func (r *RepositoryImpl) Delete(p *[]RequestDelete, d *controller_utils.Caller) *error_utils.ApplicationError {
	var err error
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err = r.validateDelete(tx, p, d)
		if err != nil {
			return err
		}

		err = r.deleteAccess(tx, p, d)
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