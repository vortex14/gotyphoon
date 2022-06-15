package errors

import Errors "errors"


var DockerCommandAlreadyStarted = Errors.New("docker command already started")

var DockerCommandNotStarted = Errors.New("docker command not started")

var DockerStoutAlreadySet = Errors.New("docker stdout already set")

var DockerStdErrAlreadySet = Errors.New("docker stderr already set")

var DockerStdInAlreadySet = Errors.New("Stdin already set")