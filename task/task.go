package task

import (
	"github.com/fatih/color"
)


type TyphoonTask struct {
	Fetcher           FetcherTask   `json:"fetcher"`
	Processor         ProcessorTask `json:"processor"`
	Scheduler         SchedulerTask `json:"scheduler"`
	ResultTransporter TransporterTask `json:"result_transporter"`
	Priority          int           `json:"priority" default:"3" fake:"{randomstring:[3]}"`
	URL               string        `json:"url" default:"https://httpstat.us/200" fake:"{url}"`
	Taskid            string        `json:"taskid" default:"task-id" fake:"{uuid}"`
	ProjectName string `json:"project_name" fake:"skip"`
	//msg *nsq.Message
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

