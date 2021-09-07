package task

import (
	"context"
	"github.com/vortex14/gotyphoon/ctx"

	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/interfaces"
)

type TyphoonTask struct {
	Fetcher           FetcherTask     `json:"fetcher"`
	Processor         ProcessorTask   `json:"processor"`
	Scheduler         SchedulerTask   `json:"scheduler"`
	ResultTransporter TransporterTask `json:"result_transporter"`
	Priority          int             `json:"priority" default:"3" fake:"{randomstring:[3]}"`
	URL               string        `json:"url" default:"https://httpstat.us/200" fake:"{url}"`
	Taskid            string        `json:"taskid" default:"task-id" fake:"{uuid}"`
	ProjectName string `json:"project_name" fake:"skip"`
	//msg *nsq.Message
}

func (t *TyphoonTask) GetFetcherMethod() string {
	return t.Fetcher.Method
}

func (t *TyphoonTask) GetFetcherTimeout() int {
	return t.Fetcher.Timeout
}

func (t *TyphoonTask) GetFetcherUrl() string {
	return t.URL
}

func (t *TyphoonTask) SetFetcherUrl(url string)  {
	t.URL = url
}

func (t *TyphoonTask) SetProxyServerUrl(url string)  {
	t.Fetcher.IsProxyRequired = true
	t.Fetcher.ProxyServer = url
}

func (t *TyphoonTask) SetUserAgent(agent string)  {
	if t.Fetcher.Headers == nil {
		t.Fetcher.Headers = make(map[string]string)
	}
	t.Fetcher.Headers["User-Agent"] = agent
}

func (t *TyphoonTask) SetProxyAddress(address string)  {
	t.Fetcher.Proxy = address
}

func (t *TyphoonTask) IsProxyRequired() bool {
	return t.Fetcher.IsProxyRequired
}

func (t *TyphoonTask) GetProxyAddress() string {
	return t.Fetcher.Proxy
}

func (t *TyphoonTask) GetUserAgent() string {
	return t.Fetcher.Headers["User-Agent"]
}

func (t *TyphoonTask) GetProxyServerUrl() string  {
	return t.Fetcher.ProxyServer
}

func (t *TyphoonTask) SetStatusCode(code int)  {
	t.Fetcher.Response.Code = code
}

func Get(c context.Context) (bool, *TyphoonTask) {
	taskInstance, ok := ctx.Get(c, interfaces.TASK).(*TyphoonTask)
	return ok, taskInstance
}

func GetTaskCtx(context context.Context) *TyphoonTask {
	return ctx.Get(context, interfaces.TASK).(*TyphoonTask)
}

func NewTaskCtx(task *TyphoonTask) context.Context {
	return ctx.Update(context.Background(), interfaces.TASK, task)
}

func PatchCtx(context context.Context, task *TyphoonTask) context.Context {
	return ctx.Update(context, interfaces.TASK, task)
}

func init()  {
	//fmt.Println("TEST STATUSES", errorStatuses)
}

func (t *TyphoonTask) IsMaxRetry() bool {
	status := false
	if t.Fetcher.IsMaxFailedRetry() {
		status = true
	} else if t.Fetcher.IsBadStatus() && t.Fetcher.IsResponseRetry() {
		status = true
	} else if !t.Fetcher.IsBadStatus() && t.Processor.IsMaxProcessorRetry(){
		status = true
	} else if t.Fetcher.IsMaxResponseRetry() {
		status = true
	}

	color.Red(`
		DEBUG IsMaxRetry: %t 
		
		IsMaxFailedRetry: %t
		IsBadStatus && IsMaxResponseRetry: %t
			
	`, status, t.Fetcher.IsMaxFailedRetry(), t.Fetcher.IsBadStatus() && t.Processor.IsMaxProcessorRetry())



	return status


}


func (t *TyphoonTask) UpdateRetriesCounter() {

	if t.Fetcher.Response.Code == 599 {
		t.Fetcher.Save.System.Failed += 1
	} else if t.Fetcher.Response.Code == 200 && t.Processor.ErrorResponse {
		t.Processor.Save.System.ProcessorRetries += 1
	} else {
		t.Fetcher.Save.System.Retries += 1
	}



}

func (t *TyphoonTask) IsRetry() bool{
	var status = false

	if t.Fetcher.IsBadStatus() && t.Fetcher.IsFailedRetry() {
		status = true
	} else if t.Fetcher.IsBadStatus() && t.Fetcher.IsResponseRetry() {
		status = true
	} else if t.Fetcher.IsBadStatus() && t.Processor.IsMaxProcessorRetry() {
		status = true
	}


	color.Red(`
		DEBUG FETCHER RESPONSE. Task id %s
			
			Is retry? - %t
			Is bad Status? - %t
			Is failed retry? - %t
			Is max processor retry? - %t
			Is response retry? - %t
			Is max processor retry? - %t
			
			t.Fetcher.Save.System.Failed - %d
			t.Fetcher.Save.System.Retries - %d
	`, t.Taskid,
		status,
		t.Fetcher.IsBadStatus(),
		t.Fetcher.IsFailedRetry(),
		t.Processor.IsMaxProcessorRetry(),
		t.Fetcher.IsResponseRetry(),
		t.Processor.IsMaxProcessorRetry(),
		t.Fetcher.Save.System.Failed,
		t.Fetcher.Save.System.Retries,
	)


	return status
}
