package interfaces

type FetcherSettings struct {
	Port   int `yaml:"port"`
	Queues struct {
		Priority          Queue `yaml:"priority"`
		ProcessorPriority Queue `yaml:"processor_priority"`
		Deferred          Queue `yaml:"deferred"`
		Retries           Queue `yaml:"retries"`
		Exceptions        Queue `yaml:"exceptions"`
	} `yaml:"queues"`
}

type ProcessorSettings struct {
	Port   int `yaml:"port"`
	Queues struct {
		Priority                  Queue `yaml:"priority"`
		SchedulerPriority         Queue `yaml:"scheduler_priority"`
		ResultTransporterPriority Queue `yaml:"result_transporter_priority"`
		FetcherRetries            Queue `yaml:"fetcher_retries"`
		Deferred                  Queue `yaml:"deferred"`
		Exceptions                Queue `yaml:"exceptions"`
	}
}

type TransporterSettings struct {
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
}

type SchedulerSettings struct {
	Port   int `yaml:"port"`
	Queues struct {
		Priority          Queue `yaml:"priority"`
		FetcherPriority   Queue `yaml:"fetcher_priority"`
		ProcessorPriority Queue `yaml:"processor_priority"`
		ProcessorDeferred Queue `yaml:"processor_deferred"`
		FetcherDeferred   Queue `yaml:"fetcher_deferred"`
		Exceptions        Queue `yaml:"exceptions"`
	} `yaml:"queues"`
}


type DonorSettings struct {
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
}


func (q *Queue) SetGroupName(name string)  {
	q.group = name
}

func (q *Queue) GetGroupName() string {
	return q.group
}

func (q *Queue) SetComponentName(name string)  {
	q.component = name
}

func (q *Queue) GetComponentName() string {
	return q.component
}

func (q *Queue) GetPriority() int {
	return q.priority
}

func (q *Queue) SetPriority(number int) {
	q.priority = number
}