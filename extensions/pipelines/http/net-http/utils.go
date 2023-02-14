package net_http

import (
	"fmt"
	"net/http"

	"github.com/vortex14/gotyphoon/extensions/data/fake"
)

func FormattingProxy(proxy string) string {
	return fmt.Sprintf("http://%s", proxy)
}

func MakeBasicRequest(url string) (error, *http.Response, *string) {
	task := fake.CreateDefaultTask()
	task.SetFetcherUrl(url)
	request, err := NewRequest(task)
	if err != nil {
		return err, nil, nil
	}
	client := GetHttpClient(task)

	err, data, response := GetBody(client, request)
	if err != nil {
		return err, nil, nil
	}

	return nil, response, data
}
