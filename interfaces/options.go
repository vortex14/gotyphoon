package interfaces

type MessageBrokerOptions struct {
	EnabledConsumer bool
	EnabledProducer bool
	Active bool
}

type BaseServiceOptions struct {
	Active bool
}

type TyphoonIntegrationsOptions struct {
	NSQ MessageBrokerOptions

	Mongo BaseServiceOptions
	Redis BaseServiceOptions
	DGraph BaseServiceOptions
	Badger BaseServiceOptions
	Circus BaseServiceOptions
	Yugabyte BaseServiceOptions
	ArangoDB BaseServiceOptions
}