package typhoon

import (
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


type Cluster struct {
	Name              string
	Config            string
	Description       string
	Typhoon           interfaces.ClusterLabel
	Meta              interfaces.ClusterMeta
	clusterPath       string
	clusterConfigPath string
	Projects          []*interfaces.ClusterProject
}

func (c *Cluster) GetClusterConfigPath() string {
	settings := c.GetEnvSettings()
	clusterConfigPath := settings.Clusters + "/" + c.Name + "/" + c.Config
	return clusterConfigPath
}

func (c *Cluster) Create()  {
	settings := c.GetEnvSettings()

	if len(settings.Clusters) == 0 {
		color.Red("Cluster path not found. Need set env variable: export TYPHOON_CLUSTERS=")
		os.Exit(1)
	}
	c.clusterPath = settings.Clusters + "/" + c.Name

	c.clusterConfigPath = c.clusterPath + "/" + c.Config

	if _, err := os.Stat(c.clusterPath); !os.IsNotExist(err) {
		// path/to/whatever exists
		color.Yellow("Cluster path (%s) already exists!", c.clusterPath)


		if _, err := os.Stat(c.clusterConfigPath); !os.IsNotExist(err) {
			color.Red("Cluster config already exists!")
			os.Exit(1)
		}

	} else {
		color.Yellow("Create a new %s cluster. Cluster dir: %s", c.Name, settings.Clusters + "/" + c.Name)
		errClusterDir := os.MkdirAll(c.clusterPath, os.ModePerm)
		if errClusterDir != nil {
			color.Red("%s", errClusterDir)
			os.Exit(1)
		}
	}

	emptyConfig, _ := yaml.Marshal(c)

	u := utils.Utils{}
	_ = u.DumpToFile(&interfaces.FileObject{
		Data: string(emptyConfig),
		Path: c.clusterConfigPath,

	})
}

func (c *Cluster) getAllProjects(settings *environment.Settings) []string {
	files, _ := ioutil.ReadDir(settings.Projects)

	var allProjects []string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if f.Name()[:1] == "." {
			continue
		}

		if _, err := os.Stat(settings.Projects + "/" + f.Name() + "/.git"); os.IsNotExist(err) {
			continue
		}

		allProjects = append(allProjects, f.Name())
	}

	return allProjects

}

func (c *Cluster) ChangeClusterBranch() {

}

func (c *Cluster) findProjectsByTerm(term string, projects []string) [] string{
	var foundProjects []string
	for _, project := range projects {
		if strings.Contains(strings.ToLower(project), strings.ToLower(term)) {
			foundProjects = append(foundProjects, project)
		}
	}
	return foundProjects
}

func (c *Cluster) renderClusterList(table *tview.Table, projects []string, settings *environment.Settings)  {
	table.SetCell(0, 0,
		tview.NewTableCell("№").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).SetMaxWidth(10))

	table.SetCell(0, 1,
		tview.NewTableCell("Typhoon Projects").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).SetMaxWidth(40))

	for r := 0; r < len(projects); r++ {
		for c := 0; c < 2; c++ {
			color := tcell.ColorWhite
			rowNumber := r + 1
			var value string
			projectName := projects[r]

			if c == 0 {
				value = strconv.Itoa(rowNumber)
			} else if c == 1 {
				value = projectName
			}
			table.SetCell(rowNumber, c, tview.NewTableCell(value).SetTextColor(color))
		}
	}
}

func (c *Cluster) removeProject(projectName string)  {
	for i, project := range c.Projects {
		if project.Name == projectName {
			c.Projects = append(c.Projects[:i], c.Projects[i+1:]...)
			break
		}
	}
}

func (c *Cluster) isExistProject(projectName string) bool {
	status := false
	for _, project := range c.Projects {
		if project.Name == projectName {
			status = true
			break
		}
	}
	return status
}

func (c *Cluster) SaveConfig()  {
	settings := c.GetEnvSettings()
	c.clusterConfigPath = settings.Clusters + "/" + c.Name + "/" + c.Config
	data, _ := yaml.Marshal(&c)
	u := utils.Utils{}
	_ = u.DumpToFile(&interfaces.FileObject{
		Data: string(data),
		Path: c.clusterConfigPath,

	})

	//color.Yellow("Config %s updated. ", c.clusterConfigPath)
}

