package case_access

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	utils "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Handler interface {
	Create(c* fiber.Ctx) error
	Delete(c* fiber.Ctx) error
}

type HandlerImpl struct {
	Service Service
}

func NewHandler(Service Service) Handler {
	return &HandlerImpl{Service: Service}
}

func (h *HandlerImpl) Create(c* fiber.Ctx) error {
	var err error
	var body []RequestCreate

	err = c.BodyParser(&body)
	if err != nil {
		return utils.RespondError(c, "INSERT_FAILED", &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}

	appError := h.Service.Create(&body, utils.GetCallerDetails(c))
	if appError != nil {
		return utils.RespondError(c, config.InsertFailed, appError)
	}

	return utils.Respond(c, config.InsertSuccess, "Successfully added the case access", "Successfully added the case access")
}

func (h *HandlerImpl) Delete(c* fiber.Ctx) error {
	var body []RequestDelete
	err := c.BodyParser(&body)
	if err != nil {
		return utils.RespondError(c, config.DeleteFailed, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	appError := h.Service.Delete(&body, utils.GetCallerDetails(c))
	if appError != nil {
		return utils.RespondError(c, config.InsertFailed, appError)
	}
	return utils.Respond(c, config.DeleteSuccess, "Successfully deleted the case access", "Successfully deleted the case access")
}