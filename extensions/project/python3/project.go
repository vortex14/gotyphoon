package python3

import (
	"context"
	"fmt"
	lessWatcher "github.com/radovskyb/watcher"
	BaseWatcher "github.com/vortex14/gotyphoon/elements/models/watcher"
	"github.com/vortex14/gotyphoon/extensions/models/watcher"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/folder"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/environment"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	tyLog "github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/migrates/v1.1"
	"github.com/vortex14/gotyphoon/services"
	"github.com/vortex14/gotyphoon/utils"
)

type Services struct {
	Mongo map[string]mongo.Service
	Redis map[string]redis.Service
}

type TestMongo struct {
}

type Project struct {
	forms.Project

	Watcher *watcher.EventFSLess
}

func (p *Project) GetService(name string) interfaces.Service {
	switch name {
	case interfaces.NSQ:
		return p.Services.Collections.Nsq
		//case interfaces.MONGO:
		//	return p.Services.Collections.Mongo

	}
	return nil
}

func (p *Project) RunTestServices() {

	typhoonServices := services.Services{Project: p}
	typhoonServices.RunTestServices()
}

func (p *Project) Run() interfaces.Project {
	p.Construct(func() {
		p.Folder = &folder.Folder{Path: p.GetProjectPath()}

		p.CheckProject()

		p.LOG = tyLog.New(tyLog.D{"project": p.Name})

		p.Add()
		go p.Watch()

		p.LOG.Info(p.Folder.Path)

		if !p.Folder.IsExist("typhoon") {
			p.LOG.Info("init typhoon symlink path")
			_ = p.CreateSymbolicLink()
		}

		//
		//
		//color.Magenta("start components")
		p.LOG.Info("start components ...")
		//p.AddPromise()
		p.StartComponents(true)
		////
		//p.AddPromise()
		//go p.task.Run()
		//
		//c := make(chan os.Signal, 1)
		//signal.Notify(c, os.Interrupt)
		//go p.Watch()
		////go Watch(&task.wg, typhoonComponent, project.GetConfigFile())
		//sig := <-c
		//fmt.Printf("Got %s signal. Aborting...\n", sig)
		//p.AddPromise()
		//go p.Close()
		//p.task.Stop()
	})

	return p

}

func (p *Project) CreateProject() {
	color.Yellow("creating project...")
	u := utils.Utils{}
	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1/project",
	}

	err := u.CopyDir(p.Name, fileObject)

	if utils.NotNill(err) {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	gitIgnore := &interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: ".gitignore",
	}
	errCopyIgnore := u.CopyFile(p.Name+"/.gitignore", gitIgnore)
	if errCopyIgnore != nil {
		color.Red("Error copy %s", err)
	}

	_, confT := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config.goyaml",
	})
	goTemplate := interfaces.GoTemplate{
		Source:     confT,
		ExportPath: p.Name + "/config.local.yaml",
		Data: map[string]string{
			"projectName": p.Name,
			"nsqdAdd":     "localhost:4150",
			"redisHost":   "localhost",
			"mongoHost":   "localhost",
			"redisPort":   "6379",
			"debug":       "true",
		},
	}

	_ = u.GoRunTemplate(&goTemplate)
	goTemplateCompose := interfaces.GoTemplate{
		Source:     confT,
		ExportPath: p.Name + "/config.prod.yaml",
		Data: map[string]string{
			"projectName": p.Name,
			"nsqdAdd":     "nsqd:4150",
			"redisHost":   "redis",
			"redisPort":   "6379",
		},
	}

	_ = u.GoRunTemplate(&goTemplateCompose)
	//color.Green("Teplate status: %b", status)

	_, dataTDockerLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1", Name: "docker-compose.local.goyaml"})

	dataConfig := map[string]string{
		"projectName": p.GetName(),
		"tag":         p.GetTag(),
	}

	goTemplateComposeLocal := interfaces.GoTemplate{
		Source:     dataTDockerLocal,
		ExportPath: p.Name + "/docker-compose.local.yaml",
		Data:       dataConfig,
	}

	u.GoRunTemplate(&goTemplateComposeLocal)
	color.Green("Project %s created !", p.Name)

}

