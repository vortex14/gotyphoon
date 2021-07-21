package v1_1

import (
	"github.com/fatih/color"
	interfaces2 "github.com/vortex14/gotyphoon/interfaces"
	utils2 "github.com/vortex14/gotyphoon/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


type Import struct {
	Old string
	New string
}

type ProjectMigrate struct {
	Project interfaces2.Project
	Dir     *interfaces2.FileObject
}

func (m *ProjectMigrate) MigrateV11()  {
	m.migrateComponents()
	m.CreateDockerComposeConfigs()
	u := utils2.Utils{}
	_, dataT := u.GetGoTemplate(&interfaces2.FileObject{Path: m.Dir.Path, Name: "component.gopy"})

	_, confT := u.GetGoTemplate(&interfaces2.FileObject{Path: m.Dir.Path, Name: "config.goyaml"})
	goTemplate := interfaces2.GoTemplate{
		Source: confT,
		ExportPath: "config.local.yaml",
		Data: map[string]string{
			"projectName": m.Project.GetName(),
			"nsqdAdd": "localhost:4150",
			"redisHost": "localhost",
			"redisPort": "6380",
			"mongoHost": "localhost",
		},
	}

	u.GoRunTemplate(&goTemplate)

	if _, err := os.Stat("fetcher.py"); os.IsNotExist(err) {
		fetcherConfig := map[string]string{
			"component":   "fetcher",
			"executeFile": "fetcher",
			"componentClass": "Fetcher",
			"apiClass": "FetcherApi",
		}
		goTemplateFetcher := interfaces2.GoTemplate{
			Source: dataT,
			ExportPath: "fetcher.py",
			Data: fetcherConfig,
		}


		u.GoRunTemplate(&goTemplateFetcher)
	}

	if _, err := os.Stat("processor.py"); os.IsNotExist(err) {
		processorConfig := map[string]string{
			"component":   "processor",
			"executeFile": "processor",
			"componentClass": "Processor",
			"apiClass": "ProcessorApi",
		}

		goTemplateProcessor := interfaces2.GoTemplate{
			Source: dataT,
			ExportPath: "processor.py",
			Data: processorConfig,
		}


		u.GoRunTemplate(&goTemplateProcessor)
	}

	if _, err := os.Stat("donor.py"); os.IsNotExist(err) {

		donorConfig := map[string]string{
			"component":   "donor",
			"executeFile": "donor",
			"componentClass": "Donor",
			"apiClass": "DonorApi",
		}

		goTemplateDonor := interfaces2.GoTemplate{
			Source: dataT,
			ExportPath: "donor.py",
			Data: donorConfig,
		}


		u.GoRunTemplate(&goTemplateDonor)

	}


	if _, err := os.Stat("scheduler.py"); os.IsNotExist(err) {

		SchedulerConfig := map[string]string{
			"component":   "scheduler",
			"executeFile": "scheduler",
			"componentClass": "Scheduler",
			"apiClass": "SchedulerApi",
		}
		goTemplateScheduler := interfaces2.GoTemplate{
			Source: dataT,
			ExportPath: "scheduler.py",
			Data: SchedulerConfig,
		}


		u.GoRunTemplate(&goTemplateScheduler)

	}


	if _, err := os.Stat("result_transporter.py"); os.IsNotExist(err) {

		rtConfig := map[string]string{
			"component":   "result_transporter",
			"executeFile": "resulttransporter",
			"componentClass": "ResultTransporter",
			"apiClass": "ResultWorkerApi",
		}

		goTemplateTransporter := interfaces2.GoTemplate{
			Source: dataT,
			ExportPath: "result_transporter.py",
			Data: rtConfig,
		}


		u.GoRunTemplate(&goTemplateTransporter)

	}


	_ = filepath.Walk("project/fetcher", m.VisitAndReplace)

	_ = filepath.Walk("project/processor", m.VisitAndReplace)

	_ = filepath.Walk("project/result_transporter", m.VisitAndReplace)

	_ = filepath.Walk("project/donor", m.VisitAndReplace)

	_ = filepath.Walk("project/scheduler", m.VisitAndReplace)
	color.Yellow("Migrated.")
	return
}

func (m *ProjectMigrate) CreateDockerComposeConfigs()  {
	u := utils2.Utils{}
	_, dataTDockerLocal := u.GetGoTemplate(&interfaces2.FileObject{Path: m.Dir.Path, Name: "docker-compose.local.goyaml"})

	dataConfig := map[string]string{
		"projectName": m.Project.GetName(),
		"tag": m.Project.GetTag(),
	}

	goTemplateComposeLocal := interfaces2.GoTemplate{
		Source: dataTDockerLocal,
		ExportPath: "docker-compose.local.yaml",
		Data: dataConfig,
	}


	u.GoRunTemplate(&goTemplateComposeLocal)
}

func (m *ProjectMigrate) migrateComponents()  {
	if _, err := os.Stat("project"); os.IsNotExist(err) {
		_ = os.Mkdir("project", 755)
	}

	if _, err := os.Stat("fetcher"); !os.IsNotExist(err) {
		_ = os.Rename("fetcher", "project/fetcher")
	}

	if _, err := os.Stat("processor"); !os.IsNotExist(err) {
		_ = os.Rename("processor", "project/processor")
	}

	if _, err := os.Stat("donor"); !os.IsNotExist(err) {
		_ = os.Rename("donor", "project/donor")
	}

	if _, err := os.Stat("result_transporter"); !os.IsNotExist(err) {
		_ = os.Rename("result_transporter", "project/result_transporter")
	}

	if _, err := os.Stat("scheduler"); !os.IsNotExist(err) {
		_ = os.Rename("scheduler", "project/scheduler")
	}
}

//func (m *ProjectMigrate) getComponentTemplate() (error error, data string)  {
//	box := packr.NewBox(".")
//
//	data, err := box.FindString("component.gopy")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return err, data
//}

//func (m *ProjectMigrate) getConfigTemplate() (error error, data string) {
//	box := packr.NewBox(".")
//
//	data, err := box.FindString("config.goyaml")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return err, data
//}


func (m *ProjectMigrate) VisitAndReplace(path string, fi os.FileInfo, err error) error {

	var newImports [] *Import

	newImports = append(newImports, &Import{
		Old: "from executions.strategies",
		New: "from typhoon.components.fetcher.executions.strategies",
	},
	&Import{
		Old: "from responses.strategies",
		New: "from typhoon.components.fetcher.responses.strategies",
	},
	&Import{
		Old: "from executions.strategies.base_pre_fetch",
		New: "from typhoon.components.fetcher.executions.strategies.base_pre_fetch",
	},
	&Import{
		Old: "from executable.pipeline_group",
		New: "from typhoon.components.processor.executable.pipeline_group",
	},
	&Import{
		Old: "from executable.base_handler",
		New: "from typhoon.components.processor.executable.base_handler",
	},
	&Import{
		Old: "from project.executable",
		New: "from project.processor.executable",
	},
	&Import{
		Old: "from project.resource",
		New: "from project.processor.resource",
	},
	&Import{
		Old: "from executable.text_pipelines.base_pipeline import BasePipeline",
		New: "from typhoon.components.processor.executable.text_pipelines.base_pipeline import BasePipeline",
	},
	&Import{
		Old: "from project.consumers",
		New: "from project.result_transporter.consumers",
	},
	&Import{
		Old: "from executions.base_consumer",
		New: "from typhoon.components.result_transporter.executions.base_consumer",
	},
	&Import{
		Old: "from project.result_transporter.consumers import exceptions",
		New: "from typhoon.extensions.result_transporter import exceptions",
	},
	&Import{
		Old: "from project.v1",
		New: "from project.donor.v1",
	},
	)

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	if strings.Contains(path, ".pyc") || strings.Contains(path, ".json"){
		return nil
	}

	data, err := ioutil.ReadFile(path)
	newContents := string(data)
	if err != nil {
		panic(err)
	}

	for _, v := range newImports {
		newContents = strings.Replace(newContents, v.Old, v.New, -1)
	}
	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}

	return nil
}

func (m *ProjectMigrate) CreateTransporterManifest()  {
	color.Yellow("Creating manifest for transporter ...")

	//dir := &components.Directory{Path: "project/result_transporter/consumers"}
	//dataDir := dir.GetDataFromDirectory("project/result_transporter/consumers")
	//for _, v := range dataDir {
	//	color.Red("k %s, v %s", v.Path, v.Type)
	//}
}