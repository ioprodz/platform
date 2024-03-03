package db

import (
	"context"
	"ioprodz/common/config"
	"ioprodz/common/policies"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbInstance *mongo.Database
var dbInit = false

func GetInstance() *mongo.Database {

	if dbInit {
		return dbInstance
	}
	uri := config.Load().DB_MONGO_URI
	if uri == "" {
		panic("Mongo Uri missing")
	}
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	db := client.Database("ioprodz")
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	dbInit = true
	dbInstance = db
	return db

}

type BaseMongoRepository[T policies.Entity] struct {
	collection *mongo.Collection
}

func (repo *BaseMongoRepository[T]) List() ([]T, error) {
	var list []T

	cur, err := repo.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return list, err
	}

	for cur.Next(context.TODO()) {
		var b T
		err := cur.Decode(&b)
		if err != nil {
			return list, err
		}

		list = append(list, b)
	}

	return list, nil
}

func (repo *BaseMongoRepository[T]) Get(id string) (T, error) {

	var data *T
	err := repo.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&data)
	if err != nil {
		var empty T
		return empty, &policies.StorageError{Message: "Element not found by id: " + id + "(" + err.Error() + ")"}
	}
	return *data, nil

}

func (repo *BaseMongoRepository[T]) Delete(id string) error {
	_, err := repo.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return &policies.StorageError{Message: "Element could not be deleted by id: " + id}
	} else {
		return nil
	}
}

func (repo *BaseMongoRepository[T]) Create(entity T) error {
	_, err := repo.collection.InsertOne(context.TODO(), entity)
	return err
}

func (repo *BaseMongoRepository[T]) Update(entity T) error {
	_, err := repo.collection.UpdateOne(context.TODO(), bson.M{"id": entity.GetId()}, bson.M{"$set": entity})
	if err != nil {
		return &policies.StorageError{Message: "Element not found by id: " + entity.GetId()}
	} else {
		return nil
	}
}

func CreateMongoRepo[T policies.Entity](collectionName string) *BaseMongoRepository[T] {
	repo := &BaseMongoRepository[T]{collection: GetInstance().Collection(collectionName)}
	return repo
}
