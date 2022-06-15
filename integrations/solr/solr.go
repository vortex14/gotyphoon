package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/integrations/docker"
	"github.com/vortex14/gotyphoon/utils"
	"strings"

	JQ "github.com/itchyny/gojq"
	SolrSF9V "github.com/stevenferrer/solr-go"
	SolrVan "github.com/vanng822/go-solr/solr"

	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Document = SolrSF9V.M

type DockerOptions struct {
	MatchName string
}

type ConnectOptions struct {
	RemoteSSHConnectionURL string
	Collection             string
	Endpoint               string
}

type SchemaField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SchemaFields []SchemaField

const (
	Integer = "pint"
	String  = "string"
	Float   = "pdouble"
	MFloat  = "pdoubles"

	SchemaFieldsJQPath = ".schema.fields"
)

type Client struct {
	singleton.Singleton
	Options  *ConnectOptions
	DOptions *DockerOptions

	clientSF9V SolrSF9V.Client
	clientVan  *SolrVan.SolrInterface

	LOG interfaces.LoggerInterface
}

func (c *Client) ConnectToDockerNode() {
	d := &docker.Docker{RemoteSSHUrl: c.Options.RemoteSSHConnectionURL}
	err, containers := d.GetRemoteActiveContainersList()
	if err != nil {
		color.Red("%+v", err)
		return
	}
	//var matchedContainer types.Container
	for _, container := range containers {
		if strings.Contains(container.Image, c.DOptions.MatchName) {
			//matchedContainer = container
			break
		}
	}

	//d.RunRemoteCommandInContainer(matchedContainer)

}

func (c *Client) Init() *Client {
	c.Construct(func() {
		c.LOG = log.New(log.D{"connection": "solr"})
		if c.Options == nil {
			c.LOG.Error(Errors.SolrConnectionsOptionsNotFound.Error())
			return
		}
		solrUrl := fmt.Sprintf("http://%s", c.Options.Endpoint)

		solrDetailUrl := fmt.Sprintf("%s/solr", solrUrl)
		si, err := SolrVan.NewSolrInterface(solrDetailUrl, c.Options.Collection)
		if err != nil {
			c.LOG.Error(Errors.SolrConnectionEndpointError.Error())
			return
		}
		c.clientVan = si
		//baseURL := strings.ReplaceAll(c.Options.Endpoint, "/solr", "/")
		solrClient := SolrSF9V.NewJSONClient(solrUrl)
		c.clientSF9V = solrClient
		c.LOG.Debug("connected!")
	})
	return c
}

func (c *Client) GetClientVan() *SolrVan.SolrInterface {
	return c.clientVan
}

func (c *Client) GetClientSF9V() SolrSF9V.Client {
	return c.clientSF9V
}

func (c *Client) GetFullSchema() (error, *SolrVan.SolrResponse) {
	schema, err := c.clientVan.Schema()
	if err != nil {
		return err, nil
	}
	dump, errD := schema.All()
	if errD != nil {
		return errD, nil
	}
	return nil, dump
}

func (c *Client) AddSchemaFields(fields []SolrSF9V.Field) error {
	err := c.clientSF9V.AddFields(context.Background(), c.Options.Collection, fields...)

	return err
}

func (c *Client) UpdateDoc(context context.Context, doc SolrSF9V.M) (error, *SolrSF9V.UpdateResponse) {
	err, s := c.UpdateDocs(context, []SolrSF9V.M{doc})
	return err, s
}

func (c *Client) RemoveAllDocs() (error, *SolrVan.SolrUpdateResponse) {
	all, err := c.clientVan.DeleteAll()
	return err, all
}

func (c *Client) UpdateDocs(context context.Context, docs []SolrSF9V.M) (error, *SolrSF9V.UpdateResponse) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).
		Encode(docs)
	update, err := c.clientSF9V.Update(context, c.Options.Collection, SolrSF9V.JSON, buf)
	return err, update
}

