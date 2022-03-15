package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"grpc-api/protobuf"
)

type userDetails struct {
	Uid         int32
	Name        string
	Nationality string
	Zip         int32
}

var users = []userDetails{
	{
		Uid:         1,
		Name:        "Josh Winters",
		Nationality: "American",
		Zip:         10111,
	},
	{
		Uid:         2,
		Name:        "Brian Stone",
		Nationality: "British",
		Zip:         20212,
	},
}

func toUser(data userDetails) *protobuf.User {
	return &protobuf.User{
		Uid:         data.Uid,
		Name:        data.Name,
		Nationality: data.Nationality,
		Zip:         data.Zip,
	}
}

func fromUser(user *protobuf.User) userDetails {
	return userDetails{
		Uid:         user.GetUid(),
		Name:        user.GetName(),
		Nationality: user.GetNationality(),
		Zip:         user.GetZip(),
	}
}

type server struct {
	protobuf.UserServiceServer
}

func (s *server) FetchUser(ctx context.Context, req *protobuf.FetchUserRequest) (*protobuf.FetchUserResponse, error) {
	fmt.Println("Fetching User")
	uid := req.GetUid()

	for _, user := range users {
		if user.Uid == uid {
			return &protobuf.FetchUserResponse{
				User: toUser(user),
			}, nil
		}
	}

	return nil, errors.New("User not found")
}

func (s *server) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.CreateUserResponse, error) {
	fmt.Println("Creating User")
	user := req.GetUser()

	data := fromUser(user)

	users = append(users, data)

	return &protobuf.CreateUserResponse{
		User: toUser(data),
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *protobuf.UpdateUserRequest) (*protobuf.UpdateUserResponse, error) {
	fmt.Println("Updating User")
	user := req.GetUser()

	data := fromUser(user)

	for i, user := range users {
		if user.Uid == data.Uid {
			users[i] = data
			return &protobuf.UpdateUserResponse{
				User: toUser(data),
			}, nil
		}
	}

	return nil, errors.New("Couldn't update user")
}

func (s *server) DeleteUser(ctx context.Context, req *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	fmt.Println("Deleting User")
	uid := req.GetUid()
	var tmpUsers []userDetails

	for i, user := range users {
		if user.Uid == uid {
			tmpUsers = append(users[:i], users[i+1:]...)
			users = tmpUsers
			fmt.Printf("User with id %d has been deleted.\n", uid)
			return &protobuf.DeleteUserResponse{
				Uid: uid,
			}, nil
		}
	}

	return nil, errors.New("User does not exist")
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Timeout: 5 * time.Second,
			},
		),
	}

	s := grpc.NewServer(opts...)
	protobuf.RegisterUserServiceServer(s, &server{})

	fmt.Println("Starting server...")
	fmt.Printf("Hosting server on: %s\n", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
