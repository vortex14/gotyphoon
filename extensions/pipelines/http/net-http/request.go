package net_http

import (
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"
)

func Request(client *http.Client, request *http.Request, logger interfaces.LoggerInterface) (error, *http.Response, *string) {

	err, body, response := GetBody(client, request)

	if err != nil { return Errors.ResponseReadError, nil, nil }

	if len(*body) == 0 { return Errors.ResponseEmptyError, nil, nil }

	return nil, response, body
}