func (c *Client) Commit(context context.Context) error {
	err := c.clientSF9V.Commit(context, c.Options.Collection)
	return err
}

func (c *Client) UpdateDocsAndCommit(context context.Context, docs []SolrSF9V.M) (error, *SolrSF9V.UpdateResponse) {
	err, response := c.UpdateDocs(context, docs)
	if err != nil {
		return err, nil
	}

	errC := c.Commit(context)
	if errC != nil {
		return errC, nil
	}

	return nil, response
}

func (c *Client) UpdateDocAndCommit(context context.Context, doc SolrSF9V.M) (error, *SolrSF9V.UpdateResponse) {
	err, response := c.UpdateDoc(context, doc)
	if err != nil {
		return err, nil
	}

	errC := c.Commit(context)
	if errC != nil {
		return errC, nil
	}

	return nil, response
}

func (c *Client) RemoveSchema(fields []SolrSF9V.Field) error {
	// ignore early created fields for creating schema else will be error from solr rest
	var mFields []SolrSF9V.Field
	for _, originField := range fields {
		mFields = append(mFields, SolrSF9V.Field{Name: originField.Name})
	}
	return c.clientSF9V.DeleteFields(context.Background(), c.Options.Collection, mFields...)
}

func (c *Client) GetCollectionSchemaFields() (error, []SolrSF9V.Field) {
	err, schema := c.GetFullSchema()
	if err != nil {
		return err, nil
	}

	query, err := JQ.Parse(SchemaFieldsJQPath)
	if err != nil {
		return err, nil
	}

	var collectionSchemaFields []SolrSF9V.Field
	fields, ok := query.Run(schema.Response).Next()
	if ok {
		for _, field := range fields.([]interface{}) {
			fieldMap := field.(map[string]interface{})
			fieldName := fieldMap["name"].(string)
			fieldType := fieldMap["type"].(string)
			if utils.IsStrContain(fieldName, "id", "_version_") {
				continue
			}

			collectionSchemaFields = append(collectionSchemaFields,
				SolrSF9V.Field{
					Name: fieldName,
					Type: fieldType,
				})
		}
	} else {
		err = Errors.JqExecuteQueryError
	}

	return err, collectionSchemaFields
}

func (c *Client) GetInfoCore(coreName string) (error, *SolrVan.SolrResponse) {
	admin, err := c.clientVan.CoreAdmin()
	if err != nil {
		return err, nil
	}

	info, err := admin.Status(coreName)
	if err != nil {
		return err, nil
	}

	return nil, info
}

func (c *Client) RenameCore(lastName string, newName string) (error, *SolrVan.SolrResponse) {
	admin, err := c.clientVan.CoreAdmin()
	if err != nil {
		return err, nil
	}
	response, err := admin.Rename(lastName, newName)
	if err != nil {
		return err, nil
	}

	return nil, response
}

func (c *Client) CreateNewCore(name string) error {
	return c.clientSF9V.CreateCore(context.Background(), SolrSF9V.NewCreateCoreParams(name))
}

func (c *Client) CreateCoreConfig(name string) {
	//c.clientSF9V.CreateCore()
}

func (c *Client) CreateCollection() {
	params := SolrSF9V.NewCollectionParams()
	params.Name("test-new-collections")
	//err := c.clientSF9V.CreateCollection(context.Background(), params)
	c.clientVan.SetCore("test-new-collections")
	admin, err := c.clientVan.CoreAdmin()
	if err != nil {
		return
	}
	err2 := c.clientSF9V.CreateCore(context.Background(), SolrSF9V.NewCreateCoreParams("test-test-test-core"))
	if err2 != nil {
		color.Red("%+v", err2)
		return
	}

	admin.SetBasicAuth("tst", "test2")

	//_, s := utils.DumpPrettyJson(reload)
	//println(s, admin)
	//if err != nil {
	//	color.Red("%s", err.Error())
	//	return
	//}
}
