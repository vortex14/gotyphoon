package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/grafana-tools/sdk"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type DashBoard struct {
	Id         string
	FolderId   string
	ConfigName string
	Project    interfaces.Project
	Cluster    interfaces.Cluster
}

type ResponseImportDashboard struct {
	PluginID         string `json:"pluginId"`
	Title            string `json:"title"`
	Imported         bool   `json:"imported"`
	ImportedURI      string `json:"importedUri"`
	ImportedURL      string `json:"importedUrl"`
	Slug             string `json:"slug"`
	DashboardID      int    `json:"dashboardId"`
	FolderID         int    `json:"folderId"`
	ImportedRevision int    `json:"importedRevision"`
	Revision         int    `json:"revision"`
	Description      string `json:"description"`
	Path             string `json:"path"`
	Removed          bool   `json:"removed"`
	Message			 string `json:"message"`
}


type DashboardGrafana struct {
	Message   string `json:"message"`
	Overwrite bool   `json:"overwrite"`
	FolderId int `json:"folderId"`
	Dashboard struct {
		Annotations struct {
			List []struct {
				BuiltIn    int    `json:"builtIn"`
				Datasource string `json:"datasource"`
				Enable     bool   `json:"enable"`
				Hide       bool   `json:"hide"`
				IconColor  string `json:"iconColor"`
				Name       string `json:"name"`
				Type       string `json:"type"`
			} `json:"list"`
		} `json:"annotations"`
		Editable     bool          `json:"editable"`
		GnetID       interface{}   `json:"gnetId"`
		GraphTooltip int           `json:"graphTooltip"`
		Links        []interface{} `json:"links"`
		Panels       interface {} `json:"panels"`
		Refresh       string        `json:"refresh"`
		SchemaVersion int           `json:"schemaVersion"`
		Style         string        `json:"style"`
		Tags          []interface{} `json:"tags"`
		Templating    struct {
			List []interface{} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
		} `json:"timepicker"`
		Timezone string `json:"timezone"`
		Title    string `json:"title"`
		Version  int    `json:"version"`
	} `json:"dashboard"`
}


type Config struct {
	Name string
	Endpoint string
	//token string
	DashBoardUrl string
}

func (d *DashBoard) getClient(configProject *interfaces.ConfigProject) (context.Context, *sdk.Client) {
	settings := d.Project.GetEnvSettings()
	c, _ := sdk.NewClient(settings.GrafanaEndpoint, settings.GrafanaToken, sdk.DefaultHTTPClient)
	ctx := context.Background()
	return ctx, c
}

func (d *DashBoard) GetGrafanaDashboard() *DashboardGrafana {
	var configData DashboardGrafana
	_ = d.Project.LoadConfig()
	rawBoard, _ := ioutil.ReadFile(d.Project.GetProjectPath() + "/" +d.ConfigName)
	_ = json.Unmarshal(rawBoard, &configData)
	return &configData
}

