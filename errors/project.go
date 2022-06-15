package errors

import Errors "errors"

var ProjectGitlabNotFound = Errors.New("project not found in Gitlab")

var ProjectNotFound = Errors.New("project does not exists in the current directory")

var ProjectInvalidEnv = Errors.New("we need set valid environment variables like TYPHOON_PATH and TYPHOON_PROJECTS")