func (p *Project) BuildCIResources() {
	color.Green("Build CI Resources for %s !", p.Name)
	u := utils.Utils{}
	_, confCi := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: ".gitlab-ci.yml",
	})
	goTemplate := interfaces.GoTemplate{
		Source:     confCi,
		ExportPath: ".gitlab-ci.yml",
	}

	_ = u.GoRunTemplate(&goTemplate)

	_, dockerFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "Dockerfile",
	})
	goTemplateDocker := interfaces.GoTemplate{
		Source:     dockerFile,
		ExportPath: "Dockerfile",
		Data: map[string]string{
			"TYPHOON_IMAGE": p.Version,
		},
	}

	_ = u.GoRunTemplate(&goTemplateDocker)

	_, helmFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "helm-review-values.yml",
	})
	goTemplateHelmValues := interfaces.GoTemplate{
		Source:     helmFile,
		ExportPath: "helm-review-values.yml",
	}

	_ = u.GoRunTemplate(&goTemplateHelmValues)

	_, configFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config-stage.goyaml",
	})
	goTemplateConfig := interfaces.GoTemplate{
		Source:     configFile,
		ExportPath: "config.kube-stage.yaml",
		Data: map[string]string{
			"projectName": p.GetName(),
		},
	}

	_ = u.GoRunTemplate(&goTemplateConfig)

}

func (p *Project) GetEnvSettings() *environment.Settings {
	return p.EnvSettings
}

func (p *Project) Watch() {
	color.Green("watch for project ..")

	p.LOG.Debug(p.GetProjectPath())

	p.Watcher = &watcher.EventFSLess{
		Watcher: BaseWatcher.Watcher{Path: p.GetProjectPath()},
		Callback: func(log interfaces.LoggerInterface, event *lessWatcher.Event) {
			//println("111!!! !!! >>>>> >>  >> > >")
			//log.Error(event.)

			componentChanged := ""

			for _, component := range p.SelectedComponent {
				if strings.Contains(event.Path, strings.ToLower(component)) {
					color.Yellow("reloading %s ... !", component)
					componentChanged = component
					break
				}

			}

			if _, ok := p.Components.ActiveComponents[componentChanged]; ok {

				color.Yellow("Reload %s ...", componentChanged)
				color.Yellow("event %+v", event)
				component := p.Components.ActiveComponents[componentChanged]

				//p.AddPromise()
				component.Restart(p)

				// "example" is not in the map
			} else {
				color.Yellow("%s isn't running", componentChanged)
			}

		},
		Timeout: 100,
	}
	p.LOG.Warning("Start watching ...")
	p.Watcher.Watch()
	p.LOG.Warning("Stop watching ...")

	//watcher, _ = fsnotify.NewWatcher()
	//defer watcher.Close()
	//
	//if err := filepath.Walk("project", watchDirTeet); err != nil {
	//	fmt.Println("ERROR", err)
	//}
	//
	//done := make(chan bool)
	//
	//go func() {
	//	for {
	//		select {
	//		case event := <-watcher.Events:
	//
	//			if strings.Contains(event.Name, ".pyc") {
	//				continue
	//			}
	//
	//			if strings.Contains(event.String(), "CHMOD") {
	//				continue
	//			}
	//
	//			if strings.Contains(event.Name, ".py~") {
	//				continue
	//			}
	//
	//			if strings.Contains(event.Name, "__pycache__") {
	//				continue
	//			}
	//
	//			componentChanged := ""
	//
	//			for _, component := range p.SelectedComponent {
	//				if strings.Contains(event.Name, strings.ToLower(component)) {
	//					color.Yellow("reloading %s ... !", component)
	//					componentChanged = component
	//					break
	//				}
	//
	//			}
	//
	//			if _, ok := p.components.ActiveComponents[componentChanged]; ok {
	//
	//				color.Yellow("Reload %s ...", componentChanged)
	//				color.Yellow("event %+v", event)
	//				component := p.components.ActiveComponents[componentChanged]
	//
	//				//p.AddPromise()
	//				go component.Restart(p)
	//
	//				// "example" is not in the map
	//			} else {
	//				color.Yellow("%s isn't running", componentChanged)
	//			}
	//
	//			//
	//
	//			//p.AddPromise()
	//			//go component.Restart(p)
	//
	//			//go component.Restart(p)
	//
	//			//initComponent(wg, tcomponents, componentChanged, configFile)
	//
	//			// watch for errors
	//		case err := <-watcher.Errors:
	//			color.Red("ERROR---->", err)
	//		}
	//	}
	//}()
	//
	//<-done
}

func (p *Project) Close() {
	// Close Watcher Promise
	p.Done()

	color.Yellow("close project ...")
	for _, component := range p.Components.ActiveComponents {
		if component.IsActive() {
			p.Add()
			go component.Close(p)
		}

	}
	color.Yellow("await done...")
	p.Await()
	color.Yellow("await done !")

}