func (c *Cluster) LoadConfig(settings *environment.Settings) *Cluster {
	configCluster := settings.Clusters + "/" + c.Name + "/" + c.Config
	c.clusterPath = settings.Clusters + "/" + c.Name
	dat, err := ioutil.ReadFile(configCluster)
	if err != nil && len(dat) > 0 {
		color.Red("%s", err)
		os.Exit(1)
	}
	var configClusterLoad Cluster
	errDecode := yaml.Unmarshal(dat, &configClusterLoad)
	if errDecode != nil {
		color.Red("%s", errDecode)
		os.Exit(1)
	}
	c.Description = configClusterLoad.Description
	c.Projects = configClusterLoad.Projects
	c.Meta = configClusterLoad.Meta
	c.Typhoon = configClusterLoad.Typhoon
	if len(c.Typhoon.Kind) == 0 {
		c.Typhoon.Kind = "Cluster"
	}

	return &configClusterLoad

}

func (c *Cluster) GetConfigName() string {
	return c.Config
}

func (c *Cluster) GetProjects() []*interfaces.ClusterProject {
	settings := c.GetEnvSettings()
	clusterConfig := c.LoadConfig(settings)
	return clusterConfig.Projects
}

func (c *Cluster) GetEnvSettings() *environment.Settings {
	env := environment.Environment{}
	env.Load()
	_, settings := env.GetSettings()
	return settings
}

func (c *Cluster) GetName() string  {
	return c.Name
}

func (c *Cluster) renderTable(data [][]string)  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Sign", "Rating"})
	table.AppendBulk(data)
	table.Render()
}

func (c *Cluster) Show()  {
	settings := c.GetEnvSettings()

	if len(settings.Clusters) == 0 {
		color.Red("Cluster Path not found !")
		os.Exit(1)
	} else {
		color.Green("Cluster Path: %s", settings.Clusters)
	}

	files, _ := ioutil.ReadDir(settings.Clusters)

	var allClusters [][]string
	for n, f := range files {
		if !f.IsDir() {
			continue
		}
		if f.Name()[:1] == "." {
			continue
		}

		allClusters = append(allClusters, []string{strconv.Itoa(n), f.Name()})
	}

	u := utils.Utils{}

	header := []string{"№", "Name"}

	u.RenderTableOutput(header, allClusters)

}

func (c *Cluster) Deploy() {
	var tableData [][]string
	settings := c.GetEnvSettings()
	color.Green("Deploy ... to %s", settings.Gitlab)
	gitlabClient, err := gitlab.NewClient(settings.GitlabToken, gitlab.WithBaseURL(settings.Gitlab))
	if err != nil {
		color.Red("%s", err)
	}

	projects, _, err := gitlabClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		ListOptions:              gitlab.ListOptions{
			PerPage: 100,
			Page: 1,
		},
	})

	if err != nil {
		color.Red("%s", err)
		return
	}

	for i, project := range projects {
		tableData = append(tableData, []string{string(i), project.Name, project.WebURL})
	}

	c.renderTable(tableData)

}

func (c *Cluster) GetMeta() *interfaces.ClusterMeta {
	return &c.Meta
}

