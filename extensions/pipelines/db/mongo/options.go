package mongo

type Options struct {
	AuthSource string `yaml:"authSource,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
}

type Mongo struct {
	Name    string   `yaml:"name"`
	Details Options  `yaml:"details"`
	DbNames []string `yaml:"db_names"`
}
