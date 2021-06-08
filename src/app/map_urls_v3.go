package app

import (
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/pages/quiz"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/middlewares/user"
	crud2 "github.com/dembygenesis/quiz_maker_auth/src/v3/api/pages/cases"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/pages/user_management/case_access"
	"github.com/dembygenesis/quiz_maker_auth/src/v3/api/pages/user_management/crud"
)

func mapUrlsV3(app *fiber.App) {
	api := app.Group("/api/v1", cors.New(), logger.New())

	userManagementCrud := crud.Initialize()
	caseAccessCrud := case_access.Initialize()
	casesV2 := crud2.Initialize()
	quiz := quiz.Initialize()
	
	// Quiz
	api.Get("/quiz/:id",
		// user.RoleMiddlewareV2([]string{"Admin"}),
		quiz.FetchQuiz)

	api.Post("/quiz/:id",
		// user.RoleMiddlewareV2([]string{"Admin"}),
		quiz.AnswerQuiz)

	// Crud
	api.Get("/user-management/crud", user.RoleMiddlewareV2([]string{"Admin"}), userManagementCrud.ListUsers)
	api.Post("/user-management/crud", user.RoleMiddlewareV2([]string{"Admin"}), userManagementCrud.CreateUser)

	api.Get("/user-management/crud/user-types",
		user.RoleMiddlewareV2([]string{"Admin", "Organization Member", "Law Firm"}),
		userManagementCrud.ListUserTypes)

	// Case access
	api.Post("/user-management/case-access",
		user.RoleMiddlewareV2([]string{"Admin"}),
		// cases.HasCaseAccess(),
		caseAccessCrud.Create)
	api.Delete("/user-management/case-access",
		user.RoleMiddlewareV2([]string{"Admin"}),
		caseAccessCrud.Delete)

	// Cases
	api.Get("/cases-v2", user.RoleMiddlewareV2([]string{"Admin"}), casesV2.ListCases)


}
