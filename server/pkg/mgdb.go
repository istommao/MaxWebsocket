package pkg

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholders with your credentials
// const uri = "mongodb://127.0.0.1:27017//?retryWrites=true"

type MgDB struct {
	client     *mongo.Client
	ctx        context.Context
	dbInstance *mongo.Database
}

func (w *MgDB) Close() {
	if err := w.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func NewMgDB() *MgDB {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/web3db"))
	if err != nil {
		panic(err)
	}

	db := &MgDB{
		client:     client,
		ctx:        ctx,
		dbInstance: client.Database("web3db"),
	}

	return db
}

func (w *MgDB) SaveMessage(key string, value []byte) error {

	collection := w.dbInstance.Collection("messages")

	res, err := collection.InsertOne(w.ctx, bson.M{"key": key, "value": len(value)})

	fmt.Println(res.InsertedID)
	return err
}
