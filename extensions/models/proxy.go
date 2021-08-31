package models

type Proxy struct {
	Proxy string `json:"proxy"`
	Agent string `json:"agent"`
	Success bool `json:"success"`
	AgentMobile string `json:"agent_mobile"`
}
