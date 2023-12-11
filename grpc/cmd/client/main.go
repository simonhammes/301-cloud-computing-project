package main

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/simonhammes/301-cloud-computing-project/grpc/api"
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

func generateFakeStudents(n int) []*api.Student {
	// Preallocate
	students := make([]*api.Student, n)

	for j := 0; j < n; j++ {
		name := fmt.Sprintf("%s %s", faker.FirstName(), faker.LastName())
		// No ID on purpose
		students[j] = &api.Student{Name: name}
	}

	return students
}

func unaryExample(client api.StudentsServiceClient) {
	request := api.GetStudentByIdRequest{Id: 3}

	log.Print("Calling GetStudentById()")
	response, err := client.GetStudentById(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Student: ID = %d, Name = %s", response.Id, response.Name)
}

func serverStreamingExample(client api.StudentsServiceClient) {
	request := api.GetStudentsRequest{PerMessage: 5}

	log.Print("Calling GetStudents()")
	stream, err := client.GetStudents(context.Background(), &request)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		for _, student := range response.Students {
			log.Printf("Student: %s", student.Name)
		}
	}
}

func clientStreamingExample(client api.StudentsServiceClient) {
	log.Print("Calling ImportStudents()")
	stream, err := client.ImportStudents(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i := 1; i <= 10; i++ {
		message := api.ImportStudentsRequest{Students: generateFakeStudents(5)}
		log.Printf("Importing %d students", len(message.Students))
		err := stream.Send(&message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Do some work
		time.Sleep(500 * time.Millisecond)
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Summary: Imported %d students", summary.Count)
}

func bidirectionalStreamingExample(client api.StudentsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Print("Calling ImportStudentsV2()")
	stream, err := client.ImportStudentsV2(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create a channel
	waitc := make(chan struct{})

	// Start goroutine to receive messages from the server
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// No more messages
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Printf("Received %d students with generated IDs", len(in.Students))
		}
	}()

	// Send messages
	for i := 0; i < 5; i++ {
		message := api.ImportStudentsV2Request{Students: generateFakeStudents(5)}
		log.Printf("Importing %d students", len(message.Students))

		err := stream.Send(&message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Do some work
		time.Sleep(3 * time.Second)
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

	client := api.NewStudentsServiceClient(connection)

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
