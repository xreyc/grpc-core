package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authv1 "github.com/xreyc/grpc-core/internal/gen/go/auth/v1"
	grpcClient "github.com/xreyc/grpc-core/internal/grpc"
)

func GetUserDetails(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		username = "xreyc" // default for testing
	}

	resp, err := grpcClient.AuthClient.GetUserDetails(c, &authv1.GetUserRequest{
		Username: username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  resp.GetUsername(),
		"email":     resp.GetEmail(),
		"full_name": resp.GetFullName(),
	})
}
