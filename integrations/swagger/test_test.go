package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/utils"
	"reflect"
	"strings"
	"testing"
	"time"
)

type (
	SomeStruct struct {
		Bool    bool                      `json:"bool"`
		Int     int                       `json:"int"`
		Int64   int64                     `json:"int64"`
		Float64 float64                   `json:"float64"`
		String  string                    `json:"string"`
		Bytes   []byte                    `json:"bytes"`
		JSON    json.RawMessage           `json:"json"`
		Time    time.Time                 `json:"time"`
		Slice   []SomeOtherType           `json:"slice"`
		Map     map[string]*SomeOtherType `json:"map"`

		Struct struct {
			X string `json:"x"`
		} `json:"struct"`

		EmptyStruct struct {
			Y string
		} `json:"structWithoutFields"`

		Ptr *SomeOtherType `json:"ptr"`
	}

	SomeOtherType string
)

type Role struct {
	IsAdmin bool `json:"isAdmin"`
}

type Link struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}

type Img struct {
	Href string `json:"href"`
	Link *Link  `json:"link"`
}

type User struct {
	Name   string `json:"name" description:"test descr" required:"!"`
	Role   *Role  `json:"role"`
	Images []*Img `json:"images" description:"photo for profile"`
	Links  []*Link
}

const ResultGen = `{
  "properties": {
    "bool": {
      "type": "boolean"
    },
    "bytes": {
      "format": "byte",
      "type": "string"
    },
    "float64": {
      "format": "double",
      "type": "number"
    },
    "int": {
      "type": "integer"
    },
    "int64": {
      "format": "int64",
      "type": "integer"
    },
    "json": {},
    "map": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "ptr": {
      "type": "string"
    },
    "slice": {
      "items": {
        "type": "string"
      },
      "type": "array"
    },
    "string": {
      "type": "string"
    },
    "struct": {
      "properties": {
        "x": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "structWithoutFields": {},
    "time": {
      "format": "date-time",
      "type": "string"
    }
  },
  "type": "object"
}`

func TestStructGen(t *testing.T) {
	Convey("test gen swagger model from go struct", t, func() {

		_struct := &SomeStruct{}

		So(fmt.Sprintf("%s", DumpStructSchema(_struct)), ShouldEqual, ResultGen)

	})

}

func TestBaseOpenAPITemplate(t *testing.T) {

	Convey("test generate swagger contract", t, func() {
		tmpl := ConstructorNewFromArgs("demo v1.1", "test description", "3.1.0", []string{"https", "localhost"})

		contract, err := tmpl.swagger.MarshalJSON()

		So(err, ShouldBeNil)
		println(fmt.Sprintf("%s", contract))
		So(
			fmt.Sprintf("%s", contract),
			ShouldEqual,
			utils.ClearStrTabAndN(`
				{"components":{},
				"info":{"description":"test description","title":"demo v1.1","version":"3.1.0"},
				"openapi":"3.1.0",
				"paths":null,
				"servers":[{"url":"https://localhost/"}]}`),
		)
	})

}

func TestAddNewOperation(t *testing.T) {

	Convey("create a first operation", t, func() {

		tmpl := ConstructorNewFromArgs(
			"demo v1.1",
			"test description",
			"3.1.0",
			[]string{"https", "localhost"})

		operation := &openapi3.Operation{
			Tags:        []string{"test tag"},
			Summary:     "short description",
			Responses:   openapi3.Responses{},
			Description: "some description",
		}

		tmpl.swagger.AddOperation("/", "GET", operation)

		contract, err := tmpl.swagger.MarshalJSON()

		So(err, ShouldBeNil)
		println(fmt.Sprintf("%s", contract))

		So(
			fmt.Sprintf("%s", contract),
			ShouldEqual,
			utils.ClearStrTabAndN(`{"components":{},
				"info":{"description":"test description","title":"demo v1.1","version":"3.1.0"},
				"openapi":"3.1.0",
				"paths":{"/":{"get":{"description":"some description","responses":{},
				"summary":"short description","
				tags":["test tag"]}}},
				"servers":[{"url":"https://localhost/"}]}`))

	})

}

