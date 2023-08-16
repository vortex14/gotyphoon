package swagger

import "github.com/vortex14/gotyphoon/interfaces"

type OpenApi struct {
	Title       string `yaml:"title" json:"title" mapstructure:"title"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
	Contacts    struct {
		Name  string `yaml:"name" json:"name" mapstructure:"name"`
		URL   string `yaml:"url" json:"url" mapstructure:"url"`
		Email string `yaml:"email" json:"email" mapstructure:"email"`
	} `yaml:"contacts" json:"contacts" mapstructure:"contacts"`
	Host string `yaml:"host" json:"host" mapstructure:"host"`
	Json struct {
		Local string `yaml:"local" json:"local" mapstructure:"local"`
		Web   string `yaml:"web" json:"web" mapstructure:"web"`
	} `yaml:"json" json:"json" mapstructure:"json"`
	Security map[string]struct {
		Type string `yaml:"type" json:"type" mapstructure:"type"`
		Name string `yaml:"name" json:"name" mapstructure:"name"`
		In   string `yaml:"in" json:"in" mapstructure:"in"`
	}
	CodeSamples map[string]struct {
		Label  string `yaml:"label" json:"label" mapstructure:"label"`
		Source string `yaml:"source" json:"source" mapstructure:"source"`
	} `yaml:"codeSamples,omitempty" json:"codeSamples,omitempty" mapstructure:"codeSamples,omitempty"`
}

type Property struct {
	Name        string        `yaml:"name" json:"name" mapstructure:"name"`
	Description string        `yaml:"description" json:"description" mapstructure:"description"`
	Type        string        `yaml:"type" json:"type" mapstructure:"type"`
	Required    string        `yaml:"required" json:"required" mapstructure:"required"`
	Default     interface{}   `yaml:"default" json:"default" mapstructure:"default"`
	Example     interface{}   `yaml:"example" json:"example" mapstructure:"example"`
	Enum        []interface{} `yaml:"enum" json:"enum" mapstructure:"enum"`
	Deprecated  bool          `yaml:"deprecated" json:"deprecated" mapstructure:"deprecated"`
}

type Model struct {
	Name        string     `yaml:"name" json:"name" mapstructure:"name"`
	Description string     `yaml:"description" json:"description" mapstructure:"description"`
	Properties  []Property `yaml:"properties" json:"properties" mapstructure:"properties"`
}

type Server struct {
	Host string `yaml:"host" json:"host" mapstructure:"host"`
	Port int    `yaml:"port" json:"port" mapstructure:"port"`

	Scheme       string                                  `yaml:"scheme" json:"scheme" mapstructure:"scheme"`
	DefaultRoute string                                  `yaml:"defaultRoute" json:"defaultRoute" mapstructure:"defaultRoute"`
	Routes       map[string]interfaces.ResourceInterface `yaml:"routes" json:"routes" mapstructure:"routes"`
	Timeout      string                                  `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
	Domains      []string                                `yaml:"domains" json:"domains" mapstructure:"domains"`
	OpenApi      OpenApi                                 `yaml:"openapi" json:"openapi,omitempty" mapstructure:"openapi,omitempty"`
	JwtSecret    string                                  `yaml:"jwt-secret" json:"jwt-secret" mapstructure:"jwt-secret"`
	ApiVersion   string                                  `yaml:"apiVersion" json:"apiVersion" mapstructure:"apiVersion"`
	Models       map[string]Model                        `yaml:"models" json:"models" mapstructure:"models"`
	GeoIP        struct {
		DB string `yaml:"db,omitempty" json:"db,omitempty" mapstructure:"db,omitempty"`
	} `yaml:"geoip,omitempty" json:"geoip,omitempty" mapstructure:"geoip,omitempty"`
	ExcludeRequestKeys []string `yaml:"excludeRequestKeys" json:"excludeRequestKeys" mapstructure:"excludeRequestKeys"`
}

type Config struct {
	Parameters Parameters `yaml:"parameters" json:"parameters" mapstructure:"parameters"`
	Server     Server     `yaml:"server" json:"server" mapstructure:"server"`
}

type Parameters struct {
	Authentication struct {
		AuthHeader         string `yaml:"authHeader" json:"authHeader" mapstructure:"authHeader"`
		AuthCookie         string `yaml:"authCookie" json:"authCookie" mapstructure:"authCookie"`
		JwtUserIdKey       string `yaml:"jwtUserIdKey" json:"jwtUserIdKey" mapstructure:"jwtUserIdKey"`
		JwtSessionTokenKey string `yaml:"jwtSessionTokenKey" json:"jwtSessionTokenKey" mapstructure:"jwtSessionTokenKey"`
		UserIdKey          string `yaml:"userIdKey" json:"userIdKey" mapstructure:"userIdKey"`
		SessionTokenKey    string `yaml:"sessionTokenKey" json:"sessionTokenKey" mapstructure:"sessionTokenKey"`
	} `yaml:"authentication" json:"authentication" mapstructure:"authentication"`
}
