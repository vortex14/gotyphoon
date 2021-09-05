package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"net/http"
	"strings"

	"github.com/vortex14/gotyphoon/interfaces"
)

type HTTPResponseDefaultMiddleware struct {
	*label.MetaInfo
}

func (m *HTTPResponseDefaultMiddleware) Run(
	task *task.TyphoonTask,
	response *http.Response,
	) error {

	responseHeaders := make(map[string]string)
	for key, value := range response.Header {
		if key == "" && len(value) > 0 {
			responseHeaders[key] = strings.Join(value, "")
		}
	}

	task.Fetcher.Response.Headers = responseHeaders
	task.Fetcher.Response.Code = response.StatusCode

	return nil
}

func (m *HTTPResponseDefaultMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	catch func(err error),
	next func(ctx context.Context),

	) {
	taskInstance, _ := context.Value(TASK).(*task.TyphoonTask)
	response, _ := context.Value(RESPONSE).(*http.Response)
	if err := m.Run(taskInstance, response); err != nil {
		catch(err)
	}
}


func ConstructorHTTPResponseDefaultMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HTTPResponseDefaultMiddleware{
		MetaInfo: &label.MetaInfo{
			Required:    required,
			Name:        NAMEHttpBasicAuthMiddleware,
			Description: DESCRIPTIONHttpBasicAuthMiddleware,
		},
	}
}