func TestAddResponseForOperation(t *testing.T) {
	Convey("add a new response for any operation", t, func() {
		operation := &openapi3.Operation{
			Tags:        []string{"test tag"},
			Summary:     "short description",
			Responses:   openapi3.Responses{},
			Description: "some description",
		}

		type SuccessResponse struct {
			Message string `json:"message"`
			Status  bool   `json:"status"`
		}

		type ErrResponse struct {
			Message string `json:"message"`
			Status  bool   `json:"status"`
			Code    int    `json:"code"`
		}

		responseSuccessDescription := "Success response"

		errResponseDescription := "error response"

		_schema, _ := CreateSchemaFromStruct(&SuccessResponse{})

		successResponse := &openapi3.Response{
			Description: &responseSuccessDescription,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: _schema,
				},
			},
		}

		_schemaErr, _ := CreateSchemaFromStruct(&ErrResponse{})

		errorResponse := &openapi3.Response{
			Description: &errResponseDescription,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: _schemaErr,
				},
			},
		}

		operation.AddResponse(200, successResponse)
		operation.AddResponse(400, errorResponse)

		data, err := operation.MarshalJSON()

		So(err, ShouldBeNil)
		println(fmt.Sprintf("%s", data))
		So(fmt.Sprintf("%s", data), ShouldEqual, utils.ClearStrTabAndN(`
			{"description":"some description",
			"responses":{"200":{"content":{"application/json":{"schema":{"properties":
			{"message":{"type":"string"},"status":{"type":"boolean"}},"type":"object"}}},
			"description":"Success response"},"400":{"content":{"application/json":
			{"schema":{"properties":{"code":{"type":"integer"},"message":{"type":"string"},
			"status":{"type":"boolean"}},"type":"object"}}},"description":"error response"}},
			"summary":"short description","tags":["test tag"]}`))
	})
}

