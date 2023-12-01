package main

import (
	"context"
	"fmt"
	"github.com/simonhammes/301-cloud-computing-project/grpc/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"time"
)

func usage() string {
	return "Usage: client <unary|server-streaming|client-streaming|bidirectional>"
}

func unaryExample(client search.SearchServiceClient) {
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

func serverStreamingExample(client search.SearchServiceClient) {
	request := search.PrimeGenerationRequest{Limit: 100}

	stream, err := client.GeneratePrimes(context.Background(), &request)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		prime, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Next Prime: %d\n", prime.Prime)
	}
}

func clientStreamingExample(client search.SearchServiceClient) {
	stream, err := client.RecordLogMessages(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i := 1; i <= 10; i++ {
		message := search.LogMessage{
			Message: fmt.Sprintf("Log message %d", i),
		}
		err := stream.Send(&message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Do some work
		time.Sleep(time.Second)
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Summary: %v", summary)
}

func bidirectionalStreamingExample(client search.SearchServiceClient) {
	messages := []*search.LogMessage{
		{Message: "1st message"},
		{Message: "2nd message"},
		{Message: "3rd message"},
		{Message: "4th message"},
		{Message: "5th message"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.LogChat(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	waitc := make(chan struct{})

	// Start goroutine to receive messages
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Received message: %s", in.Message)
		}
	}()

	// Send messages
	for _, message := range messages {
		log.Printf("Sending message: %s", message.Message)
		err := stream.Send(message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Do some work
		time.Sleep(time.Second)
	}

	stream.CloseSend()
	<-waitc
}

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

	if len(os.Args) < 2 {
		println(usage())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "unary":
		unaryExample(client)
	case "server-streaming":
		serverStreamingExample(client)
	case "client-streaming":
		clientStreamingExample(client)
	case "bidirectional":
		bidirectionalStreamingExample(client)
	default:
		println(usage())
		os.Exit(1)
	}
}
