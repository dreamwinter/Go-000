package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"xll.com/go-000/Week04/api"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.GetUser(ctx, &api.GetUserReq{Id: "1"})
	if err != nil {
		log.Printf("GetUser error: %v\n", err)
	} else {
		log.Printf("GetUser: %+v\n", r.GetUser())
	}
	r1, err := c.CreateUser(ctx, &api.CreateUserReq{
		User: &api.UserMessage{
			Id:     "1",
			Name:   "XLL",
			Gender: "M",
			Age:    int32(39),
		},
	})
	if err != nil {
		log.Printf("CreateUser error: %v\n", err)
		return
	}
	log.Printf("CreateUser: %+v\n", r1.GetUser())

	r2, err := c.ListUsers(ctx, &api.ListUsersReq{})
	if err != nil {
		log.Printf("ListUsers error: %v\n", err)
	} else {
		log.Printf("List User:\n")
		for {
			listUsersRes, err := r2.Recv()
			if err != nil {
				break
			}
			log.Printf("%+v\n", listUsersRes)
		}
	}

	r3, err := c.UpdateUser(ctx, &api.UpdateUserReq{
		User: &api.UserMessage{
			Id:     "1",
			Name:   "CDQ",
			Gender: "F",
			Age:    int32(39),
		},
	})
	if err != nil {
		log.Printf("UpdateUser error: %v\n", err)
		return
	}
	log.Printf("UpdateUser: %+v\n", r3.GetUser())

	r4, err := c.ListUsers(ctx, &api.ListUsersReq{})
	if err != nil {
		log.Printf("ListUsers error: %v\n", err)
	} else {
		log.Printf("List User after update:\n")
		for {
			listUsersRes, err := r4.Recv()
			if err != nil {
				break
			}
			log.Printf("%+v\n", listUsersRes)
		}
	}
	r5, err := c.DeleteUser(ctx, &api.DeleteUserReq{
		Id: "1",
	})
	if err != nil {
		log.Printf("DeleteUser error: %v\n", err)
		return
	}
	log.Printf("DeleteUser: %+v\n", r5.GetSuccess())

	r6, err := c.ListUsers(ctx, &api.ListUsersReq{})
	if err != nil {
		log.Printf("ListUsers error: %v\n", err)
	} else {
		log.Printf("List User after delete:\n")
		for {
			listUsersRes, err := r6.Recv()
			if err != nil {
				break
			}
			log.Printf("%+v\n", listUsersRes)
		}
	}
}
