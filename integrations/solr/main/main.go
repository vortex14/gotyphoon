package main

import (
	"context"
	"github.com/fatih/color"
	Solr "github.com/stevenferrer/solr-go"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/integrations/solr"
	"github.com/vortex14/gotyphoon/utils"

	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func main() {
	options := &solr.ConnectOptions{}
	client := (&solr.Client{Options: options, DOptions: &solr.DockerOptions{MatchName: "bitnami/solr:8-debian-10"}}).
		Init()

	//_, _ = client.GetCollectionSchemaFields()
	client.ConnectToDockerNode()

	return
	//_, res := utils.DumpPrettyJson(fields)

	//err := client.RemoveSchema(fields)
	//if err != nil {
	//	color.Red("%s", err.Error())
	//	return
	//}
	//color.Yellow("%s", res)

	errR, s2 := client.RemoveAllDocs()
	if errR != nil {
		color.Red("%+v", errR)
		return
	}

	_, s2P := utils.DumpPrettyJson(s2)
	color.Red("%s", s2P)
	//return

	product := fake.CreateProductWithId()

	product.Id = "947cb24a-53c0-4078-bbcd-3e4262dc8ae0"

	schema := []Solr.Field{
		{
			Name:        "description",
			Type:        "string",
			MultiValued: false,
		},
		//{
		//	Name: "upc",
		//	Type: "string",
		//	//MultiValued: false,
		//},
		//{
		//	Name: "price_listing",
		//	Type: "pdouble",
		//	//MultiValued: false,
		//
		//},
		//{
		//	Name: "price_offer",
		//	Type: "pdouble",
		//	//MultiValued: false,
		//},
	}

	//err := client.RemoveSchema(schema)
	//if err != nil {
	//	color.Red(">>>>>",err.Error())
	//}

	errI := client.AddSchemaFields(schema)
	color.Red(">>>>>> added: %+v", errI)
	return

	//
	//errI := client.AddSchemaFields([]Solr.Field{
	//	//{
	//	//	Name: "description",
	//	//	Type: "string",
	//	//},
	//	//{
	//	//	Name: "upc",
	//	//	Type: "string",
	//	//},
	//	{
	//		Name: "price_listing",
	//		Type: "pdouble",
	//		MultiValued: false,
	//
	//	},
	//	{
	//		Name: "price_offer",
	//		Type: "pdouble",
	//		MultiValued: false,
	//	},
	//})
	//
	//
	//if errI != nil {
	//	color.Red(errI.Error(), ">>>>")
	//	//return
	//}
	//
	//return

	product.Price.ListingPrice = 10.5

	err, s := client.UpdateDocAndCommit(context.Background(), solr.Document{
		"id":            product.Id,
		"url":           product.Url,
		"title":         product.Title,
		"description":   product.Description,
		"upc":           product.Upc,
		"price_listing": product.Price.ListingPrice,
		"price_offer":   product.Price.OfferPrice,
	})

	if err != nil {
		color.Red(">>>>>>>> : %+v", err.Error())
		//return
	}
	color.Yellow("%+v", s.BaseResponse)

	return
	//fields := []Solr.Field{
	//	{
	//		Name: "url",
	//		Type: "string",
	//	},
	//}
	//status, _ := solrClient.CoreStatus(context.Background(), &Solr.CoreParams{})
	//
	//err = solrClient.AddFields(context.Background(), "go_products", fields...)

	//color.Red("%s", status, err)
	//status, qtime, err := client.Ping()

	//schema, err := client.Schema()
	//if err != nil {
	//	return
	//}

	//schema.SetCore("go_products")
	//_, _ = schema.Fields("url", false, true)

	//post, err := schema.Post("fields", schema)
	//if err != nil {
	//	return
	//
	//}

	//_, prettyS := utils.DumpPrettyJson(status)
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

	//err, pretty := utils.DumpPrettyJson(update.Result)
	//_, prettyProduct := utils.DumpPrettyJson(product)
	//if err != nil {
	//	return
	//}
	//admin, errAdmin := client.CoreAdmin()
	//admin.SetBasicAuth("vortex", "324252")

	//err, data := utils.JsonDump(update.Result)
	//color.Yellow("%s", prettyProduct)
	//color.Red("%s", pretty)
	//println(update.Result, update.Success, product.Id, errUpdate)
}
