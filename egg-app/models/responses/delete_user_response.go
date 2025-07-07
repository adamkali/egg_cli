
package responses

import (
	"github.com/adamkali/egg-app/db/repository"
)

type DashboardResponse struct {
	AuthenticatedUser           *repository.User
	PresignedUserProfilePicture *string
}

type DashboardDetailedResponse struct {
	Data    DashboardResponse `json:"data"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
}
