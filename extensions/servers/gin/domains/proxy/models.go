package proxy

import (
	"errors"
	"net/url"

	"github.com/vortex14/gotyphoon/interfaces"
)

type Proxy struct {
	Address  string
	Region   string
	Provider string
}

type Settings struct {
	BlockedTime      int
	CheckTime        int
	CheckBlockedTime int
	ConcurrentCheck  int
	Port             int
	PrefixNamespace  string
	CheckHosts       []string

	interfaces.RedisDetails
}

type Stats struct {
	ObservableHosts interface{}
	Stats           interface{}
	Active          interface{}
	Blocked         interface{}
	Locked          interface{}
	List            interface{}
	Allowed         interface{}
	TestEndpoint    *UpdateCheckPayload
}

type UpdateCheckPayload struct {
	Headers   map[string]string `json:"headers"`
	UrlSource string            `json:"url"`
	Url       *url.URL
}

func (u *UpdateCheckPayload) ParseUrl() (error, *url.URL) {

	var exception error

	if len(u.UrlSource) == 0 && u.Url != nil {
		u.UrlSource = u.Url.String()
	} else if len(u.UrlSource) > 0 {
		urlDecoded, e := url.Parse(u.UrlSource)
		exception = e
		if e == nil {
			u.Url = urlDecoded
		}
	} else {
		exception = errors.New("not found source url or *url.URL")
	}

	return exception, u.Url
}
