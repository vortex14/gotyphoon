package proxy

import "github.com/vortex14/gotyphoon/interfaces"

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
	TestEndpoint    string
}
