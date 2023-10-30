package db

import (
	"context"
	"despacho/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var connectionString = ""
var dbName = ""
var colName = ""
var db *mongo.Database
var Collection *mongo.Collection

func InitMongoConnection() error {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panicf("Error al conectar a la base de datos: %v", err)
		return err
	}
	db = client.Database(dbName)
	Collection = db.Collection(colName)
	return nil
}

func GetMongoCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

/** Connect to MongoDB
 * @param connectionString string
 * @param dbName string
 * @param colName string
 * @return error
 */
func SetMongoString(mongo_string string, mongo_db string, mongo_col string) {
	connectionString = mongo_string
	dbName = mongo_db
	colName = mongo_col
}

func AddDeliveryToOrder(id string, deliveries []models.Delivery) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"deliveries": deliveries}}
	_, err = Collection.UpdateOne(ctx, filter, update)
	return err
}