func (p *Project) Down() {
	p.LoadConfig()
	commandDropProject := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", p.GetName())
	color.Red("Running: %s: ", commandDropProject)
	ctxP, cancelP := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelP()

	if err := exec.CommandContext(ctxP, "bash", "-c", commandDropProject).Run(); err != nil {
		color.Yellow("Project components killed!")
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}

	commandDropTyphoon := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", "typhoon")
	ctxT, cancelT := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelT()

	if err := exec.CommandContext(ctxT, "bash", "-c", commandDropTyphoon).Run(); err != nil {
		color.Red("%s", err.Error())
	}
}

func (p *Project) Migrate() {

	color.Yellow("Migrate project to %s !", p.GetVersion())

	if p.Version == "v1.1" {
		prMigrates := v1_1.ProjectMigrate{
			Project: p,
			Dir: &interfaces.FileObject{
				Path: "../builders/v1.1",
			},
		}
		prMigrates.MigrateV11()
	}
}

func (p *Project) Build() {
	color.Yellow("builder run... options %+v", p.BuilderOptions)
}

func (p *Project) RunQueues() {
	if len(p.SelectedComponent) == 0 {
		color.Red("No set components for project")
		return
	}
	p.Services.RunNSQ()
}

func (p *Project) StartComponents(promise bool) {

	fmt.Printf(tyLog.DOG)

	if p.Components.ActiveComponents == nil {
		p.InitComponents()
	}

	//if promise {
	//	defer p.PromiseDone()
	//}

	for _, componentName := range p.SelectedComponent {
		p.Add()
		go p.Components.ActiveComponents[componentName].Start(p)
	}

	p.LOG.Debug("Components started")
}

func (p *Project) InitComponents() {
	p.Components.ActiveComponents = make(map[string]interfaces.Component)
	p.Name = p.Config.ProjectName
	for _, componentName := range p.SelectedComponent {

		var componentFileName string

		switch componentName {
		case interfaces.FETCHER:
			componentFileName = interfaces.TYPHOON2PYTHON2FETCHER
		case interfaces.PROCESSOR:
			componentFileName = interfaces.TYPHOON2PYTHON2PROCESSOR
		case interfaces.DONOR:
			componentFileName = interfaces.TYPHOON2PYTHON2DONOR
		case interfaces.TRANSPORTER:
			componentFileName = interfaces.TYPHOON2PYTHON2TRANSPORTER
		case interfaces.SCHEDULER:
			componentFileName = interfaces.TYPHOON2PYTHON2SCHEDULER

		}

		p.Components.ActiveComponents[componentName] = &Component{
			Component: forms.Component{
				MetaInfo: label.MetaInfo{Name: componentName},
				Language: interfaces.PYTHON,
				FileExt:  fmt.Sprintf("%s.py", componentFileName),
			},
		}
	}
}

func (p *Project) CreateSymbolicLink() error {
	env := &environment.Environment{}
	_, settings := env.GetSettings()

	linkTyphoonPath := fmt.Sprintf("%s/pytyphoon/typhoon", settings.Path)
	color.Yellow("TYPHOON_PATH=%s", settings.Path)
	directLink := filepath.Join(p.GetProjectPath(), "typhoon")
	color.Yellow(directLink)
	err := os.Symlink(linkTyphoonPath, directLink)

	if err != nil {
		fmt.Printf("err %s", err)
	}

	return nil
}

func (p *Project) CheckProject() {

	var status = true

	for _, componentName := range p.SelectedComponent {
		component := &Component{
			Component: forms.Component{
				MetaInfo:    label.MetaInfo{Name: componentName},
				ProjectPath: p.Folder.Path,
			},
		}
		color.Yellow("checking: %s...", componentName)

		if !component.CheckComponent() {
			status = false
			color.Yellow("%s is: false", componentName)
		}

	}
	p.LoadConfig()

	if !status {
		color.Red("%s : %s", Errors.ProjectNotFound.Error(), p.Path)
		os.Exit(1)
	}

	env := &environment.Environment{}
	_, settings := env.GetSettings()

	if len(settings.Path) == 0 || len(settings.Projects) == 0 {
		color.Red("%s in %s", Errors.ProjectInvalidEnv.Error(), env.ProfilePath)
		os.Exit(1)
	}

	p.EnvSettings = settings

}
