package mongo

import "github.com/vortex14/gotyphoon/interfaces"

func CreateMongoServiceWithoutAuth(name string, host string, port int, dbNames []string) *Service  {
	return &Service{
		Settings: interfaces.ServiceMongo{
			DbNames: dbNames,
			Name: name,
			Details: struct {
				AuthSource string `yaml:"authSource,omitempty"`
				Username   string `yaml:"username,omitempty"`
				Password   string `yaml:"password,omitempty"`
				Host       string `yaml:"host"`
				Port       int    `yaml:"port"`
			}{
				Host: host,
				Port: port,
			},
		},
	}
}
