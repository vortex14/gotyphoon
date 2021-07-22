package task

type ProcessorTask struct {
	ErrorResponse bool `json:"error_response" default:"false"`
	MaxProcessorRetries int `json:"max_processor_retries" default:"5"`
	Callback struct {
		Name string `json:"name" default:"first_group"`
		Type string `json:"type" default:"pipelines_group"`
	} `json:"callback"`
	Strategy string `json:"strategy" default:"text"`
	Save     struct {
		Project struct {
		} `json:"project"`
		System struct {
			ProcessorRetries int `json:"processor_retries"`
		} `json:"system"`
	} `json:"save"`
	History []interface{} `json:"history"`
}


func (t *ProcessorTask) IsMaxProcessorRetry() bool{
	return t.ErrorResponse && t.MaxProcessorRetries <= t.Save.System.ProcessorRetries
}



func (t *ProcessorTask) IsProcessorRetry() bool {
	return t.ErrorResponse && t.MaxProcessorRetries > t.Save.System.ProcessorRetries
}
