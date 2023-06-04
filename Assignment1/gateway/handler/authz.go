package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/gateway/grpc/pb"
	"google.golang.org/grpc"
)

func RequestPG(c *gin.Context) {
	conn, err := grpc.Dial("auth:5052", grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	client := pb.NewAuthzClient(conn)

	var pgRequest pb.PGRequest
	err = c.BindJSON(&pgRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	response, err := client.PGAgreement(context.Background(), &pgRequest)

	if err != nil {
		c.JSON(http.StatusBadGateway, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func DiffieHellman(c *gin.Context) {
	conn, err := grpc.Dial("auth:5052", grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	client := pb.NewAuthzClient(conn)

	var diffieHellmanRequest pb.DiffieHellmanRequest
	err = c.BindJSON(&diffieHellmanRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	response, err := client.DiffieHellman(context.Background(), &diffieHellmanRequest)

	if err != nil {
		c.JSON(http.StatusBadGateway, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}
