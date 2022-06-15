package helm

import (
	"embed"
	"io/fs"
	"os"

	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type Resources struct {
	Project interfaces.Project
}

//go:embed templates
var helmTemplates embed.FS

func (r *Resources) RemoveHelmMinikubeManifests() {
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

func (r *Resources) BuildHelmMinikubeResources() {
	color.Yellow("build helm minikube resources ...")
	r.Project.LoadConfig()

	u := utils.Utils{}

	dir, err := fs.Sub(helmTemplates, "templates/v1.1")
	if err != nil {
		panic(err)
	}

	fileObject := &interfaces.FileObject{
		Path: ".",
	}

	err = u.CopyDirAndReplaceLabel("helm", &interfaces.ReplaceLabel{Label: "{{PROJECT_NAME}}", Value: r.Project.GetName()}, fileObject, dir)

	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	_, dataTDeployLocal := u.GetGoTemplate(&interfaces.FileObject{Path: ".", Name: "helm_deploy.gosh"}, dir)

	dataConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDeployLocal := interfaces.GoTemplate{
		Source:     dataTDeployLocal,
		ExportPath: "helm_deploy.sh",
		Data:       dataConfig,
	}

	u.GoRunTemplate(&goTemplateHelmDeployLocal)

	_, dataTDumpLocal := u.GetGoTemplate(&interfaces.FileObject{Path: ".", Name: "helm_dump.gosh"}, dir)

	dataDumpConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDumpLocal := interfaces.GoTemplate{
		Source:     dataTDumpLocal,
		ExportPath: "helm_dump.sh",
		Data:       dataDumpConfig,
	}

	u.GoRunTemplate(&goTemplateHelmDumpLocal)

	_, dataTDeleteLocal := u.GetGoTemplate(&interfaces.FileObject{Path: ".", Name: "helm_delete.gosh"}, dir)

	dataDeleteConfig := map[string]string{
		"projectName": r.Project.GetName(),
	}

	goTemplateHelmDeleteLocal := interfaces.GoTemplate{
		Source:     dataTDeleteLocal,
		ExportPath: "helm_delete.sh",
		Data:       dataDeleteConfig,
	}

	u.GoRunTemplate(&goTemplateHelmDeleteLocal)

	if err := os.Chmod("helm_delete.sh", 0755); err != nil {
		color.Red("%s", err)
	}

	if err := os.Chmod("helm_deploy.sh", 0755); err != nil {
		color.Red("%s", err)
	}

	if err := os.Chmod("helm_dump.sh", 0755); err != nil {
		color.Red("%s", err)
	}

	_, confT := u.GetGoTemplate(&interfaces.FileObject{
		Path: ".",
		Name: "config.minikube.goyaml",
	}, dir)
	goTemplate := interfaces.GoTemplate{
		Source:     confT,
		ExportPath: "config.minikube.yaml",
		Data: map[string]string{
			"projectName": r.Project.GetName(),
		},
	}

	_ = u.GoRunTemplate(&goTemplate)

	color.Green("Generated")

}
