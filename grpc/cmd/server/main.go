package main

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/simonhammes/301-cloud-computing-project/grpc/api"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

type server struct {
	api.StudentsServiceServer
}

func (s *server) GetStudentById(_ context.Context, request *api.GetStudentByIdRequest) (*api.Student, error) {
	response := &api.Student{
		Id:   request.Id,
		Name: "John Doe",
	}

	return response, nil
}

func (s *server) GetStudents(request *api.GetStudentsRequest, stream api.StudentsService_GetStudentsServer) error {
	id := int32(1)
	for i := 1; i <= 10; i++ {
		// Simulate database query/network request
		time.Sleep(time.Second)

		// Generate fake data
		students := make([]*api.Student, request.PerMessage)
		for j := 0; j < int(request.PerMessage); j++ {
			name := fmt.Sprintf("%s %s", faker.FirstName(), faker.LastName())
			students[j] = &api.Student{Id: id, Name: name}
			id++
		}

		response := api.GetStudentsResponse{Students: students}

		if err := stream.Send(&response); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ImportStudents(stream api.StudentsService_ImportStudentsServer) error {
	count := int32(0)

	for {
		message, err := stream.Recv()

		if err == io.EOF {
			// No more messages on the stream
			log.Print("Received EOF")
			summary := api.ImportStudentsResponse{Count: count}
			return stream.SendAndClose(&summary)
		}

		if err != nil {
			return err
		}

		// Process message
		for _, student := range message.Students {
			log.Printf("Importing student: %s", student.Name)
			time.Sleep(200 * time.Millisecond)
			count++
		}
	}
}

func (s *server) ImportStudentsV2(stream api.StudentsService_ImportStudentsV2Server) error {
	generateRandomNumber := func(min, max int32) int32 {
		return rand.Int31n(max-min+1) + min
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// No more incoming messages
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Importing %d students", len(in.Students))

		// Do some work
		time.Sleep(500 * time.Millisecond)

		log.Print("Generating IDs...")
		for _, student := range in.Students {
			student.Id = generateRandomNumber(1, 10000)
		}

		log.Print("Sending response...")
		message := api.ImportStudentsV2Response{Students: in.Students}
		err = stream.Send(&message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := server{}

	api.RegisterStudentsServiceServer(grpcServer, &server)

	log.Print("Starting server...")
	log.Printf("Listening on %s", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
