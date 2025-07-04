package templates_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/templates"
)

const ResultMODELS_RESPONSE_DeleteUserResponseTemplate = `
package responses

import (
	"github.com/adamkali/egg/internal/repository"
)

type DashboardResponse struct {
	AuthenticatedUser           *repository.User
	PresignedUserProfilePicture *string
}

type DashboardDetailedResponse struct {
	Data    DashboardResponse ` + "`" + `json:"data"` + "`" + `
	Success bool              ` + "`" + `json:"success"` + "`" + `
	Message string            ` + "`" + `json:"message"` + "`" + `
}
`

func TestMODELS_RESPONSE_DeleteUserResponseTemplate(t *testing.T) {
	// load the template
	temp := templates.MODELS_RESPONSE_DeleteUserResponseTemplate
	templateTest := template.Must(template.New("delete_user_response.go").Parse(temp))

	// execute the template
	stringWriter := new(bytes.Buffer)
	err := templateTest.ExecuteTemplate(stringWriter, "delete_user_response.go", createConfiguration())
	if err != nil {
		t.Error(err)
	}

	// check the result
	if stringWriter.String() != ResultMODELS_RESPONSE_DeleteUserResponseTemplate {
		diff := Diff(stringWriter.String(), ResultMODELS_RESPONSE_DeleteUserResponseTemplate)
		for i, v := range diff {
			t.Errorf("line %d: expected %s, got %s", i, v.Expected, v.Actual)
		}
	}
}
