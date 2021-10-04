package errors

import Errors "errors"

var MongoNotFoundDB = Errors.New("not found mongo db")
var MongoNotFoundDBMap = Errors.New("not found mongo db map. need call .init")

var SolrConnectionsOptionsNotFound = Errors.New("solr connection options not found")
var SolrConnectionEndpointError = Errors.New("solr connection has error. need check endpoint or other *ConnectOptions")