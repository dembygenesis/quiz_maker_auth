package quiz

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/error_utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type Handler interface {
	ListCases(c *fiber.Ctx) error
	FetchQuiz(c *fiber.Ctx) error
	AnswerQuiz(c *fiber.Ctx) error
}

type HandlerImpl struct {
	Service Service
}

func NewHandler(service Service) Handler {
	return &HandlerImpl{Service: service}
}

func (h *HandlerImpl) ListCases(c *fiber.Ctx) error {
	users, pagination, appError := h.Service.List(controller_utils.GetCallerDetails(c), controller_utils.GetPaginationDetails(c))
	fmt.Println("pagination", pagination)
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	// return controller_utils.Respond(c, config.FetchSuccess, "Successfully fetched the user", users)
	return controller_utils.RespondWithPagination(c, config.FetchSuccess, "Successfully fetched the user", users, &pagination)
}

func (h *HandlerImpl) FetchQuiz(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return controller_utils.RespondError(c, config.InsertFailed, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	users, pagination, appError := h.Service.FetchQuiz(id)
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	return controller_utils.RespondWithPagination(c, config.FetchSuccess, "Successfully fetched the quiz", users, &pagination)
}

func (h *HandlerImpl) AnswerQuiz(c *fiber.Ctx) error {
	var body []RequestAnswerQuiz
	err := c.BodyParser(&body)
	if err != nil {
		return controller_utils.RespondError(c, config.InsertFailed, &error_utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	appError := h.Service.AnswerQuiz(&body)
	if appError != nil {
		return controller_utils.RespondError(c, config.InsertFailed, appError)
	}
	return controller_utils.RespondWithPagination(c, config.FetchSuccess, "Successfully answered the quiz", nil, nil)
}