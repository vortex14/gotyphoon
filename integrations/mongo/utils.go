package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetObjectID(_id string) (primitive.ObjectID, error) {
	id := primitive.ObjectID{}

	return id, id.UnmarshalText([]byte(_id))
}
