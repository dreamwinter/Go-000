package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"xll.com/go-000/Week02/internal/service"
)

// GetUserHander creates the UserHandler
func GetUserHander(userService service.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			log.Printf("Encounter error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": "Invalid request",
			})
			return
		}
		log.Printf("[GetUser] Request id: %v", id)
		user, err := userService.GetUser(id)
		if err != nil {
			if service.IsUserNotFound(err) {
				log.Printf("Cannot find the user with id: %v", id)
				c.JSON(http.StatusNotFound, gin.H{
					"Message": "Cannot find the user",
				})

			} else {
				log.Printf("Encountered error: %T %v\n", errors.Cause(err), errors.Cause(err))
				log.Printf("Stack trace:\n%+v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"Message": "Server Internal Error",
				})
			}
			return
		}

		c.JSON(http.StatusOK, user)
		log.Printf("[GetUser] Response: %#v", user)
	}
}
