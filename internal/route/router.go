package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xreyc/grpc-core/internal/handler/rest"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/get-user-details", rest.GetUserDetails)
	return r
}
