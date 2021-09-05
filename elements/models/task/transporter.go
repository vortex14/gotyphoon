package task

type TransporterTask struct {
	Consumer string `json:"consumer" fake:"{randomstring:[exceptions]}"`
	ForceUpdate bool `json:"force_update" fake:"{randomstring:[true]}"`
	Age int `json:"age" fake:"{randomstring:[86400]}"`
	Save interface{} `json:"save"`
	Result interface{} `json:"result"`
}