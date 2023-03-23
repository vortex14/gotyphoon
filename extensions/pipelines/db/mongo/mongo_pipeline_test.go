package mongo

import (
	Context "context"
	"fmt"
	NSQ "github.com/segmentio/nsq-go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"go.mongodb.org/mongo-driver/bson"
	M "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"testing"

	"encoding/json"
	"github.com/vortex14/gotyphoon/integrations/mongo"
)

func init() {
	log.InitD()
}

//func CreateMongoServiceWithoutAuth(name string, host string, port int, dbNames []string) *mongo.Service  {
//	return &mongo.Service{
//		Settings: interfaces.ServiceMongo{
//			DbNames: dbNames,
//			Name: name,
//			Details: struct {
//				AuthSource string `yaml:"authSource,omitempty"`
//				Username   string `yaml:"username,omitempty"`
//				Password   string `yaml:"password,omitempty"`
//				Host       string `yaml:"host"`
//				Port       int    `yaml:"port"`
//			}{
//				Host: host,
//				Port: port,
//			},
//		},
//	}
//}

func TestCtxPipeline(t *testing.T) {

	Convey("test", t, func() {

		p := &Pipeline{
			BasePipeline: &forms.BasePipeline{
				NotIgnorePanic: true,
				MetaInfo: &label.MetaInfo{
					Name:        "Mongo",
					Description: "Mongo pipeline",
				},
			},
			opts: &interfaces.ServiceMongo{
				DefaultCollection: "test-collection",
				DefaultDatabase:   "test-db",
				DbNames:           []string{"test-db"},
				Name:              "database",
				Details: struct {
					AuthSource string `yaml:"authSource,omitempty"`
					Username   string `yaml:"username,omitempty"`
					Password   string `yaml:"password,omitempty"`
					Host       string `yaml:"host"`
					Port       int    `yaml:"port"`
				}{
					Host: "localhost",
					Port: 27017,
				},
			},

			Fn: func(context Context.Context,
				task interfaces.TaskInterface,
				logger interfaces.LoggerInterface,
				service *mongo.Service,
				database *M.Database,
				collection *M.Collection) (error, Context.Context) {

				So(collection.Name(), ShouldEqual, "test-collection")
				So(database.Name(), ShouldEqual, "test-db")
				return nil, context

			},
			Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {

			},
		}

		L := log.New(map[string]interface{}{"test": "test"})

		ctx := log.NewCtx(Context.Background(), L)

		ctx = task.PatchCtx(ctx, fake.CreateDefaultTask())

		p.Run(ctx, func(pipeline interfaces.BasePipelineInterface, err error) {
			L.Errorf("%+v", err)
		}, func(ctx Context.Context) {

		})
	})

}

func TestMongoGroup(t *testing.T) {
	Convey("test mongo group", t, func() {

		opts := &interfaces.ServiceMongo{
			DefaultCollection: "test-collection",
			DefaultDatabase:   "test-db",
			DbNames:           []string{"test-db", "data"},
			Name:              "db",
			Details: struct {
				AuthSource string `yaml:"authSource,omitempty"`
				Username   string `yaml:"username,omitempty"`
				Password   string `yaml:"password,omitempty"`
				Host       string `yaml:"host"`
				Port       int    `yaml:"port"`
			}{
				Host: "localhost",
				Port: 27017,
			},
		}

		p := forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{Name: "mongo group"},
			Options:  forms.GetNotRetribleOptions(),
			Stages: []interfaces.BasePipelineInterface{
				&Pipeline{
					BasePipeline: &forms.BasePipeline{
						NotIgnorePanic: true,
						MetaInfo: &label.MetaInfo{
							Name:        "Mongo",
							Description: "Mongo pipeline",
						},
					},
					opts: opts,

					Fn: func(context Context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface,
						service *mongo.Service,
						database *M.Database,
						collection *M.Collection) (error, Context.Context) {

						So(collection.Name(), ShouldEqual, "test-collection")
						So(database.Name(), ShouldEqual, "test-db")
						return nil, context

					},
					Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {

					},
				},
				&Pipeline{
					BasePipeline: &forms.BasePipeline{
						NotIgnorePanic: true,
						MetaInfo: &label.MetaInfo{
							Name:        "Mongo",
							Description: "Mongo pipeline",
						},
					},
					opts: opts.RenameDefaultDatabase("data").RenameDefaultCollection("LA-LA-DBA-DBA"),

					Fn: func(context Context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface,
						service *mongo.Service,
						database *M.Database,
						collection *M.Collection) (error, Context.Context) {

						So(collection.Name(), ShouldEqual, "LA-LA-DBA-DBA")
						So(database.Name(), ShouldEqual, "data")
						return nil, context

					},
					Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {

					},
				},
			},
		}

		L := log.New(map[string]interface{}{"test": "test"})

		ctx := log.NewCtx(Context.Background(), L)

		ctx = task.PatchCtx(ctx, fake.CreateDefaultTask())

		err := p.Run(ctx)

		So(err, ShouldBeNil)
	})
}

