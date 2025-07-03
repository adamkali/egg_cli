package templates

import (
	"text/template"

	"github.com/adamkali/egg_cli/pkg/configuration"
)

func Mapping(config *configuration.Configuration) map[string]*template.Template {
	return map[string]*template.Template{
		"main.go":                                                 template.Must(template.New("main.go").Parse(MainGoTemplate)),
		"openapitools.json":                                       template.Must(template.New("openapitools.yaml").Parse(OpenapitoolsJSONTemplate)),
		"sqlc.yaml":                                               template.Must(template.New("sqlc.yaml").Parse(SQLCYamlTemplate)),
		"README.md":                                               template.Must(template.New("README.md").Parse(READMETemplate)),
		"Makefile":                                                template.Must(template.New("Makefile").Parse(MakefileTemplate)),
		"Dockerfile":                                              template.Must(template.New("Dockerfile").Parse(DockerfileTemplate)),
		".gitignore":                                              template.Must(template.New(".gitignore").Parse(GitignoreTemplate)),
		".dockerignore":                                           template.Must(template.New(".dockerignore").Parse(DockerignoreTemplate)),
		".air.toml":                                               template.Must(template.New(".air.toml").Parse(AirTomlTemplate)),
		"./cmd/configuration/configuration.go":                    template.Must(template.New("./cmd/configuration/configuration.go").Parse(CmdConfigurationConfigurationTemplate)),
		"./cmd/db.go":                                             template.Must(template.New("./cmd/db.go").Parse(DBCmdTemplate)),
		"./cmd/migrate.go":                                        template.Must(template.New("./cmd/migrate.go").Parse(MigrateCmdTemplate)),
		"./cmd/root.go":                                           template.Must(template.New("./cmd/root.go").Parse(RootCmdTemplate)),
		"./cmd/swag.go":                                           template.Must(template.New("./cmd/swag.go").Parse(SwagCmdTemplate)),
		"./cmd/up.go":                                             template.Must(template.New("./cmd/up.go").Parse(UpCmdTemplate)),
		"./cmd/down.go":                                           template.Must(template.New("./cmd/down.go").Parse(DownCmdTemplate)),
		"./cmd/generate.go":                                       template.Must(template.New("./cmd/generate.go").Parse(GenerateCmdTemplate)),
		"./cmd/version.go":                                        template.Must(template.New("./cmd/version.go").Parse(VersionCmdTemplate)),
		"./cmd/bump.go":                                           template.Must(template.New("./cmd/bump.go").Parse(BumpCmdTemplate)),
		"./services/auth_service.go":                              template.Must(template.New("./services/auth_service.go").Parse(SERVICES_AuthServiceTemplate)),
		"./services/minio_service.go":                             template.Must(template.New("./services/minio_service.go").Parse(SERVICES_MinioServiceTemplate)),
		"./services/redis_service.go":                             template.Must(template.New("./services/redis_service.go").Parse(SERVICES_RedisServiceTemplate)),
		"./services/user_service.go":                              template.Must(template.New("./services/user_service.go").Parse(SERVICES_UserServiceTemplate)),
		"./services/validator_service.go":                         template.Must(template.New("./services/validator_service.go").Parse(SERVICES_ValidatorServiceTemplate)),
		"./services/mock_auth_service.go":                         template.Must(template.New("./services/mock_auth_service.go").Parse(SERVICES_MockAuthServiceTemplate)),
		"./services/mock_user_service.go":                         template.Must(template.New("./services/mock_user_service.go").Parse(SERVICES_MockUserServiceTemplate)),
		"./services/i_auth_service.go":                            template.Must(template.New("./services/i_auth_service.go").Parse(SERVICES_IAuthServiceTemplate)),
		"./services/i_minio_service.go":                           template.Must(template.New("./services/i_minio_service.go").Parse(SERVICES_IMinioServiceTemplate)),
		"./services/i_redis_service.go":                           template.Must(template.New("./services/i_redis_service.go").Parse(SERVICES_IRedisServiceTemplate)),
		"./services/i_user_service.go":                            template.Must(template.New("./services/i_user_service.go").Parse(SERVICES_IUserServiceTemplate)),
		"./controllers/controller.go":                             template.Must(template.New("./controllers/controller.go").Parse(CONTROLLER_ControllerTemplate)),
		"./controllers/routes.go":                                 template.Must(template.New("./controllers/routes.go").Parse(CONTROLLER_RoutesTemplate)),
		"./controllers/user_controller.go":                        template.Must(template.New("./controllers/user_controller.go").Parse(CONTROLLERS_UserControllerTemplate)),
		"./middlewares/configs/auth.go":                           template.Must(template.New("./middlewares/configs/auth.go").Parse(MIDDLEWARES_CONFIGS_AuthConfigTemplate)),
		"./middlewares/configs/static.go":                         template.Must(template.New("./middlewares/configs/static.go").Parse(MIDDLEWARES_CONFIGS_StaticConfigTemplate)),
		"./models/requests/login_request.go":                      template.Must(template.New("./models/requests/login_request.go").Parse(MODELS_REQUESTS_LoginRequestTemplate)),
		"./models/requests/new_user_request.go":                   template.Must(template.New("./models/requests/new_user_request.go").Parse(MODELS_REQUESTS_NewUserRequestTemplate)),
		"./models/responses/delete_user_response.go":              template.Must(template.New("./models/responses/delete_user_response.go").Parse(MODELS_RESPONSE_DeleteUserResponseTemplate)),
		"./models/responses/login_response.go":                    template.Must(template.New("./models/responses/login_response.go").Parse(MODELS_RESPONSE_LoginResponseTemplate)),
		"./models/responses/user_response.go":                     template.Must(template.New("./models/responses/user_response.go").Parse(MODELS_RESPONSE_UserResponseTemplate)),
		"./models/responses/users_response.go":                    template.Must(template.New("./models/responses/users_response.go").Parse(MODELS_RESPONSE_UsersResponseTemplate)),
		"./models/handlers/login_handler.go":                      template.Must(template.New("./models/handlers/login_handler.go").Parse(MODELS_HANDLERS_LoginHandlerTemplate)),
		"./models/handlers/register_handler.go":                   template.Must(template.New("./models/handlers/register_handler.go").Parse(MODELS_HANDLERS_RegisterHandlerTemplate)),
		"./models/handlers/delete_user_handler.go":                template.Must(template.New("./models/handlers/delete_user_handler.go").Parse(MODELS_HANDLERS_DeleteUserHandlerTemplate)),
		"./models/handlers/get_current_logged_in_user_handler.go": template.Must(template.New("./models/handlers/get_current_logged_in_user_handler.go").Parse(MODELS_HANDLERS_GetCurrentLoggedInUserHandlerTemplate)),
		"./models/handlers/get_profile_picture_handler.go":        template.Must(template.New("./models/handlers/get_profile_picture_handler.go").Parse(MODELS_HANDLERS_GetProfilePictureHandlerTemplate)),
		"./models/handlers/get_users_handler.go":                  template.Must(template.New("./models/handlers/get_users_handler.go").Parse(MODELS_HANDLERS_GetUsersHandlerTemplate)),
		"./models/handlers/upload_profile_picture_handler.go":     template.Must(template.New("./models/handlers/upload_profile_picture_handler.go").Parse(MODELS_HANDLERS_UploadProfilePictureHandlerTemplate)),
		"./" + config.Database.Migration.Destination + "/0001_init.sql": template.Must(
			template.New("./" + config.Database.Migration.Destination + "/migrations/0001_init.sql").Parse(DATABASE_MIGRATIONS_INITTemplate),
		),
		"./" + config.Database.QueriesLocation + "/token.go": template.Must(
			template.New("./" + config.Database.QueriesLocation + "/token.go").Parse(DATABASE_QUERIES_TokenTemplate),
		),

		"./" + config.Database.QueriesLocation + "/user.go": template.Must(
			template.New("./" + config.Database.QueriesLocation + "/user.go").Parse(DATABASE_QUERIES_UserTemplate),
		),
	}
}
