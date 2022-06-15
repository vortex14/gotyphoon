# Package typhoon/integrations/solr/main

## Imports

    import (
        "context"
        "github.com/fatih/color"
        Solr "github.com/sf9v/solr-go"
        "github.com/vortex14/gotyphoon/extensions/data/fake"
        "github.com/vortex14/gotyphoon/integrations/solr"
        "github.com/vortex14/gotyphoon/log"
        "github.com/vortex14/gotyphoon/utils"
        "path/filepath"
    )

# Structs

## ParseConfig

    type ParseConfig struct {
        Test          bool
        Internal      bool
        UnderScore    bool
    }

ParseConfig to use when invoking ParseAny, ParseSingleFileWalker, and ParseSinglePackageWalker.

**ParserConfig**

    package main
    
    import (
            "context"
            "path/filepath"
    
            "github.com/fatih/color"
            Solr "github.com/sf9v/solr-go"
            "github.com/vortex14/gotyphoon/extensions/data/fake"
            "github.com/vortex14/gotyphoon/integrations/solr"
            "github.com/vortex14/gotyphoon/utils"
    
            "github.com/vortex14/gotyphoon/log"
    )
    
    func init()  {
            log.InitD()
    }
    
    func CreateSchema(client *solr.Client,schema []Solr.Field, collection string)  {
            errI := client.AddSchemaFields(collection, schema)
            color.Red(">>>>>> added: %+v", errI)
            return
    
    }


    type CollectionSchema struct {
            Collection string
            Schema []Solr.Field
    }
    
    func GetFieldsFromTemplate(template *solr.TemplateConfig) chan CollectionSchema {
            out := make(chan CollectionSchema)
            go func(ch chan CollectionSchema) {
                    for _, collection := range template.Collections {
                            var schema []Solr.Field
                            for _, field := range collection.SchemaFields {
                                    schema = append(schema, Solr.Field{
                                            Name:                 field.Name,
                                            Type:                 field.Type,
                                            Stored:               field.Stored,
                                            Indexed:              field.Indexed,
                                            Required:             field.Required,
                                            DocValues:            field.DocValues,
                                            MultiValued:          field.MultiValued,
                                            UseDocValuesAsStored: field.UseDocValuesAsStored,
                                    })
                            }; ch <- CollectionSchema{Schema: schema, Collection: collection.Name}
                    }
                    close(ch)
            }(out)
            return out
    }


    func main()  {
            options := &solr.ConnectOptions{
                    Collection: "5520.products",
                    Endpoint:   "solr.typhoon-s1.ru",
                    RemoteSSHConnectionURL: "ssh://root@195.201.108.45:22",
                    //Endpoint: "localhost:8983",
            }
            client := (
                    &solr.Client{Options: options, DOptions: &solr.DockerOptions{
                            MatchImageName: "bitnami/solr:8-debian-10",
                            //MatchContainerName: "solr-typhoon-instance",
                    }}).
                    Init()
    
            pathT := filepath.Join(utils.GetCurrentDir(), "templates", "collection.yaml")
            data := utils.ReadFile(pathT)
            //var collections solr.TemplateConfig
    
            //println(data)
            var o solr.TemplateConfig
    
            err := utils.YamlLoad(&o, []byte(data))
            if err != nil {
                    color.Red(err.Error())
                    return
            }
    
            container := client.ConnectToDockerNode()
            ping, err := container.Client.Ping(context.Background())
            if err != nil {
                    color.Red("%s", err.Error())
                    return
            }
    
            color.Yellow("%+v", ping)
            err, e := container.Exec(context.Background(), []string{"ls"})
            if err != nil {
                    color.Red("%s", err.Error())
                    return
            }
            println(e.Stdout())
            //for collectionData := range GetFieldsFromTemplate(&o) {
            //      color.Red("%+v", len(collectionData.Schema))
            //      //_, _ = container.Exec(context.Background(), []string{"/opt/bitnami/solr/bin/solr", "delete", "-c", collectionData.Collection})
            //      //_, _ = container.Exec(context.Background(), []string{"/opt/bitnami/solr/bin/solr", "create", "-c", collectionData.Collection})
            //      //CreateSchema(client, collectionData.Schema, collectionData.Collection)
            //}




            return
    
            //_, _ = client.GetCollectionSchemaFields()
            //container := client.ConnectToDockerNode()
    
            //err, e := container.Exec(context.Background(), []string{"/opt/bitnami/solr/bin/solr", "delete", "-c", "5520.products"})
            //err, e := container.Exec(context.Background(), []string{"/opt/bitnami/solr/bin/solr", "create", "-c", "5520.products"})
            //err, e := container.Exec(context.Background(), []string{"ls", "/opt/bitnami/solr/data/"})
            //if err != nil {
            //      return
            //}
    
            //CreateSchema(client)
            //println(e.Stdout())
    
            //return
    
            //return
            //_, res := utils.PrintPrettyJson(fields)
    
            //err := client.RemoveSchema(fields)
            //if err != nil {
            //      color.Red("%s", err.Error())
            //      return
            //}
            //color.Yellow("%s", res)


            errR, s2 := client.RemoveAllDocs()
            if errR != nil {
                    color.Red("%+v", errR)
                    return
            }


            _, s2P := utils.PrintPrettyJson(s2)
            color.Red("%s", s2P)
            //return
    
            product := fake.CreateProductWithId()
    
            product.Id = "947cb24a-53c0-4078-bbcd-3e4262dc8ae0"


            product.Price.ListingPrice = 10.5
    
            err, s := client.UpdateDocAndCommit(context.Background(), solr.Document{
                    "id":    product.Id,
                    "url":   product.Url,
                    "title": product.Title,
                    "description": product.Description,
                    "upc": product.Upc,
                    "listing_price": product.Price.ListingPrice,
                    "offer_price": product.Price.OfferPrice,
            })



            if err != nil {
                    color.Red(">>>>>>>> : %+v", err.Error())
                    //return
            }
            color.Yellow("%+v", s.BaseResponse)






            return
            //fields := []Solr.Field{
            //      {
            //              Name: "url",
            //              Type: "string",
            //      },
            //}
            //status, _ := solrClient.CoreStatus(context.Background(), &Solr.CoreParams{})
            //
            //err = solrClient.AddFields(context.Background(), "go_products", fields...)


            //color.Red("%s", status, err)
            //status, qtime, err := client.Ping()
    
            //schema, err := client.Schema()
            //if err != nil {
            //      return
            //}
    
            //schema.SetCore("go_products")
            //_, _ = schema.Fields("url", false, true)
    
            //post, err := schema.Post("fields", schema)
            //if err != nil {
            //      return
            //
            //}
    
            //_, prettyS := utils.PrintPrettyJson(status)
            //color.Red("%s", prettyS)
            //return
            //
            //product := fake.CreateProductWithId()
            //
            //
            //
            //product.Id = "947cb24a-53c0-4078-bbcd-3e4262dc8ae0"


            //var docs []Solr.Document
            //
            //docs = append(docs, Solr.Document{"product": product})
            //update, errUpdate:= client.Add(docs, 10, &url.Values{})
    
            //_, r := utils.JsonDumpStr(update.Result)
    
            //err, pretty := utils.PrintPrettyJson(update.Result)
            //_, prettyProduct := utils.PrintPrettyJson(product)
            //if err != nil {
            //      return
            //}
            //admin, errAdmin := client.CoreAdmin()
            //admin.SetBasicAuth("vortex", "324252")
    
            //err, data := utils.JsonDump(update.Result)
            //color.Yellow("%s", prettyProduct)
            //color.Red("%s", pretty)
            //println(update.Result, update.Success, product.Id, errUpdate)
    }

-   These are usually excluded since many testcases is not documented anyhow

-   As of *go 1.16* it is recommended to **only** use module based parsing

### Test bool

Test denotes if test files (ending with \_test.go) should be included or not (default not included)

### Internal bool

Internal determines if internal folders are included or not (default not)

### UnderScore bool

UnderScore, when set to true it will include directories beginning with \_

## CollectionSchema

    type CollectionSchema struct {
        Collection    string
        Schema        []Solr.Field
    }

CollectionSchema to use when invoking ParseAny, ParseSingleFileWalker, and ParseSinglePackageWalker. .CollectionSchema &lt;1> These are usually excluded since many testcases is not documented anyhow

### Collection string

### Schema \[\]Solr.Field

# Functions

## CreateSchema

    func CreateSchema(client *solr.Client,schema []Solr.Field, collection string)

## GetFieldsFromTemplate

    func GetFieldsFromTemplate(template *solr.TemplateConfig) chan CollectionSchema
