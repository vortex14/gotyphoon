package interfaces

import (
	"github.com/vortex14/gotyphoon/builders"
	"github.com/xanzy/go-gitlab"
)

type GoTemplate struct {
	Source string
	ExportPath string
	Data interface{}
}

type Discovery struct {
	Port    int    `yaml:"port,omitempty"`
	Host    string `yaml:"host,omitempty"`
	Cluster string `yaml:"cluster,omitempty"`
}


type ClusterLabel struct {
	Kind string
	Version string
}

type ClusterGitlab struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Variables []*gitlab.PipelineVariable `yaml:"variables,omitempty"`
}

type ClusterGrafana struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	FolderId string `yaml:"folder_id,omitempty"`
}

type ClusterDocker struct {
	Image string
}

type ClusterMeta struct {
	Gitlab  ClusterGitlab  `yaml:"gitlab,omitempty"`
	Grafana ClusterGrafana `yaml:"grafana,omitempty"`
	Docker  ClusterDocker  `yaml:"docker,omitempty"`
}



type DockerLabel struct {

}

type FileObject struct {
	Type string
	Path string
	Name string
	Data string
	FileSystem
}

type Yield struct {
	Name string
	Topic string
	Channel string
}

type GrafanaConfig struct {
	Id string `yaml:"id,omitempty"`
	Name string `yaml:"name,omitempty"`
	FolderId string `yaml:"folder_id,omitempty"`
	DashboardUrl string `yaml:"dashboard_url,omitempty"`
}






type Producers map[string] *Producer
type Consumers map[string] [] *Consumer
type Pipelines map[string] BasePipelineInterface

type ReplaceLabel struct {
	Label string
	Value string
}

type ReplaceLabels []*ReplaceLabel


type MapFileObjects map[string]*FileObject
type BuilderOptions builders.BuildOptions

type ServiceRedis struct {
	Name    string `yaml:"name"`
	Details struct {
		Host     string      `yaml:"host"`
		Port     int         `yaml:"port"`
		Password interface{} `yaml:"password"`
	} `yaml:"details"`
	Service
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
	Service
}


type Services struct {
	Mongo struct {
		Production []ServiceMongo `yaml:"production"`
		Debug []ServiceMongo      `yaml:"debug"`
	} `yaml:"mongo"`
	Redis struct {
		Production []ServiceRedis `yaml:"production"`
		Debug []ServiceRedis      `yaml:"debug"`
	} `yaml:"redis"`
}


type ConfigProject struct {
	ConfigInterface

	configFile string
	configPath string
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
	Grafana             []GrafanaConfig
	WaitingTasks        int       `yaml:"waiting_tasks"`
	PauseTime           int       `yaml:"pause_time"`
	MaxProcessorRetries int       `yaml:"max_processor_retries"`
	MaxFailed           int       `yaml:"max_failed"`
	MemoryLimit         float64   `yaml:"memory_limit"`
	RetryingDelay       int       `yaml:"retrying_delay"`
	RegisterService     Discovery `yaml:"register_service,omitempty"`
	TyComponents        struct {
		Fetcher FetcherSettings `yaml:"fetcher"`
		ResultTransporter TransporterSettings `yaml:"result_transporter"`
		Scheduler SchedulerSettings `yaml:"scheduler"`
		Processor ProcessorSettings `yaml:"processor"`
		Donor DonorSettings `yaml:"donor"`
	} `yaml:"ty_components"`
	Services Services `yaml:"services"`
}

type Queue struct {
	priority   int
	component  string
	group 	   string
	Concurrent int    `yaml:"concurrent"`
	MsgTimeout int    `yaml:"msg_timeout"`
	Channel    string `yaml:"channel"`
	Topic      string `yaml:"topic"`
	Share      bool   `yaml:"share"`
	Writable   bool   `yaml:"writable"`
	Readable   bool   `yaml:"readable"`
}