func serializeTask(doc bson.M) []byte {
	_doc, _ := json.Marshal(doc["task"])
	return _doc
}

func deserializeTask(doc bson.M) *task.TyphoonTask {
	docM, _ := json.Marshal(doc["task"])
	var _task task.TyphoonTask
	_ = json.Unmarshal(docM, &_task)
	return &_task
}

func TestMigrateProcessorExceptions(t *testing.T) {
	Convey("test migrate exceptions to NSQ ", t, func() {

		producer, _ := NSQ.StartProducer(NSQ.ProducerConfig{
			Topic:   "hello10100101010101010",
			Address: "localhost:4150",
		})

		opts := &interfaces.ServiceMongo{
			DefaultCollection: "processor_exceptions",
			DefaultDatabase:   "",
			DbNames:           []string{""},
			Name:              "db",
			Details: struct {
				AuthSource string `yaml:"authSource,omitempty"`
				Username   string `yaml:"username,omitempty"`
				Password   string `yaml:"password,omitempty"`
				Host       string `yaml:"host"`
				Port       int    `yaml:"port"`
			}{},
		}

		p := forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{Name: "mongo group"},
			Options:  forms.GetNotRetribleOptions(),
			Stages: []interfaces.BasePipelineInterface{
				&Pipeline{
					BasePipeline: &forms.BasePipeline{
						NotIgnorePanic: true,
						MetaInfo: &label.MetaInfo{
							Name:        "Mongo check",
							Description: "Mongo pipeline for check collection",
						},
					},
					opts: opts,

					Fn: func(context Context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface,
						service *mongo.Service,
						database *M.Database,
						collection *M.Collection) (error, Context.Context) {

						query := bson.D{{"exception_id", "b860f84703ec2a0b79ae8de74e1f70aa"}}

						c, _ := collection.CountDocuments(context, query)
						logger.Warning(c)

						var results []bson.M
						limit := int64(100)

						cursor, err := collection.Find(context, query, &options.FindOptions{Limit: &limit})
						if err != nil {
							return err, context
						}

						if err = cursor.All(Context.TODO(), &results); err != nil {
							return err, context
						}

						task.SetSaveData("data", &results)

						return nil, context

					},
					Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {

					},
				},
				&Pipeline{
					BasePipeline: &forms.BasePipeline{
						Options:        &forms.Options{ProgressBar: true},
						NotIgnorePanic: true,
						MetaInfo: &label.MetaInfo{
							Name:        "processing exceptions",
							Description: "Mongo pipeline for processing exceptions",
						},
					},
					opts: opts,
					Fn: func(context Context.Context,
						task interfaces.TaskInterface,
						logger interfaces.LoggerInterface,
						service *mongo.Service,
						database *M.Database,
						collection *M.Collection) (error, Context.Context) {

						list := task.GetSaveData("data").(*[]bson.M)

						wg := sync.WaitGroup{}

						_, _bar := forms.GetBar(context)

						//logger.Warning(_bar)
						//
						//return nil, context

						_bar.NewOption(int64(0), int64(len(*list)))

						//for i := 0; i < 100; i++ {
						//	_bar.Increment()
						//	time.Sleep(10 * time.Millisecond)
						//}

						for _, doc := range *list {
							_t := serializeTask(doc)
							wg.Add(1)
							go func(w *sync.WaitGroup) {
								err := producer.Publish(_t)
								if err != nil {
									logger.Error(err)
								} else {
									_bar.Increment()
									//logger.Debug("Published !")
								}
								w.Done()

							}(&wg)
							//logger.Debug(_t.Taskid)
						}

						wg.Wait()
						_bar.Finish()
						logger.Info("DONE !!")

						//logger.Debug(len(*list))

						return nil, context

					},
					Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {

					},
				},
			},
		}

		L := log.New(map[string]interface{}{"test": "test"})

		ctx := log.NewCtx(Context.Background(), L)

		ctx = task.PatchCtx(ctx, fake.CreateDefaultTask())

		err := p.Run(ctx)

		So(err, ShouldBeNil)
	})
}

