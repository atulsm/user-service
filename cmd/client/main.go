package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/atulsm/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Parse command line flags
	page := flag.Int("page", 1, "Page number")
	pageSize := flag.Int("page_size", 10, "Number of items per page")
	flag.Parse()

	// Set up connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make the request
	req := &pb.GetUsersRequest{
		Page:     int32(*page),
		PageSize: int32(*pageSize),
	}

	resp, err := client.GetUsers(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	// Print the response
	log.Printf("Total users: %d", resp.Total)
	log.Printf("Page: %d, Page Size: %d", resp.Page, resp.PageSize)
	log.Println("\nUsers:")
	for _, user := range resp.Users {
		log.Printf("ID: %s", user.Id)
		log.Printf("Email: %s", user.Email)
		log.Printf("Name: %s %s", user.FirstName, user.LastName)
		log.Printf("Phone: %s", user.PhoneNumber)
		log.Printf("Created: %s", user.CreatedAt)
		log.Printf("Updated: %s", user.UpdatedAt)
		log.Println("---")
	}
}
