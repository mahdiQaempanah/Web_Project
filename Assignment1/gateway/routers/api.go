package routers

import (
	"github.com/mahdiQaempanah/Web_Project/Assignment1/gateway/handler"

	"github.com/gin-gonic/gin"
)

func Api() *gin.Engine {
	router := gin.Default()

	router.GET("/api/v1/authz/pg", handler.RequestPG)
	router.GET("/api/v1/authz/dh", handler.DiffieHellman)

	router.GET("/api/v1/biz/get", handler.GetUsers)
	router.GET("/api/v1/biz/getwithinj", handler.GetUsersWithSQLInject)

	// router.RunTLS("0.0.0.0:6433", "/certs/server.crt", "/certs/server.key")
	router.Run("0.0.0.0:6433")
	return router
}
