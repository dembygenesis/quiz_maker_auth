package crud

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Handler interface {
	CreateUser(c *fiber.Ctx) error
	ListUsers(c *fiber.Ctx) error
	ListUserTypes(c *fiber.Ctx) error
}

type HandlerImpl struct {
	Service Service
}

func NewHandler(service Service) Handler {
	return &HandlerImpl{Service: service}
}

func (h *HandlerImpl) CreateUser(c *fiber.Ctx) error {
	var body RequestCreate

	err := c.BodyParser(&body)
	if err != nil {
		return controller_utils.RespondError(c, "INSERT_FAILED", &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	appError := h.Service.Create(&body, controller_utils.GetCallerDetails(c))
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	return controller_utils.Respond(c, config.InsertSuccess, "Successfully added the user", nil)
}

func (h *HandlerImpl) ListUsers(c *fiber.Ctx) error {
	users, pagination, appError := h.Service.List(controller_utils.GetCallerDetails(c), controller_utils.GetPaginationDetails(c))
	fmt.Println("pagination", pagination)
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	return controller_utils.RespondWithPagination(c, config.FetchSuccess, "Successfully fetched the user", users, &pagination)
}

func (h *HandlerImpl) ListUserTypes(c *fiber.Ctx) error {
	users, pagination, appError := h.Service.ListUserTypes(controller_utils.GetCallerDetails(c), controller_utils.GetPaginationDetails(c))
	fmt.Println("pagination", pagination)
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	return controller_utils.RespondWithPagination(c, config.FetchSuccess, "Successfully fetched the user", users, &pagination)
}