func (d *DashBoard) ImportGrafanaConfigLowLevel(jsonConfig []byte, folderId string) *interfaces.GrafanaConfig {
	configProject := d.Project.LoadConfig()
	settings := d.Project.GetEnvSettings()
	ctx, c := d.getClient(configProject)
	bearer := fmt.Sprintf("Bearer %s", settings.GrafanaToken)
	url := settings.GrafanaEndpoint
	importUrl := fmt.Sprintf("%s/api/dashboards/import", url)
	var configData DashboardGrafana
	_ = json.Unmarshal(jsonConfig, &configData)

	folderID := d.getFolderId(ctx, c, folderId)

	configData.FolderId = folderID
	dashboardName := fmt.Sprintf("Dashboard of %s", configData.Dashboard.Title)
	configData.Dashboard.Title = dashboardName

	requestBody, _ := json.Marshal(configData)
	req, err := http.NewRequest("POST", importUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		color.Red("%s", err.Error())
		return nil
	}
	req.Header.Set("Authorization", bearer)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("%s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response ResponseImportDashboard
	_ = json.Unmarshal(body, &response)
	if !response.Imported {
		color.Red("%+v", response)
		os.Exit(1)
	}

	configDashboard := interfaces.GrafanaConfig{
		Id:           strings.Split(strings.Split(response.ImportedURL, "d/")[1], "/")[0],
		Name:         dashboardName,
		FolderId:     folderId,
		DashboardUrl: settings.GrafanaEndpoint + response.ImportedURL,
	}

	configProject.Grafana = append(configProject.Grafana, configDashboard)

	//color.Yellow("%+v", response)
	configDumpData, err := yaml.Marshal(&configProject)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil
	}
	u := &utils.Utils{}
	err = u.DumpToFile(&interfaces.FileObject{
		Name: configProject.GetConfigName(),
		Data: string(configDumpData),
		Path: configProject.GetConfigName(),
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil
	}
	if response.Imported {
		color.Green("%s created !", dashboardName)
	}

	return &configDashboard
}


func (d *DashBoard) getFolderId(ctx context.Context, c *sdk.Client,folderUID string) int {
	var FolderId int
	if len(folderUID) > 0 && folderUID != "0" {
		data, _ := c.GetFolderByUID(ctx, folderUID)
		if data.ID == 0 {
			color.Red("Folder not found. UUID %s", folderUID)
			os.Exit(1)
			return 0
		}

		FolderId = data.ID
	} else {
		FolderId = sdk.DefaultFolderId
	}

	return FolderId



}
func (d *DashBoard) ImportGrafanaConfig(folderId string) *interfaces.GrafanaConfig {
	_ = d.Project.LoadConfig()
	rawBoard, _ := ioutil.ReadFile(d.Project.GetProjectPath() + "/" +d.ConfigName)

	configDashboard := d.ImportGrafanaConfigLowLevel(rawBoard, folderId)

	return configDashboard
}


func (d *DashBoard) RemoveGrafanaDashboard() (error, *interfaces.GrafanaConfig) {
	configProject := d.Project.LoadConfig()
	grafanaDashboard := d.GetGrafanaDashboard()
	var removedDashboard interfaces.GrafanaConfig
	ctx, c := d.getClient(configProject)
	var dashboardId string
	dashboardName := "Dashboard of " + grafanaDashboard.Dashboard.Title
	for i, dashboard := range configProject.Grafana {
		if dashboard.Name == dashboardName {
			configProject.Grafana = append(configProject.Grafana[:i], configProject.Grafana[i+1:]...)
			dashboardId = dashboard.Id
			removedDashboard = interfaces.GrafanaConfig{
				Name: dashboardName,
				Id: dashboardId,
			}
			break
		}
	}


	_, err := c.DeleteDashboardByUID(ctx, dashboardId)
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}

	color.Green("%s was be removed.", dashboardName)

	configDumpData, err := yaml.Marshal(&configProject)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err, nil
	}
	u := &utils.Utils{}
	err = u.DumpToFile(&interfaces.FileObject{
		Name: d.ConfigName,
		Data: string(configDumpData),
		Path: configProject.GetConfigName(),
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return err, nil
	}

	return nil, &removedDashboard
	//_, data := c.GetAllFolders(ctx)
	//color.Red("%+v", data)
}

func (d *DashBoard) CreateGrafanaMonitoringTemplates()  {
	d.Project.LoadConfig()
	u := utils.Utils{}

	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "grafana-template.gojson",
	}
	validProjectName := strings.ReplaceAll(d.Project.GetName(), "-", "_")
	exportPath := fmt.Sprintf("%s/monitoring-grafana.json", d.Project.GetProjectPath())
	err := u.CopyFileAndReplaceLabel(exportPath, &interfaces.ReplaceLabel{Label: "{{.projectName}}", Value: validProjectName}, fileObject)

	if err != nil {

		color.Red("Error: %s", err)
		os.Exit(0)

	}
	color.Green("Generated Grafana monitoring template for %s", d.Project.GetName())
	color.Yellow("%s", exportPath)
	fmt.Printf("========================  ")
}

func (d *DashBoard) CreateBaseGrafanaConfig()  {
	color.Yellow("Creating base grafana properties into typhoon project config.yaml")
	configProject := d.Project.LoadConfig()
	configProject.Grafana = append(configProject.Grafana, interfaces.GrafanaConfig{
		Name: "Typhoon project dashboard",
		Id: "0000000",
	})

	configDumpData, _ := yaml.Marshal(&configProject)

	u := &utils.Utils{}
	err := u.DumpToFile(&interfaces.FileObject{
		Name: d.Project.GetConfigFile(),
		Data: string(configDumpData),
		Path: configProject.GetConfigName(),
	})

	if err != nil {
		return
	}

	color.Green("%s updated.", d.Project.GetConfigFile())
}

func (d *DashBoard) CreateGrafanaNSQMonitoringTemplates()  {
	d.Project.LoadConfig()
	color.Yellow("Creating NSQ Grafana monitoring templates ...")
	u := utils.Utils{}

	exportPath := fmt.Sprintf("%s/grafana-nsq-monitoring.json", d.Project.GetProjectPath())

	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "grafana-nsq-template.gojson",
	}
	err := u.CopyFileAndReplaceLabel(exportPath, &interfaces.ReplaceLabel{Label: "{{.projectName}}", Value: d.Project.GetName()}, fileObject)

	if err != nil {
		color.Red("Error %s", err)
		os.Exit(0)
	}
	color.Green("Generated Grafana NSQ monitoring template for %s", d.Project.GetName())
}
