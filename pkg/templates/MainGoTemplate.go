package templates
const MainGoTemplate = `
/*
Copyright © {{.Copyright.Year}} {{.Copyright.Author}}  

This is made by the Full Stack Template
*/
package main

import (
	"{{.Namespace}}/cmd"
)

// @Title {{.Name}} 
// @Version {{.Semver}}
// @Description This is the swagger page for the Project {{.Name}} generated with Egg-go. use this to test your database connection
// @Contact.name Adam Kalinowski 
// @Contact.url https://{{.Namespace}}
// @License.name {{.License}}
// @BasePath /api
func main() {
	cmd.Execute()
}
`

