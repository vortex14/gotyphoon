package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vortex14/gotyphoon/elements/models/label"
)

type DataBase struct {
	*label.MetaInfo
	db          *mongo.Database
	client      *mongo.Client
	Collections map[string]*Collection
}

func (m *DataBase) GetMongoCollections() {

}

func (m *DataBase) GetMongoCollection(collectionName string) *mongo.Collection {
	//if m.db == nil {  }
	return m.db.Collection(collectionName)
}

func (m *DataBase) Export() {

}

func (m *DataBase) Import() {

}
