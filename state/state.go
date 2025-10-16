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

package state

type (
	ErrMsg error
)

const (
	ProjectHostName       = "PROJECTHOSTNAME"
	ProjectUsernameName   = "PROJECTUSERNAMENAME"
	ProjectNameName       = "PROJECTNAMENAME"
	DatabaseURLName       = "DATABASEURLNAME"
	DatabaseRootName      = "DATABASEROOTNAME"
	DatabaseSqlcOrGoName  = "DATABASESQLCORGONAME"
	LicenseName           = "LICENSENAME"
	CopyrightYearName     = "COPYRIGHTYEARNAME"
	CopyrightAuthorName   = "COPYRIGHTAUTHORNAME"
	ServerPortName        = "SERVERPORTNAME"
	ServerJWTName         = "SERVERJWTSECRETNAME"
	ServerFrontendDirName = "SERVERFRONTENDDIRNAME"
	ServerFrontendApiName = "SERVERFRONTENDAPINAME"
	MinioURLName          = "MINIOURLNAME"
	MinioAccessKeyName    = "MINIOACCESSKEYNAME"
	MinioSecretKeyName    = "MINIOSECRETKEYNAME"
)

var (
	ProjectNamespace   = ""
	ProjectHost        = ""
	ProjectUsername    = ""
	ProjectName        = ""
	ServerPort         = ""
	ServerJWT          = ""
	ServerFrontendDir  = ""
	ServerFrontendApi  = ""
	DatabaseURL        = ""
	DatabaseRoot       = ""
	DatabaseSqlcOrGo   = ""
	CacheURL           = ""
	MinioURL           = ""
	MinioAccessKey     = ""
	MinioSecretKey     = ""
	License            = ""
	CopyrightYear      = ""
	CopyrightAuthor    = ""
	ProjectSettingsMap = map[int]string{
		0: ProjectHostName,
		1: ProjectUsernameName,
		2: ProjectNameName,
	}
	ProjectSeverMap = map[int]string{
		0: ServerJWTName,
		1: ServerPortName,
	}
	DatabaseMap = map[int]string{
		0: DatabaseURLName,
		1: DatabaseSqlcOrGoName,
		2: DatabaseRootName,
	}
	LicenseMap = map[int]string{
		0: LicenseName,
		1: CopyrightYearName,
		2: CopyrightAuthorName,
	}
	MinioMap = map[int]string{
		0: MinioURLName,
		1: MinioAccessKeyName,
		2: MinioSecretKeyName,
	}
)
