package swagger

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"reflect"
	"strings"
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

func (oa *OpenApi) AddOperation(path string, method string, operation *openapi3.Operation) {
	oa.swagger.AddOperation(path, method, operation)
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

func (oa *OpenApi) MoveRequiredFieldsToTopLevel() {
	for _, schemaRef := range oa.swagger.Components.Schemas {
		var required []string
		for _, val := range schemaRef.Value.Properties {
			if len(val.Value.Required) > 0 {
				required = append(required, val.Value.Required...)
				val.Value.Required = nil
			}

		}
		schemaRef.Value.Required = required
	}
}

func (oa *OpenApi) CreateBaseSchemasFromStructure(source interface{}) *openapi3.SchemaRef {
	customizer := openapi3gen.SchemaCustomizer(
		func(name string, ft reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {

			schema.Title = ft.Name()

			if len(tag.Get("description")) > 0 {
				schema.Description = tag.Get("description")
			}

			if tag.Get("binding") == "required" {
				schema.Required = append(schema.Required, name)
			}

			if strings.Contains(ft.String(), ".") {
				if utils.IsFirstUpLetter(ft.Name()) && !oa.IsExistsSchema(ft.Name()) {
					oa.AddComponent(ComponentTypeSchema, ft.Name(), schema.NewRef())
				}

				for key, val := range schema.Properties {

					if utils.IsFirstUpLetter(val.Ref) {
						if !oa.IsExistsSchema(val.Ref) {
							oa.AddComponent(ComponentTypeSchema, val.Ref, schema.NewRef())
						} else {
							val.Ref = fmt.Sprintf("#/components/schemas/%s", val.Ref)
						}

					} else {
						val.Ref = ""
						val.Value.Title = key
					}

					if val.Value.Items != nil {
						val.Value.Items.Ref = fmt.Sprintf("#/components/schemas/%s", val.Value.Items.Ref)
					}
				}
			}

			return nil
		})

	generator := openapi3gen.NewGenerator(customizer)

	sc, err := generator.GenerateSchemaRef(reflect.TypeOf(source))
	if err != nil {
		panic(err)
	}

	oa.MoveRequiredFieldsToTopLevel()

	sc.Ref = fmt.Sprintf("#/components/schemas/%s", sc.Ref)

	println(fmt.Sprintf("%+v", sc))

	return sc
}

func (oa *OpenApi) AddSwaggerOperation(
	resource interfaces.ResourceInterface,
	action interfaces.ActionInterface,
	method, path string) *openapi3.Operation {

	operation := &openapi3.Operation{
		Tags:        action.GetTags(),
		Summary:     action.GetSummary(),
		Responses:   openapi3.Responses{},
		Description: action.GetDescription(),
		OperationID: strings.ToLower(
			fmt.Sprintf("Handle_%s_%s_%s",
				strings.ReplaceAll(resource.GetPath(), "/", "_"),
				action.GetName(),
				method)),
	}

	requestModel := action.GetRequestModel()
	if requestModel != nil {

		requestSchema := oa.CreateBaseSchemasFromStructure(requestModel)

		operation.RequestBody = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: action.IsRequiredRequestModel(),
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: requestSchema,
					},
				},
			},
		}
	}

	params := action.GetParams()

	if params != nil {

		schemaRef := CreateRefSchemaFromStruct(params)

		for name, ref := range schemaRef.Value.Properties {

			status := false
			for i := 0; i < reflect.TypeOf(params).Elem().NumField(); i++ {
				field := reflect.TypeOf(params).Elem().Field(i)
				ref.Value.Title = name
				if strings.ToLower(field.Name) != name {
					continue
				}
				if field.Tag.Get("binding") == "required" {
					status = true
				}
			}

			operation.AddParameter(&openapi3.Parameter{Required: status, Schema: ref, In: "query", Name: name})

			ref.Value.Title = name
		}

	}

	oa.AddOperation(path, method, operation)

	return operation
}

func (oa *OpenApi) AddSwaggerResponse(

	title *string,
	code int,
	action interfaces.ActionInterface,
	operation *openapi3.Operation,

	response interface{}) {

	_schema := oa.CreateBaseSchemasFromStructure(response)

	swaggerResponse := &openapi3.Response{
		Description: title,
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: _schema,
			},
		},
	}
	operation.AddResponse(code, swaggerResponse)

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
			OpenAPI: "3.0.1",
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
