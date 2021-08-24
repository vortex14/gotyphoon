package interfaces


type GrafanaInterface interface {
	ImportGrafanaConfig()
	RemoveGrafanaDashboard()
	CreateBaseGrafanaConfig()
	CreateGrafanaMonitoringTemplates()
}
