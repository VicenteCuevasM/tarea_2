package main

import (
	"log"
	"net"
	"os"

	db "ventas/mongo"
	order "ventas/order"
	rabbit "ventas/rabbitmq"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func main() {

	//Load config file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	mongo_string := "mongodb://localhost:" + os.Getenv("MONGO_PORT")

	// Conecta con MongoDB
	db.SetMongoString(mongo_string, "tarea2", "orders")
	err = db.InitMongoConnection()
	if err != nil {
		failOnError(err, "Error al conectar a MongoDB")
		return
	}

	// Inicia el exchange de RabbitMQ
	err = rabbit.InitializeExchange("tarea2")
	if err != nil {
		failOnError(err, "Error al iniciar el exchange de RabbitMQ")
		return
	}

	grpcPort := ":" + os.Getenv("GRPC_PORT")

	// Crea el listener
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		failOnError(err, "Error al crear el listener")
		return
	}
	defer lis.Close()

	s := order.Server{}
	// Crea el servidor gRPC
	grpcServer := grpc.NewServer()
	// Registra el servicio de ventas
	order.RegisterOrderServiceServer(grpcServer, &s)

	// Inicia el servidor gRPC

	if err := grpcServer.Serve(lis); err != nil {
		failOnError(err, "Error al iniciar el servidor gRPC")
		return
	}

}
