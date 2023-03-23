package task

type ProjectSave map[string]interface{}

type ProcessorTask struct {
	ErrorResponse       bool `json:"error_response" default:"false" fake:"{randomstring:[false]}"`
	MaxProcessorRetries int  `json:"max_processor_retries" default:"5"  fake:"{number:4,5}"`
	Callback            struct {
		Name string `json:"name" default:"first_group" fake:"{randomstring:[first_group]}"`
		Type string `json:"type" default:"pipelines_group" fake:"{randomstring:[pipelines_group]}"`
	} `json:"callback"`
	Strategy string `json:"strategy" default:"text" fake:"{randomstring:[text]}"`
	Save     struct {
		Project ProjectSave `json:"project"`
		System  struct {
			ProcessorRetries int `json:"processor_retries" fake:"skip"`
		} `json:"system"`
	} `json:"save"`
	History []interface{} `json:"history" fake:"skip"`
}

func (t *ProcessorTask) IsMaxProcessorRetry() bool {
	return t.ErrorResponse && t.MaxProcessorRetries <= t.Save.System.ProcessorRetries
}

func (t *ProcessorTask) IsProcessorRetry() bool {
	return t.ErrorResponse && t.MaxProcessorRetries > t.Save.System.ProcessorRetries
}
