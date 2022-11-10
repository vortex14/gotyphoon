package proxy

type Settings struct {
	BlockedTime      int
	CheckTime        int
	CheckBlockedTime int
	RedisHost        string
	ConcurrentCheck  int
	Port             int
	PrefixNamespace  string
	CheckHosts       []string
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
