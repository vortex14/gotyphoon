package swagger

import (
	"encoding/json"
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
		required := []string{}
		for _, val := range schemaRef.Value.Properties {
			if len(val.Value.Required) > 0 {
				required = append(required, val.Value.Required...)
				val.Value.Required = nil
			}

		}

		schemaRef.Value.Required = required
	}
}
