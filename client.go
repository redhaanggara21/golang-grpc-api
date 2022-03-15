package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"grpc-api/protobuf"
)

func main() {
	fmt.Println("Starting Client\n")

	cc, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	// Close connection before exiting app
	defer cc.Close()

	c := protobuf.NewUserServiceClient(cc)

	// Create a user
	fmt.Println("Calling CreateUser()")
	user := &protobuf.User{
		Uid:         3,
		Name:        "Sarah Connors",
		Nationality: "Canadian",
		Zip:         45015,
	}

	createUserResponse, err := c.CreateUser(context.Background(), &protobuf.CreateUserRequest{User: user})
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
	fmt.Printf("User has been created: %v\n\n", createUserResponse)

	// userID := createUserResponse.GetUser().GetId()

	// Update a user
	updateUser := &protobuf.User{
		Uid:         1,
		Name:        "Mandy Williams",
		Nationality: "American",
		Zip:         10111,
	}

	fmt.Println("Updating a user")
	updateUserResponse, err := c.UpdateUser(context.Background(), &protobuf.UpdateUserRequest{User: updateUser})
	if err != nil {
		log.Fatalf("Error occurred updating user: %v", err)
	}
	fmt.Printf("A user has been updated: %v\n\n", updateUserResponse)

	fmt.Println("Deleting user")
	var delUid int32 = 2
	deleteUserResponse, err := c.DeleteUser(context.Background(), &protobuf.DeleteUserRequest{Uid: delUid})
	if err != nil {
		log.Fatalf("Error occurred deleting user: %v", err)
	}
	fmt.Printf("User with ID: %d has been deleted.\n\n", deleteUserResponse.GetUid())

	// Get a user
	fmt.Println("Fetching a user")
	var getUid int32 = 1
	fetchUserResponse, err := c.FetchUser(context.Background(), &protobuf.FetchUserRequest{Uid: getUid})
	if err != nil {
		log.Fatalf("User not found error: %v", getUid)
	}
	fmt.Printf("User: %v\n", fetchUserResponse)
}
