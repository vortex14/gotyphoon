package config


type Queue struct {
	Concurrent int    `yaml:"concurrent"`
	MsgTimeout int    `yaml:"msg_timeout"`
	Channel    string `yaml:"channel"`
	Topic      string `yaml:"topic"`
	Share      bool   `yaml:"share"`
	Writable   bool   `yaml:"writable"`
	Readable   bool   `yaml:"readable"`
}

type Project struct {
	Config Config
	ConfigFile string
}

type GrafanaConfig struct {
	Id string `yaml:"id,omitempty"`
	Name string `yaml:"name,omitempty"`
	FolderId string `yaml:"folder_id,omitempty"`
	DashboardUrl string `yaml:"dashboard_url,omitempty"`
}

type Discovery struct {
	Port    int    `yaml:"port,omitempty"`
	Host    string `yaml:"host,omitempty"`
	Cluster string `yaml:"cluster,omitempty"`
}

type ServiceRedis struct {
	Name    string `yaml:"name"`
	Details struct {
		Host     string      `yaml:"host"`
		Port     int         `yaml:"port"`
		Password interface{} `yaml:"password"`
	} `yaml:"details"`
}

type ServiceMongo struct {
	Name    string `yaml:"name"`
	Details struct {
		AuthSource string `yaml:"authSource,omitempty"`
		Username string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"details"`
	DbNames []string `yaml:"db_names"`
}

type Config struct {
	ProjectName             string `yaml:"project_name"`
	Debug                   bool   `yaml:"debug"`
	DefaultRetriesDelay     int    `yaml:"default_retries_delay"`
	PriorityDepthCheckDelay int    `yaml:"priority_depth_check_delay"`
	TaskTimeout             int    `yaml:"task_timeout"`
	Port                    int    `yaml:"port"`
	InstancesBucketLimit    int    `yaml:"instances_bucket_limit"`
	FinishedTasks           int    `yaml:"finished_tasks"`
	ProxyManagerAPI         string `yaml:"proxy-manager-api"`
	MaxRetries              int    `yaml:"max_retries"`
	AutoThrottling          bool   `yaml:"auto_throttling"`
	IsRunning               bool   `yaml:"is_running"`
	NsqlookupdIP            string `yaml:"nsqlookupd_ip"`
	Consolidator			interface{} `yaml:"consolidator,omitempty"`
	ApiKey					string `yaml:"API_KEY,omitempty"`
	NsqdNodes               []struct {
		IP string `yaml:"ip"`
	} `yaml:"nsqd_nodes"`
	Grafana [] GrafanaConfig
	WaitingTasks        int     `yaml:"waiting_tasks"`
	PauseTime           int     `yaml:"pause_time"`
	MaxProcessorRetries int     `yaml:"max_processor_retries"`
	MaxFailed           int     `yaml:"max_failed"`
	MemoryLimit         float64 `yaml:"memory_limit"`
	RetryingDelay       int     `yaml:"retrying_delay"`
	RegisterService  	Discovery `yaml:"register_service,omitempty"`
	TyComponents struct {
		Fetcher struct {
			Port   int `yaml:"port"`
			Queues struct {
				Priority          Queue `yaml:"priority"`
				ProcessorPriority Queue `yaml:"processor_priority"`
				Deferred          Queue `yaml:"deferred"`
				Retries           Queue `yaml:"retries"`
				Exceptions        Queue `yaml:"exceptions"`
			} `yaml:"queues"`
		} `yaml:"fetcher"`
		ResultTransporter struct {
			Port   int `yaml:"port"`
			Queues struct {
				Priority            Queue `yaml:"priority"`
				SchedulerPriority   Queue `yaml:"scheduler_priority"`
				FetcherPriority     Queue `yaml:"fetcher_priority"`
				ProcessorPriority   Queue `yaml:"processor_priority"`
				Exceptions          Queue `yaml:"exceptions"`
				FetcherExceptions   Queue `yaml:"fetcher_exceptions"`
				ProcessorExceptions Queue `yaml:"processor_exceptions"`
				SchedulerExceptions Queue `yaml:"scheduler_exceptions"`
			} `yaml:"queues"`
		} `yaml:"result_transporter"`
		Scheduler struct {
			Port   int `yaml:"port"`
			Queues struct {
				Priority          Queue `yaml:"priority"`
				FetcherPriority   Queue `yaml:"fetcher_priority"`
				ProcessorPriority Queue `yaml:"processor_priority"`
				ProcessorDeferred Queue `yaml:"processor_deferred"`
				FetcherDeferred   Queue `yaml:"fetcher_deferred"`
				Exceptions        Queue `yaml:"exceptions"`
			} `yaml:"queues"`
		} `yaml:"scheduler"`
		Processor struct {
			Port   int `yaml:"port"`
			Queues struct {
				Priority                  Queue `yaml:"priority"`
				SchedulerPriority         Queue `yaml:"scheduler_priority"`
				ResultTransporterPriority Queue `yaml:"result_transporter_priority"`
				FetcherRetries            Queue `yaml:"fetcher_retries"`
				Deferred                  Queue `yaml:"deferred"`
				Exceptions                Queue `yaml:"exceptions"`
			} `yaml:"queues"`
		} `yaml:"processor"`
		Donor struct {
			Port   int `yaml:"port"`
			Queues struct {
				Priority                  Queue `yaml:"priority"`
				FetcherDeferred           Queue `yaml:"fetcher_deferred"`
				FetcherPriority           Queue `yaml:"fetcher_priority"`
				ProcessorDeferred         Queue `yaml:"processor_deferred"`
				ProcessorPriority         Queue `yaml:"processor_priority"`
				ResultTransporterPriority Queue `yaml:"result_transporter_priority"`
				SchedulerPriority         Queue `yaml:"scheduler_priority"`
			} `yaml:"queues"`
		} `yaml:"donor"`
	} `yaml:"ty_components"`
	Services struct {
		Mongo struct {
			Production [] ServiceMongo  `yaml:"production"`
			Debug [] ServiceMongo `yaml:"debug"`
		} `yaml:"mongo"`
		Redis struct {
			Production []ServiceRedis `yaml:"production"`
			Debug []ServiceRedis`yaml:"debug"`
		} `yaml:"redis"`
	} `yaml:"services"`
}

func (c *Config) GetComponentPort(name string) int {
	var port int
	switch name {
	case "donor":
		component := c.TyComponents.Donor
		port = component.Port
	case "fetcher":
		component := c.TyComponents.Fetcher
		port = component.Port
	case "processor":
		component := c.TyComponents.Processor
		port = component.Port
	case "result_transporter":
		component := c.TyComponents.ResultTransporter
		port = component.Port
	case "scheduler":
		component := c.TyComponents.Scheduler
		port = component.Port
	}

	return port
}