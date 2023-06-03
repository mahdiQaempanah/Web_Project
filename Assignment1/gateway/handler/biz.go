package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/gateway/grpc/biz"
	"google.golang.org/grpc"
)

func GetUsers(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5062", grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	client := biz.NewBizClient(conn)

	var getUserReq biz.GetUserRequest
	err = c.BindJSON(&getUserReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	response, err := client.GetUsers(context.Background(), &getUserReq)

	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}
