package crud

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/controller_utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	ListCases(c *fiber.Ctx) error
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