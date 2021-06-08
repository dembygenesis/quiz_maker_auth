package user

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	UtilityString "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	ModelUser "github.com/dembygenesis/quiz_maker_auth/src/v1/api/models/user"
	ModelUserType "github.com/dembygenesis/quiz_maker_auth/src/v1/api/models/user_type"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateMiddleware(c *fiber.Ctx) error {

	var paramsInsert ModelUser.ParamsInsert

	// Make sure all parameters are present
	err := c.BodyParser(&paramsInsert)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "something went wrong when trying to parse the update parameters",
		}
		r.AddErrors("something went wrong when trying to parse the update parameters: " + err.Error())

		return c.JSON(r)

	}

	// Make sure there are no empty params
	emptyFields := paramsInsert.NoEmptyFields()

	if len(emptyFields) != 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}

		for _, val := range emptyFields {
			r.AddErrors(val)
		}

		return c.JSON(r)
	}

	// User Type Id Must Be Valid
	userType := ModelUserType.UserType{ID: paramsInsert.UserTypeId}

	res, err := userType.ValidID()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to check the user_type_id")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("user_type_id id must be valid")

		return c.JSON(r)
	}

	// Birthday must be a valid date format
	_, err = time.Parse("2006-01-02", paramsInsert.Birthday)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to parse the birthday: " + err.Error())

		return c.JSON(r)
	}

	// Gender must be M or F
	validGender := []string{"M", "F"}

	stringInSlice := UtilityString.StringInSlice(paramsInsert.Gender, validGender)

	if stringInSlice == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("Gender must be M or F")

		return c.JSON(r)
	}

	// Make sure email is not taken
	user := ModelUser.User{Email: paramsInsert.Email}

	res, err = user.EmailNotTaken()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to validate the email: " + err.Error())

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("email already taken")

		return c.JSON(r)
	}

	return c.Next()
}

func RoleMiddlewareV2(roles []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var r response_builder.Response

		// Attempt to validate token
		token := c.Get("authorization")

		user := ModelUser.User{Token: token}
		userId, userType, err := user.ValidateTokenV2(roles)

		if err != nil {
			r.HttpCode = 401
			r.ResponseMessage = "Unauthorized"
			r.AddErrors(err.Error())

			return c.Status(401).JSON(r)
		}

		// Also fetch organizations
		/*u := ModelUser.User{ID: userId}

		res, err := u.GetDetailsById()

		if err != nil {
			r.HttpCode = 401
			r.ResponseMessage = "Something went wrong when trying to validate the user"
			r.AddErrors(err.Error())

			return c.Status(401).JSON(r)
		}*/

		c.Locals("tokenExtractedUserId", userId)
		c.Locals("tokenExtractedUserType", userType)

		return c.Next()
	}
}

func LoginValidation(c *fiber.Ctx) error {

	var paramsLogin ParamsLogin
	var response response_builder.Response

	err := c.BodyParser(&paramsLogin)

	// Parsing must be fine
	if err != nil {
		response.HttpCode = 200
		response.ResponseMessage = "Error"
		response.AddErrors("Something went wrong with parsing the arguments for Login")

		return c.JSON(response)
	}

	// No empty inputs
	if paramsLogin.Email == "" || paramsLogin.Password == "" {
		response.HttpCode = 200
		response.ResponseMessage = "Error"
		response.AddErrors("email param is required")
		response.AddErrors("password param is required")

		return c.JSON(response)
	}

	// Email must exist in record
	user := ModelUser.UserLogin{Email: paramsLogin.Email}

	exists, err := user.ValidEmail()

	if err != nil {
		response.HttpCode = 200
		response.ResponseMessage = "Syntax Error"
		response.AddErrors(err.Error())

		return c.JSON(response)
	}

	if exists == false {
		response.HttpCode = 200
		response.ResponseMessage = "Syntax Error"
		response.AddErrors("Email does not exist")

		return c.JSON(response)
	}

	return c.Next()
}