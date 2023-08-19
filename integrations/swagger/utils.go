package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/vortex14/gotyphoon/utils"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

func DumpStructSchema(someIns interface{}) []byte {
	schemaRef, err := CreateSchemaFromStruct(someIns)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(schemaRef, "", "  ")
	if err != nil {
		panic(err)
	}

	return data
}

func CreateSchemaFromStruct(someIns interface{}) (*openapi3.SchemaRef, error) {
	schemaRef, err := openapi3gen.NewSchemaRefForValue(someIns, nil)
	if err != nil {
		return nil, err
	}

	return schemaRef, nil
}

func MoveRequiredFieldsToTopLevel(swagger *openapi3.T) {
	for _, schemaRef := range swagger.Components.Schemas {
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

func CreateBaseSchemasFromStructure(tmpl *OpenApi, source interface{}) {
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
				if utils.IsFirstUpLetter(ft.Name()) && !tmpl.IsExistsSchema(ft.Name()) {
					tmpl.AddComponent(ComponentTypeSchema, ft.Name(), schema.NewRef())
				}

				for key, val := range schema.Properties {

					if utils.IsFirstUpLetter(val.Ref) {
						if !tmpl.IsExistsSchema(val.Ref) {
							tmpl.AddComponent(ComponentTypeSchema, val.Ref, schema.NewRef())
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

	_, err := generator.GenerateSchemaRef(reflect.TypeOf(source))
	if err != nil {
		panic(err)
	}
}

func CreateRefSchemaFromStruct(instance interface{}) *openapi3.SchemaRef {
	_schemas := make(openapi3.Schemas)

	schParam, _ := openapi3gen.NewSchemaRefForValue(instance, _schemas)

	return schParam
}

func CreateFileSchema() *openapi3.Schema {
	ref := openapi3.NewSchema().WithFormat("binary")
	ref.Type = "string"
	return ref
}
