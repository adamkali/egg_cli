package responses

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DeleteUserResponse struct {
    Data uuid.UUID `json:"data"`
	Success bool `json:"success"`
    Message string `json:"message"`
} // @name DeleteUserResponse

func NewDeleteUserResponse() *DeleteUserResponse {
    return &DeleteUserResponse{ Success: false, Message: "" }
}

func (DeleteUserResponse *DeleteUserResponse) Fail(ctx echo.Context, code int, err error) error {
    DeleteUserResponse.Message = err.Error()
    return ctx.JSON(code, DeleteUserResponse)
}

func (DeleteUserResponse *DeleteUserResponse) Successful(ctx echo.Context, delete_user_id uuid.UUID) error {
    DeleteUserResponse.Data = delete_user_id
    DeleteUserResponse.Success = true
    return ctx.JSON(200, DeleteUserResponse)
}



