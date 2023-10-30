package db

import (
	"context"
	"log"
	"ventas/models"

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

func AddOrder(order models.OrderGRPCMessage) (string, error) {
	res, err := Collection.InsertOne(ctx, order)
	if err != nil {
		log.Panicf("Error al insertar en la base de datos: %v", err)
		return "", err
	}

	orderID := res.InsertedID.(primitive.ObjectID).Hex()

	return orderID, nil

}
