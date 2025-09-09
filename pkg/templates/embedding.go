/*
Copyright © 2025 Adam Kalinowski <adam.kalilarosa@proton.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package templates

import "embed"

//go:embed *.tmpl cmd/*.tmpl controllers/*.tmpl middlewares/configs/*.tmpl services/*.tmpl cmd/configuration/*.tmpl database/migrations/*.tmpl database/queries/*tmpl models/handlers/*tmpl models/requests/*.tmpl models/responses/*.tmpl services/*.tmpl web/dist/index.html.tmpl web/dist/egg.png 
var templates embed.FS

func GetTemplates() embed.FS {
	return templates
}
