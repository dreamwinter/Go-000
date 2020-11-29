package main

import (
	"log"

	"xll.com/go-000/Week02/cmd/handler"
	"xll.com/go-000/Week02/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Default already use Logger(), Recovery() middleware
	r := gin.Default()
	// initialize userService instance
	userService, err := service.NewUserService(service.UserServiceInMemoryType)
	if err != nil {
		// if user service cannot be initialized, here we need exit and report fatal error
		log.Fatalf("Cannot initialize userService due to: %v", err)
	}
	r.GET("/user/:id", handler.GetUserHander(userService))
	r.Run()
}
