package interfaces

import "fmt"

type TracingOptions struct {
	EnableInfoLog bool
	JaegerHost string
	JaegerPort int
	UseBanner bool
	UseUTC bool

}

func (o *TracingOptions) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", o.JaegerHost, o.JaegerPort)
}