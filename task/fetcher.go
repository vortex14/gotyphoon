package task


type FetcherTask struct {
	Proxy string `json:"proxy"`
	Method  string `json:"method" default:"GET"`
	Timeout int    `json:"timeout" default:"30"`
	MaxRetries int `json:"max_retries" default:"15"`
	MaxFailed int `json:"max_failed" default:"5"`
	Headers map[string]string `json:"headers"`
	Cookies interface{} `json:"cookies"`
	Auth    map[string]string `json:"auth"`
	IsProxyRequired bool   `json:"is_proxy_required"`
	IsResponseCache bool   `json:"is_response_cache"`
	Strategy        string `json:"strategy"`
	Save            struct {
		Project map[string]string `json:"project"`
		System struct {
			Failed int `json:"failed"`
			Retries      int `json:"retries"`
			StatusCode   int `json:"status_code"`
			AddedAt      int `json:"added_at"`
			RetriesDelay int `json:"retries_delay"`
			ExecuteAt    int `json:"execute_at"`
			Exception    struct {
				Type            interface{} `json:"type"`
				Message         string      `json:"message"`
				ErrorDefinition interface{} `json:"error_definition"`
			} `json:"exception"`
		} `json:"system"`
	} `json:"save"`
	Data              interface{} `json:"data"`
	JSON              interface{} `json:"json"`
	Stream            bool        `json:"stream"`
	UserAgentRequired bool        `json:"user_agent_required"`
	ForceUpdate       bool        `json:"force_update" default:"true"`
	LinesLimit        interface{} `json:"lines_limit"`
	Response          struct {
		Content string `json:"content"`
		Code    int    `json:"code"`
		Headers map[string]string `json:"headers"`
		Cookies string `json:"cookies"`
		URL     string `json:"url"`
		OrigURL string `json:"orig_url"`
	} `json:"response"`

}

type Statuses map[int]bool
type Codes map[string]int

var errorStatuses = Statuses{
	200: false,
	304: false,
	400: true,
	403: true,
	404: false,
	405: true,
	406: true,
	407: true,
	500: true,
	501: true,
	502: true,
	503: true,
	504: true,
	505: true,
	506: true,
	599: true,
}

var errorCodes = Codes{
	"timeout": 100,
	"refused": 200,
	"undefined": 300,
	"prefetch": 400,
}


func (t *FetcherTask) IsBadStatus() bool{
	//fmt.Printf("code: %d, IsBadStatus: %t \n", t.Response.Code, errorStatuses[t.Response.Code])
	return errorStatuses[t.Response.Code]
}


func (t *FetcherTask) IsFailedRetry() bool{

	return t.MaxFailed > t.Save.System.Failed && t.Response.Code == 599
}

func (t *FetcherTask) IsResponseRetry() bool{

	return t.MaxRetries > t.Save.System.Retries && t.Response.Code != 599
}

func (t *FetcherTask) IsMaxFailedRetry() bool{

	return t.MaxFailed <= t.Save.System.Failed && t.Response.Code == 599
}

func (t *FetcherTask) IsMaxResponseRetry() bool{

	return t.MaxRetries <= t.Save.System.Retries && t.Response.Code != 599
}
