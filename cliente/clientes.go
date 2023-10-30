package main

import (
	"cliente/order"
	"cliente/utils"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	//Load config file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var grpcServer = os.Getenv("MV2_IP")
	var grpcPort = os.Getenv("GRPC_PORT")
	var grpcAddress = grpcServer + ":" + grpcPort

	// Open our jsonFile
	jsonFile, err := utils.LoadProducts()
	if err != nil {
		log.Printf("Error opening file %s", err)
		panic(err)
	}

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	// Create a new instance of the ProductServiceClient
	// using the connection

	client := order.NewOrderServiceClient(conn)
	msg := utils.StructToProtobuf(jsonFile)

	response, err := client.CreateOrder(context.Background(), msg)
	if err != nil {
		log.Fatalf("Error when calling CreateOrder: %s", err)
	}
	log.Printf("Id de orden: %s", response)

}