func TestCountPipeline(t *testing.T) {

	m := &mongo.Service{
		Settings: interfaces.ServiceMongo{
			DbNames: []string{""},
			Name:    "",
			Details: struct {
				AuthSource string `yaml:"authSource,omitempty"`
				Username   string `yaml:"username,omitempty"`
				Password   string `yaml:"password,omitempty"`
				Host       string `yaml:"host"`
				Port       int    `yaml:"port"`
			}{
				Host:     "",
				Port:     12333,
				Password: "",
				Username: "",
			},
		},
	}
	producer, _ := NSQ.StartProducer(NSQ.ProducerConfig{
		Topic:   "hello10100101010101010",
		Address: "",
	})

	// Publishes a message to the topic that this producer is configured for,
	// the method returns when the operation completes, potentially returning an
	// error if something went wrong.
	producer.Publish([]byte("Hello World!"))

	println(len(m.GetCollections()))

	count := m.GetCountDocuments(&interfaces.MongoQuery{
		Collection: "processor_exceptions",
		Database:   "",
		Context:    Context.Background(),
		Filter:     bson.M{},
		Query:      nil,
	}) //mongo.CreateMongoServiceWithoutAuth()

	println(count)
	limit := int64(1)
	var results []bson.M

	_ = m.GetDocuments(&interfaces.MongoQuery{
		Collection: "processor_exceptions",
		Database:   "",
		Context:    Context.Background(),
		Filter:     bson.M{},
		Query:      nil,
		Options: &options.FindOptions{
			Limit: &limit,
		},
	}, &results)

	for _, doc := range results {

		//println(doc[bson.TypeObjectID.String()])
		//_id := doc["_id"].()

		//println(fmt.Sprintf("%+v", _id.ObjectId))

		docM, _ := json.Marshal(doc["task"])
		//println(fmt.Sprintf("%s", docM))
		var test task.TyphoonTask
		_ = json.Unmarshal(docM, &test)

		//_, s := utils.DumpPrettyJson(&test)

		//println(s)

		//primitive.ObjectID{Ob}

		val, _ := mongo.GetObjectID("641ae377363fd2a1cdbff580")

		del, _ := m.RemoveDocById(&interfaces.MongoQuery{
			Collection: "processor_exceptions",
			Database:   "",
			Context:    Context.TODO(),
			//Filter:     bson.D{{"_id", doc["_id"]}},
			Filter: bson.M{"_id": val},
			//Query:      bson.D{{"_id", fmt.Sprintf("1sdf %+v", doc["_id"])}},
		})

		println(del.DeletedCount)

		r, _ := m.GetDocument(&interfaces.MongoQuery{
			Collection: "processor_exceptions",
			Database:   "",
			Context:    Context.TODO(),
			//Filter:     bson.D{{"_id", doc["_id"]}},
			Filter: bson.M{"_id": val},
			//Query:      bson.D{{"_id", fmt.Sprintf("1sdf %+v", doc["_id"])}},
		})

		//println(fmt.Sprintf("%+v", r))
		var _doc bson.M
		//64159a8488a9ec73196daa2e
		println(fmt.Sprintf("%+v", r.Decode(&_doc)))

		println(fmt.Sprintf("%+v", _doc))

		//var test task.TyphoonTask
		//err := bson.Unmarshal(docM, &test)
		//println(err)
		//
		//_, s := utils.DumpPrettyJson(test)
		//println(s)

		//println(fmt.Sprintf("%s", string(docM)))
		//_task := doc["task"].(*task.TyphoonTask)
		//println(fmt.Sprintf("%+v", _task))
	}
	//fmt.Println(docs)

}
