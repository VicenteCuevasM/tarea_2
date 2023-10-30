package db

import (
	"context"
	"inventario/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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

// Esta funcion asegura atomicidad al momento de actualizar el stock de multiples libros.
func UpdateMultipleStock(order models.OrderMessage) error {

	books := order.Products

	var operations []mongo.WriteModel

	for _, book := range books {
		filter := bson.M{"title": book.Title, "author": book.Author, "genre": book.Genre, "pages": book.Pages, "publication": book.Publication, "price": book.Price}
		update := bson.M{"$inc": bson.M{"quantity": 0 - book.Quantity}} //En el pdf dice quantity no stock, de ser necesario solo cambiar el nombre del campo
		operations = append(operations, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	_, err := Collection.BulkWrite(ctx, operations)
	if err != nil {
		log.Panicf("Error al actualizar el stock de los libros: %v", err)
	}

	return nil

}
