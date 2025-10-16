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

package targets


var (
	// the tools that are required for the fullstack_app
	// air   -> hot reloading
	// swag  -> generating docs and frontend/api endpoint connections 
	// goose -> migrations
	// sqlc  -> generating queries
	RequiredTools = []string{
		"github.com/air-verse/air@latest",
		"github.com/swaggo/swag/cmd/swag@latest",
		"github.com/pressly/goose/v3/cmd/goose@latest",
		"github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
	}
)
