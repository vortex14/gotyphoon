package net_http

import (
	"net/http"
	"time"

	"github.com/vortex14/gotyphoon/interfaces"
)

func GetHttpClient(task interfaces.TaskInterface) *http.Client {
	_, client := GetHttpClientTransport(task)
	return client
}

func GetHttpClientTransport(task interfaces.TaskInterface) (*http.Transport, *http.Client) {
	transport := &http.Transport{
		ResponseHeaderTimeout: time.Duration(task.GetFetcherTimeout()) * time.Second,
		IdleConnTimeout: time.Duration(task.GetFetcherTimeout()) * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout: time.Duration(task.GetFetcherTimeout()) * time.Second,
	}

	return transport, client
}
