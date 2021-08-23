package helm

import (
	"os"

	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type Resources struct {
	Project interfaces.Project
}

func (r *Resources) RemoveHelmMinikubeManifests()  {
	u := utils.Utils{}
	u.RemoveFiles([]string{
		"helm",
		"helm_delete.sh",
		"helm_deploy.sh",
		"helm_dump.sh",
		"helm_delete.sh",
	})

	color.Green("Removed")
}

func (r *Resources) BuildHelmMinikubeResources()  {
	color.Yellow("build helm minikube resources ...")
	r.Project.LoadConfig()

	u := utils.Utils{}
	fileObject := &interfaces.FileObject{
		Path: "../builders/v1.1/helm/helm",
	}

	err := u.CopyDirAndReplaceLabel("helm", &interfaces.ReplaceLabel{Label: "{{PROJECT_NAME}}", Value: r.Project.GetName()}, fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	_, dataTDeployLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_deploy.gosh"})

	dataConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDeployLocal := interfaces.GoTemplate{
		Source: dataTDeployLocal,
		ExportPath: "helm_deploy.sh",
		Data: dataConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDeployLocal)

	_, dataTDumpLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_dump.gosh"})

	dataDumpConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDumpLocal := interfaces.GoTemplate{
		Source: dataTDumpLocal,
		ExportPath: "helm_dump.sh",
		Data: dataDumpConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDumpLocal)


	_, dataTDeleteLocal := u.GetGoTemplate(&interfaces.FileObject{Path: "../builders/v1.1/helm", Name: "helm_delete.gosh"})

	dataDeleteConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDeleteLocal := interfaces.GoTemplate{
		Source: dataTDeleteLocal,
		ExportPath: "helm_delete.sh",
		Data: dataDeleteConfig,
	}


	u.GoRunTemplate(&goTemplateHelmDeleteLocal)


	if err := os.Chmod("helm_delete.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	if err := os.Chmod("helm_deploy.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	if err := os.Chmod("helm_dump.sh", 0755); err != nil {
		color.Red("%s",err)
	}

	_, confT := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "config.minikube.goyaml",

	})
	goTemplate := interfaces.GoTemplate{
		Source: confT,
		ExportPath: "config.minikube.yaml",
		Data: map[string]string{
			"projectName": r.Project.GetName(),
		},
	}

	_= u.GoRunTemplate(&goTemplate)

	color.Green("Generated")


}