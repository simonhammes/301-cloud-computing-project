package main

import (
	"context"
	"fmt"
	"github.com/simonhammes/301-cloud-computing-project/grpc/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// Disable TLS
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, err := grpc.Dial("127.0.0.1:3000", options)
	if err != nil {
		log.Fatalf("Could not create connection: %v", err)
	}

	// Close connection before exiting
	defer connection.Close()

	client := search.NewSearchServiceClient(connection)

	request := search.SearchRequest{
		Query:          "search term",
		PageNumber:     5,
		ResultsPerPage: 10,
	}

	response, err := client.Search(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("%+v", response)
}
