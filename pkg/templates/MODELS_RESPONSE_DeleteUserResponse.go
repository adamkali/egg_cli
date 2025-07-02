package templates

const MODELS_RESPONSE_DeleteUserResponseTemplate = `
package responses

import (
	"{{.Namespace}}/internal/repository"
)

type DashboardResponse struct {
	AuthenticatedUser           *repository.User
	PresignedUserProfilePicture *string
}

type DashboardDetailedResponse struct {
	Data    DashboardResponse ` + "`" +`json:"data"` + "`" +`
	Success bool              ` + "`" +`json:"success"` + "`" +`
	Message string            ` + "`" +`json:"message"` + "`" +`
}
`