func (c *Cluster) Add()  {
	settings := c.GetEnvSettings()

	if len(settings.Clusters) == 0 {
		color.Red("Cluster path not found. Need set env variable: export TYPHOON_CLUSTERS=")
		os.Exit(1)
	}
	if len(settings.Projects) == 0 {
		color.Red("Project path not found. Need set env variable: export TYPHOON_PROJECTS=")
		os.Exit(1)
	}

	if _, err := os.Stat(settings.Clusters + "/" + c.Name); os.IsNotExist(err) {
		color.Red("Cluster %s not found", c.Name)
		color.Green("typhoon cluster create -n %s", c.Name)
		os.Exit(1)
	}

	_ = c.LoadConfig(settings)
	configData, _ := yaml.Marshal(&c)

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	var selectedBranch string
	var selectedProjectName string
	var selectedRemote string
	var selectedGitUrl string

	c.clusterConfigPath = settings.Clusters + "/" + c.Name + "/" + c.Config

	allProjects := c.getAllProjects(settings)

	c.renderClusterList(table, allProjects, settings)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true)

	grid.
		AddItem(table, 1, 0, 1, 2, 0, 0, false)

	grid.
		AddItem(table, 1, 0, 1, 2, 0, 100, false)

	inputField := tview.NewInputField().
		SetLabel("filter by project name: ").
		SetFieldBackgroundColor(tcell.ColorSkyblue).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldWidth(30).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyTab {
				table.SetSelectable(true, true)
				app.SetFocus(table)
			}
		}).
		SetChangedFunc(func(term string) {
			table.Clear()
			projects := c.findProjectsByTerm(term, allProjects)
			c.renderClusterList(table, projects, settings)
		})


	pages := tview.NewPages().
		AddPage("main", grid, true, true)

	box := tview.NewTextView().
		SetDynamicColors(true).
		//SetRegions(true).
		SetText(string(configData)).
		SetChangedFunc(func() {
			app.Draw()
		})

	box.Box.SetBorder(true)
	box.Box.SetTitle(c.Config)

	configsList := tview.NewList().SetSelectedFunc(func(i int, selectedConfig string, additional string, r rune) {
		c.Projects = append(c.Projects, &interfaces.ClusterProject{
			Name: selectedProjectName,
			Config: selectedConfig,
			Labels: interfaces.ClusterProjectLabels{
				Git: interfaces.GitLabel{
					Url:    selectedGitUrl,
					Remote: selectedRemote,
					Branch: selectedBranch,
				}},
		})
		data, _ := yaml.Marshal(c)
		c.SaveConfig()
		box.SetText(string(data))
		pages.SwitchToPage("main")
		app.SetFocus(table)
	})

	RemoteList := tview.NewList().SetSelectedFunc(func(i int, remote string, git string, r rune) {
		selectedRemote = remote
		selectedGitUrl = git
		if configsList.GetItemCount() == 0 {
			pages.SwitchToPage("remotes")
		} else {
			pages.SwitchToPage("configs")
			app.SetFocus(configsList)
		}
	})



	pages.AddPage("remotes", RemoteList, true, false)

	pages.AddPage("configs", configsList, true, false)

	RemoteList.SetRect(100,100,100,200)

	grid.AddItem(inputField, 0, 0, 1, 2, 0, 0, true)

	grid.AddItem(box, 0, 2,2,2,3,0, false)

	table.Select(1, 1).SetFixed(1, 3).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		} else if key == tcell.KeyTab {
			app.SetFocus(inputField)
			table.SetSelectable(false, false)
		} else if key == tcell.KeyCtrlW {
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		projectName := cell.Text
		projectPath := settings.Projects + "/" + projectName
		if cell.Color == tcell.ColorRed || c.isExistProject(projectName) {
			cell.SetTextColor(tcell.ColorWhite)
			
			c.removeProject(projectName)
			c.SaveConfig()
			data, _ := yaml.Marshal(c)
			box.SetText(string(data))
			
		} else {

			project := Project{
				Name: projectName,
				Path: projectPath,
			}

			remotes, err := project.GetRemotes()
			configs := project.GetConfigs()

			if err != nil {
				cell.BackgroundColor = tcell.ColorSkyblue
				return
			}

			RemoteList.Clear()

			for _, remote := range remotes {
				RemoteList.AddItem(remote.Config().Name, remote.Config().URLs[0] ,'*' , nil)
			}

			configsList.Clear()
			for _, config := range configs {
				configsList.AddItem(config, "" ,'*' , nil)
			}

 			errH, projectBranch := project.GetBranch()
 			if errH != nil {
 				cell.BackgroundColor = tcell.ColorSkyblue
				return
			}

			selectedProjectName = cell.Text
			selectedBranch = projectBranch

			cell.SetTextColor(tcell.ColorRed)

			pages.SwitchToPage("remotes")
		}
	})




	if err := app.SetRoot(pages, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}

}