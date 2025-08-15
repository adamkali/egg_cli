package templates

import "embed"

//go:embed *.tmpl cmd/*.tmpl controllers/*.tmpl middlewares/configs/*.tmpl services/*.tmpl cmd/configuration/*.tmpl database/migrations/*.tmpl database/queries/*tmpl models/handlers/*tmpl models/requests/*.tmpl models/responses/*.tmpl services/*.tmpl web/dist/index.html.tmpl web/dist/egg.png
var templates embed.FS

func GetTemplates() embed.FS {
	return templates
}
