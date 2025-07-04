package templates_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/templates"
)

const ResultMODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate = `
package handlers

import (
	"fmt"

	"github.com/adamkali/egg/internal/repository"
	"github.com/adamkali/egg/models/responses"
	"github.com/adamkali/egg/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetCurrentLoggedInUserHandler struct {
	UserID          uuid.UUID
	LoggedInUser    *repository.User
	Context         echo.Context
	Error           error
	Code            int
	Locked          bool
}

func NewGetCurrentLoggedInUserHandler(ctx echo.Context) *GetCurrentLoggedInUserHandler {
	return &GetCurrentLoggedInUserHandler{
		Context: ctx,
		Locked:  false,
		Error:   nil,
		Code:    200,
	}
}
func (h *GetCurrentLoggedInUserHandler) Lock(code int) *GetCurrentLoggedInUserHandler {
	h.Locked = true
	h.Code = code
	return h
}

func (h *GetCurrentLoggedInUserHandler) Handle(fun any) *GetCurrentLoggedInUserHandler {
	var code int
	if !h.Locked {
		switch handle := fun.(type) {
		case func(token string) error:
			jwt_token := h.Context.Get("user").(*jwt.Token)
			claims := jwt_token.Claims.(*services.CustomJwt)
			h.UserID = claims.UserId
			h.Error = handle(jwt_token.Raw)
			code = 401
		case func(user_id uuid.UUID) (*repository.User, error):
			h.LoggedInUser, h.Error = handle(h.UserID)
			code = 404
		default:
			code = 600
			h.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("Type assertion failed for type: %T\n", fun),
			)
		}
		if h.Error != nil {
			return h.Lock(code)
		}
	}
	return h
}

func (h *GetCurrentLoggedInUserHandler) JSON() error {
	var code int
	var message string
	if h.Locked && h.Error != nil {
		code = h.Code
		if code == 600 {
			message = "Misaligend handler on the server"
		} else {
			message = h.Error.Error()
		}
	} else {
		message = "OK"
		code = 200
	}
	return h.Context.JSON(code, responses.UserResponse{
		Data:    responses.UserDataFromRepository(h.LoggedInUser),
		Success: !h.Locked,
		Message: message,
	})

}
`

func TestMODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate(t *testing.T) {
	// load the template
	temp := templates.MODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate
	templateTest := template.Must(template.New("get_current_logged_in_user_handler.go").Parse(temp))

	// execute the template
	stringWriter := new(bytes.Buffer)
	err := templateTest.ExecuteTemplate(stringWriter, "get_current_logged_in_user_handler.go", createConfiguration())
	if err != nil {
		t.Error(err)
	}

	// check the result
	if stringWriter.String() != ResultMODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate {
		diff := Diff(stringWriter.String(), ResultMODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate)
		for i, v := range diff {
			t.Errorf("line %d: expected %s, got %s", i, v.Expected, v.Actual)
		}
	}
}
