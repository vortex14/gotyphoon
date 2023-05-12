package mongo

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	MongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	"github.com/mongodb/mongo-tools/common/db"
	mongoTools "github.com/mongodb/mongo-tools/common/options"
	exportOptions "github.com/mongodb/mongo-tools/mongoexport"
	importOptions "github.com/mongodb/mongo-tools/mongoimport"
)

type Service struct {
	singleton.Singleton

	LOG      interfaces.LoggerInterface
	client   *mongo.Client
	Project  interfaces.Project
	Settings interfaces.ServiceMongo
	dbs      map[string]*mongo.Database
}

func (s *Service) GetPort() int {
	return s.Settings.GetPort()
}

func (s *Service) GetHost() string {
	return s.Settings.GetHost()
}

func (s *Service) initClient() {
	s.Construct(func() {
		var connectionString string
		if len(s.Settings.Details.Username) > 0 && len(s.Settings.Details.Password) > 0 {
			connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d/?", s.Settings.Details.Username,
				s.Settings.Details.Password,
				s.Settings.GetHost(),
				s.Settings.GetPort())

		} else {
			connectionString = fmt.Sprintf("mongodb://%s:%d", s.Settings.GetHost(), s.Settings.GetPort())
		}

		ctx := context.Background()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
		if err == nil {
			s.client = client
			s.initDbs()
			s.LOG = log.New(log.D{"service": "mongo"})
			//color.Green("Mongo client success %s", connectionString)
		} else {
			s.LOG.Error("Mongo client error: %s ", connectionString)
		}
	})
}

func (s *Service) initDbs() {
	s.dbs = map[string]*mongo.Database{}
	for _, dbName := range s.Settings.DbNames {
		s.dbs[dbName] = s.client.Database(dbName)
	}
}

func (s *Service) GetMongoCollection(dbName, collectionName string) *mongo.Collection {
	if s.client == nil {
		s.initClient()
	}
	if s.dbs == nil {
		s.LOG.Error(Errors.MongoNotFoundDBMap.Error())
		return nil
	}

	if mongoDB, ok := s.dbs[dbName]; !ok {
		s.LOG.Error(Errors.MongoNotFoundDB.Error())
		return nil
	} else {
		return mongoDB.Collection(collectionName)
	}
}

func (s *Service) GetMongoDB(dbName string) *mongo.Database {

	s.initClient()

	if s.dbs == nil {
		s.LOG.Error(Errors.MongoNotFoundDBMap.Error())
		return nil
	}

	if mongoDB, ok := s.dbs[dbName]; !ok {
		s.LOG.Error(Errors.MongoNotFoundDB.Error())
		return nil
	} else {
		return mongoDB
	}
}

func (s *Service) GetFilterOptions(_id string, value interface{}) bson.M {
	return bson.M{_id: value}
}

func (s *Service) GetUpdateOptions(values bson.D) bson.D {
	return bson.D{{"$set", values}}
}

func (s *Service) GetIncOptions(values bson.D) bson.D {
	return bson.D{{"$inc", values}}
}

func (s *Service) GetIncUpdateOptions(incValues bson.D, updateValues bson.D) bson.D {
	return bson.D{{"$inc", incValues}, {"$set", updateValues}}
}

func (s *Service) GetUpsertOptions() *MongoOptions.UpdateOptions {
	upsert := true
	return &MongoOptions.UpdateOptions{Upsert: &upsert}
}

func (s *Service) GetCollections() []*Collection {
	//color.Yellow("Get collections list from project db ...")
	var results []*Collection
	s.initClient()

	for dbName := range s.dbs {
		db := s.dbs[dbName]
		query := &interfaces.MongoQuery{
			Context: context.TODO(),
			Filter:  bson.M{},
			Query:   nil,
			Options: &options.ListCollectionsOptions{},
		}
		if query.Options != nil {
			query.Options = query.Options.(*options.ListCollectionsOptions)
		}
		collections, err := db.ListCollectionNames(
			query.Context,
			query.Filter,
			nil,
		)
		if err != nil {
			color.Red("%s", err.Error())
			os.Exit(1)
		}
		for _, collection := range collections {
			mongoData := &Collection{
				Name: collection,
				Db:   dbName,
			}

			results = append(results, mongoData)
		}

	}

	return results
}

func (s *Service) FindOne(query bson.D, collection *mongo.Collection) *bson.D {
	_doc := collection.FindOne(context.TODO(), query)
	var _decoded bson.D
	_ = _doc.Decode(&_decoded)
	return &_decoded
}

