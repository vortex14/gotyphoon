package task


type SchedulerTask struct {
	Age    int    `json:"age" default:"0" fake:"{randomstring:[0]}"`
	SendTo string `json:"send_to" default:"fetcher" fake:"{randomstring:[fetcher]}"`
}