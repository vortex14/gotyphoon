package interfaces

import (
	"github.com/vortex14/gotyphoon/extensions/project/python3/builders"
	"github.com/xanzy/go-gitlab"
)

type GoTemplate struct {
	Source     string
	ExportPath string
	Data       interface{}
}

type Discovery struct {
	Port    int    `yaml:"port,omitempty"`
	Host    string `yaml:"host,omitempty"`
	Cluster string `yaml:"cluster,omitempty"`
}

type ClusterLabel struct {
	Kind    string
	Version string
}

type ClusterGitlab struct {
	Endpoint  string                     `yaml:"endpoint,omitempty"`
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

func (f *FileObject) GetPath() string {
	if f.Path == "." {
		return f.Name
	}
	return f.Path + "/" + f.Name
}

type Yield struct {
	Name    string
	Topic   string
	Channel string
}

type GrafanaConfig struct {
	Id           string `yaml:"id,omitempty"`
	Name         string `yaml:"name,omitempty"`
	FolderId     string `yaml:"folder_id,omitempty"`
	DashboardUrl string `yaml:"dashboard_url,omitempty"`
}

type Producers map[string]*Producer
type Consumers map[string][]*Consumer
type Pipelines map[string]BasePipelineInterface

type ReplaceLabel struct {
	Label string
	Value string
}

type ReplaceLabels []*ReplaceLabel

type MapFileObjects map[string]*FileObject
type BuilderOptions builders.BuildOptions

type RedisDetails struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type ServiceRedis struct {
	Name    string       `yaml:"name"`
	Details RedisDetails `yaml:"details"`
	Service
}

type MongoDetails struct {
	AuthSource string `yaml:"authSource,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
}

type ServiceMongo struct {
	DefaultCollection string       `yaml:"default_collection"`
	DefaultDatabase   string       `yaml:"default_database"`
	Name              string       `yaml:"name"`
	Details           MongoDetails `yaml:"details"`
	DbNames           []string     `yaml:"db_names"`
	Service
}

func (s *ServiceMongo) RenameDefaultDatabase(name string) *ServiceMongo {
	return &ServiceMongo{
		DefaultDatabase:   name,
		DefaultCollection: s.DefaultCollection,
		Details:           s.Details,
		DbNames:           s.DbNames,
		Name:              s.Name,
	}
}

func (s *ServiceMongo) RenameDefaultCollection(name string) *ServiceMongo {
	return &ServiceMongo{
		DefaultDatabase:   s.DefaultDatabase,
		DefaultCollection: name,
		Details:           s.Details,
		DbNames:           s.DbNames,
		Name:              s.Name,
	}
}

func (s *ServiceMongo) RenameDefaultCollectionAndDB(dbname string, collection string) *ServiceMongo {
	return &ServiceMongo{
		DefaultDatabase:   dbname,
		DefaultCollection: collection,
		Details:           s.Details,
		DbNames:           s.DbNames,
		Name:              s.Name,
	}
}

type Services struct {
	Mongo struct {
		Production []ServiceMongo `yaml:"production"`
		Debug      []ServiceMongo `yaml:"debug"`
	} `yaml:"mongo"`
	Redis struct {
		Production []ServiceRedis `yaml:"production"`
		Debug      []ServiceRedis `yaml:"debug"`
	} `yaml:"redis"`
}

type ConfigProject struct {
	ConfigInterface

	configFile              string
	configPath              string
	ProjectName             string      `yaml:"project_name"`
	Debug                   bool        `yaml:"debug"`
	DefaultRetriesDelay     int         `yaml:"default_retries_delay"`
	PriorityDepthCheckDelay int         `yaml:"priority_depth_check_delay"`
	TaskTimeout             int         `yaml:"task_timeout"`
	Port                    int         `yaml:"port"`
	InstancesBucketLimit    int         `yaml:"instances_bucket_limit"`
	FinishedTasks           int         `yaml:"finished_tasks"`
	ProxyManagerAPI         string      `yaml:"proxy-manager-api"`
	MaxRetries              int         `yaml:"max_retries"`
	AutoThrottling          bool        `yaml:"auto_throttling"`
	IsRunning               bool        `yaml:"is_running"`
	NsqlookupdIP            string      `yaml:"nsqlookupd_ip"`
	Consolidator            interface{} `yaml:"consolidator,omitempty"`
	ApiKey                  string      `yaml:"API_KEY,omitempty"`
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
		Fetcher           FetcherSettings     `yaml:"fetcher"`
		ResultTransporter TransporterSettings `yaml:"result_transporter"`
		Scheduler         SchedulerSettings   `yaml:"scheduler"`
		Processor         ProcessorSettings   `yaml:"processor"`
		Donor             DonorSettings       `yaml:"donor"`
	} `yaml:"ty_components"`
	Services Services `yaml:"services"`

	Concurrent int    `yaml:"concurrent,omitempty"`
	Channel    string `yaml:"channel,omitempty"`
	Topic      string `yaml:"topic,omitempty"`
}

type Queue struct {
	priority   int    `yaml:"priority,omitempty"`
	component  string `yaml:"component,omitempty"`
	group      string `yaml:"group,omitempty"`
	Concurrent int    `yaml:"concurrent,omitempty"`
	MsgTimeout int    `yaml:"msg_timeout,omitempty"`
	Channel    string `yaml:"channel,omitempty"`
	Topic      string `yaml:"topic,omitempty"`
	Share      bool   `yaml:"share,omitempty"`
	Writable   bool   `yaml:"writable,omitempty"`
	Readable   bool   `yaml:"readable,omitempty"`
}
