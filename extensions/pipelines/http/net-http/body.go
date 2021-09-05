package net_http

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"

	Errors "github.com/vortex14/gotyphoon/errors"
)

func GetBody(client *http.Client, request *http.Request) (error, *string, *http.Response) {
	response, err := client.Do(request)
	if err != nil { color.Red("%s", err); return Errors.ResponseReadError, nil, nil }
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil { color.Red("%s", errC.Error()) }
	}(response.Body)

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return Errors.ResponseHttpGzipDecodeError, nil, nil
		}
		defer func(reader io.ReadCloser) {
			errR := reader.Close()
			if errR != nil {
				color.Red("%s", errR.Error())
			}
		}(reader)
	default:
		reader = response.Body
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err, nil, nil
	}
	textData := string(data)
	return nil, &textData, response
}
