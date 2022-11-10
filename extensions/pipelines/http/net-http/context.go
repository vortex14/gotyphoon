package net_http

import (
	Context "context"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"

	"github.com/vortex14/gotyphoon/ctx"
)

const (
	TASK   = "task"
	CLIENT = "client"

	TRANSPORT = "transport"
	REQUEST   = "request"
	RESPONSE  = "response"
	DATA      = "data"
)

type ValidationCallback func(logger interfaces.LoggerInterface, response *http.Response, doc *goquery.Document) bool

func NewRequestCtx(context Context.Context, request *http.Request) Context.Context {
	return ctx.Update(context, REQUEST, request)
}

func GetRequestCtx(context Context.Context) (bool, *http.Request) {
	request, ok := ctx.Get(context, REQUEST).(*http.Request)
	return ok, request
}

func NewClientCtx(context Context.Context, client *http.Client) Context.Context {
	return ctx.Update(context, CLIENT, client)
}

func GetClientCtx(context Context.Context) (bool, *http.Client) {
	client, ok := ctx.Get(context, CLIENT).(*http.Client)
	return ok, client
}

func NewResponseCtx(context Context.Context, response *http.Response) Context.Context {
	return ctx.Update(context, RESPONSE, response)
}

func NewResponseDataCtx(context Context.Context, data *string) Context.Context {
	return ctx.Update(context, DATA, data)
}

func GetResponseCtx(context Context.Context) (bool, *http.Response, *string) {
	response, ok := ctx.Get(context, RESPONSE).(*http.Response)
	data, okD := ctx.Get(context, DATA).(*string)
	return ok && okD, response, data
}

func NewTransportCtx(context Context.Context, transport *http.Transport) Context.Context {
	return ctx.Update(context, TRANSPORT, transport)
}

func GetTransportCtx(context Context.Context) (bool, *http.Transport) {
	transport, ok := ctx.Get(context, TRANSPORT).(*http.Transport)
	return ok, transport
}