func TestCustomizeRequiredFields(t *testing.T) {

	Convey("test customize schema", t, func() {
		type testRequest struct {
			Payload interface{} `json:"payload" required:"!"`
			Data    string      `json:"data"`
		}

		customizer := openapi3gen.SchemaCustomizer(
			func(name string, ft reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {

				if tag.Get("required") == "!" {
					schema.Required = append(schema.Required, name)
				}
				return nil
			})

		schemaRef, err := openapi3gen.NewSchemaRefForValue(
			&testRequest{}, nil, openapi3gen.UseAllExportedFields(), customizer)

		if err != nil {
			panic(err)
		}

		var data []byte
		if data, err = json.Marshal(schemaRef); err != nil {
			panic(err)
		}

		println(fmt.Sprintf("%s", data))

		testR := `{"properties":{"data":{"type":"string"},"payload":{"required":["payload"]}},"type":"object"}`

		So(fmt.Sprintf("%s", data), ShouldEqual, testR)

	})

}

func TestDefinition(t *testing.T) {

	schemas := make(openapi3.Schemas)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(&User{}, schemas)

	if len(schemas) != 0 {
		panic(`No references should have been collected at this point`)
	}

	if schemaRef, err = openapi3gen.NewSchemaRefForValue(&User{}, schemas); err != nil {
		panic(err)
	}
	//schemaRef.Value.Items.

	println(">>>>>>", schemaRef.Ref)

	var data []byte
	if data, err = json.MarshalIndent(schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
	if data, err = json.MarshalIndent(schemas, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemas: %s\n", data)

	//sch, _ := CreateSchemaFromStruct(&User{})
	//d, _ := sch.MarshalJSON()
	//
	//println(fmt.Sprintf("%s", d))

}

func TestCyclicStructures(t *testing.T) {

	type TestS struct {
		User *TestS `json:"user"`
	}

	schemas := make(openapi3.Schemas)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(&TestS{}, schemas, openapi3gen.ThrowErrorOnCycle())
	if schemaRef != nil || err == nil {
		println(err.Error())
		panic(`With option ThrowErrorOnCycle, an error is returned when a schema reference cycle is found`)
	}
	if _, ok := err.(*openapi3gen.CycleError); !ok {
		panic(`With option ThrowErrorOnCycle, an error of type CycleError is returned`)
	}
	if len(schemas) != 0 {
		panic(`No references should have been collected at this point`)
	}

	if schemaRef, err = openapi3gen.NewSchemaRefForValue(&TestS{}, schemas); err != nil {
		panic(err)
	}

	var data []byte
	if data, err = json.MarshalIndent(schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
	if data, err = json.MarshalIndent(schemas, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemas: %s\n", data)

}

func TestGenerateSchemaRefAndSchema(t *testing.T) {
	//propSchemaRef := &openapi3.SchemaRef{}

	schemas := make(openapi3.Schemas)
	schemaRef, _ := openapi3gen.NewSchemaRefForValue(&User{}, schemas, openapi3gen.ThrowErrorOnCycle())

	propSchema := &openapi3.Schema{}

	propSchema.Items = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "TypeStructure",
		},
	}
	schemaRef.Value = propSchema

	b, _ := propSchema.MarshalJSON()

	println(fmt.Sprintf("%s", b))

}

func TestCreateSwagWithDefinitionsAndRefs(t *testing.T) {

	Convey("test create swagger json with definitions and refs", t, func() {

		tmpl := ConstructorNewFromArgs(
			"demo v1.1",
			"test description",
			"3.1.0",
			[]string{"https", "localhost"})

		//operation := &openapi3.Operation{
		//	Tags:        []string{"test tag"},
		//	Summary:     "short description",
		//	Responses:   openapi3.Responses{},
		//	Description: "some description",
		//}

		schemas := make(openapi3.Schemas)

		//schemaRef.Ref = fmt.Sprintf("#/components/schemas/%s", "TestSchema")

		customizer := openapi3gen.SchemaCustomizer(
			func(name string, ft reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {

				if strings.Contains(ft.String(), "swagger.") && name != "_root" {
					println(name, fmt.Sprintf(" ----- %s, :::: %s --- %s", ft, ft.String(), ft.Kind().String()))

					schemaRef, _ := openapi3gen.NewSchemaRefForValue(reflect.New(ft.Elem()), schemas)
					d, _ := schemaRef.MarshalJSON()
					println(fmt.Sprintf("%s", d))
				}

				//fmt.Println(reflect.TypeOf(ft))
				//fmt.Println(reflect.ValueOf(ft).Kind())

				return nil
			})

		schemaRef, _ := openapi3gen.NewSchemaRefForValue(&User{}, schemas, customizer)

		schemaRef.Value.Required = []string{"name"}

		tmpl.AddComponent(ComponentTypeSchema, "TestSchema", schemaRef)

		//_ = openapi3.NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", "TestSchema"), schemas)
		//ns := &openapi3.NewSchemaRef()
		//operation.RequestBody.Value =
		//	tmpl.swagger.AddOperation("/", "GET", operation)

		//contract, err := tmpl.swagger.MarshalJSON()

		//So(err, ShouldBeNil)

		//body := (&openapi3.RequestBody{}).
		//	WithDescription("test description").WithRequired(true).WithSchemaRef(schemaRef, []string{})

		d, _ := tmpl.swagger.MarshalJSON()

		println(fmt.Sprintf("%s, %s", d, schemaRef.Ref))
	})

}

func TestRecursiveType(t *testing.T) {
	type RecursiveType struct {
		Field1     string           `json:"field1"`
		Field2     string           `json:"field2"`
		Field3     string           `json:"field3"`
		Components []*RecursiveType `json:"children,omitempty"`
	}

	schemas := make(openapi3.Schemas)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(&RecursiveType{}, schemas)
	if err != nil {
		panic(err)
	}

	var data []byte
	if data, err = json.MarshalIndent(&schemas, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemas: %s\n", data)
	if data, err = json.MarshalIndent(&schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
}

func TestCreateContractWithComponentsDefs(t *testing.T) {

	Convey("create required schemas", t, func() {

		result := `{"components":{"schemas":{
						"Img":{"description":"photo for profile",
						"properties":{"href":{"title":"href","type":"string"},
						"link":{"$ref":"#/components/schemas/Link"}},
						"title":"Img","type":"object"},
						"Link":{"properties":{"source":{"title":"source","type":"string"},
						"title":{"title":"title","type":"string"}},"title":"Link","type":"object"},
						"Role":{"properties":{"isAdmin":{"title":"isAdmin","type":"boolean"}},
						"title":"Role","type":"object"},
						"User":{"properties":{"images":{"description":"photo for profile",
						"items":{"$ref":"#/components/schemas/Img"},"title":"images","type":"array"},
						"name":{"description":"test descr","title":"name","type":"string"},
						"role":{"$ref":"#/components/schemas/Role"}},"required":["name"],
						"title":"User","type":"object"}}},
						"info":{"description":"test description","title":"demo v1.1","version":"3.0.1"},
						"openapi":"3.0.1","paths":null,"servers":[{"url":"https://localhost/"}]}
`
		tmpl := ConstructorNewFromArgs(
			"demo v1.1",
			"test description",
			"3.0.1",
			[]string{"https", "localhost"})

		CreateBaseSchemasFromStructure(tmpl, &User{})

		MoveRequiredFieldsToTopLevel(&tmpl.swagger)

		contract := tmpl.GetDump()

		println(fmt.Sprintf("%s", contract))
		So(fmt.Sprintf("%s", contract), ShouldEqual, utils.ClearStrTabAndN(result))

	})

}

func TestGenerateLinkType(t *testing.T) {

	schemas := make(openapi3.Schemas)
	schRef, _ := openapi3gen.NewSchemaRefForValue(&User{}, schemas)
	b, _ := schRef.MarshalJSON()
	println(fmt.Sprintf("%s", b))
}