func (s *Service) GetCountDocuments(query *interfaces.MongoQuery) int64 {
	s.initClient()
	query.Options = options.Count().SetMaxTime(query.Timeout * time.Second)
	collection := s.dbs[query.Database].Collection(query.Collection)
	count, err := collection.CountDocuments(
		query.Context,
		query.Filter,
		query.Options.(*options.CountOptions),
	)

	if err != nil {
		color.Red("GetCountDocuments error: %s", err.Error())
	}

	return count
}

func (s *Service) GetDocument(query *interfaces.MongoQuery) (*mongo.SingleResult, error) {
	s.initClient()
	//var opts *options.FindOptions
	//if query.Options != nil {
	//	opts = query.Options.(*options.FindOptions)
	//}

	collection := s.dbs[query.Database].Collection(query.Collection)
	res := collection.FindOne(query.Context, query.Filter)

	return res, nil
}

func (s *Service) RemoveDocById(query *interfaces.MongoQuery) (*mongo.DeleteResult, error) {
	s.initClient()
	collection := s.dbs[query.Database].Collection(query.Collection)
	return collection.DeleteOne(query.Context, query.Filter)
}

func (s *Service) GetDocuments(query *interfaces.MongoQuery, _results *[]bson.M) error {
	s.initClient()
	var opts *options.FindOptions
	if query.Options != nil {
		opts = query.Options.(*options.FindOptions)
	}

	collection := s.dbs[query.Database].Collection(query.Collection)
	cursor, err := collection.Find(query.Context, query.Filter, opts)
	if err != nil {
		return err
	}

	if err = cursor.All(context.TODO(), _results); err != nil {
		return err
	}
	return nil
}

func (s *Service) connect() bool {
	s.initClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	status := false
	errPing := s.client.Ping(ctx, readpref.Primary())
	if errPing != nil {
		color.Red("%s", errPing)
	} else {
		status = true
	}
	return status
}

func (s *Service) Ping() bool {
	return s.connect()
}

func (s *Service) Init() {
	s.connect()
}

func (s *Service) GetConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%d", s.Settings.GetHost(), s.Settings.GetPort())
}

func (s *Service) Import(Database string, collection string, inputFile string) (error, uint64) {
	color.Yellow("Import mongo data")
	toolOptions := s.GetToolOptions(Database, collection)
	inputOptions := &importOptions.InputOptions{
		ParseGrace: "stop",
	}
	ingestOptions := &importOptions.IngestOptions{}
	provider, err := db.NewSessionProvider(*toolOptions)
	if err != nil {
		return err, 0
	}
	importDb := &importOptions.MongoImport{
		ToolOptions:     toolOptions,
		InputOptions:    inputOptions,
		IngestOptions:   ingestOptions,
		SessionProvider: provider,
	}

	importDb.IngestOptions.Mode = "insert"
	importDb.InputOptions.File = inputFile
	importDb.IngestOptions.WriteConcern = "1"

	numInserted, _, err := importDb.ImportDocuments()

	return err, numInserted

}

func (s *Service) GetToolOptions(Database string, collection string) *mongoTools.ToolOptions {
	var toolOptions *mongoTools.ToolOptions
	namespace := &mongoTools.Namespace{
		DB:         Database,
		Collection: collection,
	}
	connection := &mongoTools.Connection{
		Host: s.Settings.GetHost(),
		Port: strconv.Itoa(s.Settings.GetPort()),
	}
	toolOptions = &mongoTools.ToolOptions{
		General:    &mongoTools.General{},
		Connection: connection,
		Verbosity:  &mongoTools.Verbosity{},
		URI:        &mongoTools.URI{},
		Auth:       &mongoTools.Auth{},
		SSL: &mongoTools.SSL{
			UseSSL: false,
		},
		Namespace: namespace,
	}
	return toolOptions
}

func (s *Service) Export(Database string, collection string, outFile string) (*bufio.Writer, *os.File, int64, error) {
	f, err := os.Create(outFile)

	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
	writer := bufio.NewWriterSize(
		f,
		4096*2,
	)
	toolOptions := s.GetToolOptions(Database, collection)

	opts := exportOptions.Options{
		ToolOptions: toolOptions,
		OutputFormatOptions: &exportOptions.OutputFormatOptions{
			Type:       "json",
			JSONFormat: "canonical",
		},
		InputOptions: &exportOptions.InputOptions{},
	}

	opts.Collection = collection
	opts.DB = Database

	me, err := exportOptions.New(opts)
	if err != nil {
		return nil, nil, 0, err
	}
	defer me.Close()
	count, err := me.Export(writer)
	return writer, f, count, err

}
