package controller_utils

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/validation_utils"
	"net/http"
	"strconv"
)

// Caller stores the caller's general info
type Caller struct {
	UserId         int
	UserType       string
	OrganizationId int
	LawFirmId      int
}

// PaginationDetails stores pagination variables
type PaginationDetails struct {
	DisabledPagination int `json:"disabled_pagination"`
	Rows               int `json:"rows"`
	Page               int `json:"page"`
	Search             map[string]interface{}
}

// GetPaginationDetails returns the caller info
func GetPaginationDetails(c *fiber.Ctx) *PaginationDetails {
	var disabledPagination int
	var search map[string]interface{}
	var page int
	var rows int

	// Set default rows to 100 if not paginated
	if c.Query("page") == "" {
		page = 0
	} else {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("rows") == "" {
		rows = 200
	} else {
		rows, _ = strconv.Atoi(c.Query("rows"))

		if rows <= 0 {
			rows = 1000
		} else if rows < 50 {
			// TODO: restore this
			// rows = 50
		} else if rows > 500 {
			rows = 500
		}
	}

	if c.Query("disabled_pagination") == "1" {
		disabledPagination = 1
	} else {
		disabledPagination = 0
	}

	if c.Query("search") != "" {
		err := json.Unmarshal([]byte(c.Query("search")), &search)
		if err != nil {
			fmt.Println("--------Error parsing search JSON--------", err)
		}
	}

	return &PaginationDetails{
		DisabledPagination: disabledPagination,
		Rows:               rows,
		Page:               page,
		Search:             search,
	}
}

// GetCallerDetails returns the caller info
func GetCallerDetails(c *fiber.Ctx) *Caller {
	return &Caller{
		UserId:         c.Locals("tokenExtractedUserId").(int),
		UserType:       c.Locals("tokenExtractedUserType").(string),
		OrganizationId: c.Locals("tokenExtractedOrganizationId").(int),
		LawFirmId:      c.Locals("tokenExtractedLawFirmId").(int),
	}
}

// ValidateRequestParams is a helper function for validating
// struct rules determined in the struct's tags
func ValidateRequestParams(s interface{}) *error_utils.ApplicationError {
	validationErrors, err := validation_utils.ValidateStructParams2(s)
	if err != nil {
		return &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "error attempting to validate the assumed struct parameter passed",
			Error:      err,
		}
	}
	if len(validationErrors) > 0 {
		return &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "validation errors",
			Error:      validationErrors,
		}
	}
	return nil
}

// RespondError - returns an error formatted JSON using Fiber.Ctx
func RespondError(c *fiber.Ctx, operationStatus string, apiErr *error_utils.ApplicationError) error {
	r := response_builder.Response{
		HttpCode:        apiErr.HttpStatus,
		ResponseMessage: apiErr.Message,
		Data:            apiErr.Error,
		OperationStatus: operationStatus,
		Pagination:      nil,
	}
	r.SetErrors(apiErr.Error)
	return c.Status(apiErr.HttpStatus).JSON(r)
}

// Respond - returns a result formatted JSON using Fiber.Ctx
func Respond(c *fiber.Ctx, operationStatus string, responseMessage string, data interface{}) error {
	r := response_builder.Response{
		HttpCode:        http.StatusOK,
		ResponseMessage: responseMessage,
		Data:            data,
		OperationStatus: operationStatus,
	}
	return c.Status(http.StatusOK).JSON(r)
}

// RespondWithPagination - returns a result formatted JSON using Fiber.Ctx with the pagination library
func RespondWithPagination(
	c *fiber.Ctx,
	operationStatus string,
	responseMessage string,
	data interface{},
	pagination *response_builder.Pagination,
) error {
	r := response_builder.Response{
		HttpCode:        http.StatusOK,
		ResponseMessage: responseMessage,
		Data:            data,
		OperationStatus: operationStatus,
		Pagination:      pagination,
	}
	return c.Status(http.StatusOK).JSON(r)
}
