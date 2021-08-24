package interfaces


type HelmInterface interface {
	BuildHelmMinikubeResources()
	RemoveHelmMinikubeManifests()
}

