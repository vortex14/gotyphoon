package swagger

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"os"
)

const (
	ComponentTypeSecurity  = "security"
	ComponentTypeParameter = "parameter"
	ComponentTypeSchema    = "schema"
)

type OpenApi struct {
	swagger      openapi3.T
	cfg          Config
	securityReqs openapi3.SecurityRequirements
}

func (oa *OpenApi) AddComponent(componentType string, name string, ref interface{}) {
	switch componentType {
	case ComponentTypeSchema:
		oa.swagger.Components.Schemas[name] = ref.(*openapi3.SchemaRef)
	case ComponentTypeParameter:
		oa.swagger.Components.Parameters[name] = ref.(*openapi3.ParameterRef)
	case ComponentTypeSecurity:
		oa.swagger.Components.SecuritySchemes[name] = ref.(*openapi3.SecuritySchemeRef)
	}
}

func (oa *OpenApi) GetDump() []byte {
	contract, _ := oa.swagger.MarshalJSON()
	return contract
}

func (oa *OpenApi) IsExistsSchema(name string) bool {
	status := false
	if _, ok := oa.swagger.Components.Schemas[name]; ok {
		status = true
	}
	return status
}

func (oa *OpenApi) addSecurity() {
	//for name, security := range oa.cfg.Server.OpenApi.Security {
	//	oa.AddComponent(ComponentTypeSecurity, name, &openapi3.SecuritySchemeRef{
	//		Value: &openapi3.SecurityScheme{
	//			Type: security.Type,
	//			Name: security.Name,
	//			In:   security.In,
	//		},
	//	})
	//}

	oa.securityReqs = openapi3.SecurityRequirements{}

	securityReq := openapi3.SecurityRequirement{}
	for securityName := range oa.swagger.Components.SecuritySchemes {
		securityReq[securityName] = make([]string, 0)
	}
	oa.securityReqs = append(oa.securityReqs, securityReq)
}

func (oa *OpenApi) Generate(openApiFile string) error {

	oa.addSecurity()

	oa.addModels()

	//oa.addRoutes()

	var data []byte
	var err error
	if data, err = oa.swagger.MarshalJSON(); err != nil {
		return err
	}
	return os.WriteFile(oa.cfg.Server.OpenApi.Json.Local, data, 0644)
}

func (oa *OpenApi) addModels() {

	// Build schemas
	//for modelName, model := range oa.cfg.Server.Models {
	//
	//	properties := make(map[string]*openapi3.SchemaRef)
	//	required := make([]string, 0)
	//	for _, Prop := range model.Properties {
	//
	//		schemaRef, isOptional := oa.buildParam(Prop)
	//
	//		properties[Prop.Name] = schemaRef
	//
	//		if !isOptional {
	//			required = append(required, Prop.Name)
	//		}
	//	}
	//
	//	oa.AddComponent(ComponentTypeSchema, modelName, &openapi3.SchemaRef{
	//		Value: &openapi3.Schema{
	//			Properties:  properties,
	//			Required:    required,
	//			Description: model.Description,
	//		},
	//	})
	//}
}

func ConstructorNewFromArgs(title, description, version string, host []string) *OpenApi {
	return &OpenApi{
		swagger: openapi3.T{
			OpenAPI: version,
			Info: &openapi3.Info{
				Title:       title,
				Description: description,
				Version:     version,
			},
			Servers: openapi3.Servers{
				&openapi3.Server{
					URL: fmt.Sprintf("%s://%s/", host[0], host[1]),
				},
			},
			Components: &openapi3.Components{
				Schemas:         make(map[string]*openapi3.SchemaRef),
				SecuritySchemes: make(map[string]*openapi3.SecuritySchemeRef),
				Parameters:      make(map[string]*openapi3.ParameterRef),
			},
		},
	}
}
