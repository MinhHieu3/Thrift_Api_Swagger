package main

import (
	"Demo_Api/gen-go/example"
	"context"
	"fmt"
	"net/http"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
)

type ThriftClient struct {
	client *example.UserServiceClient
}

func NewThriftClient() (*ThriftClient, error) {
	transport, err := thrift.NewTSocket("localhost:9090")
	if err != nil {
		return nil, err
	}

	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport, err := transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	client := example.NewUserServiceClientFactory(useTransport, protocolFactory)

	if err := transport.Open(); err != nil {
		return nil, err
	}

	return &ThriftClient{client: client}, nil
}

var thriftClient *ThriftClient
var err error

func init() {
	thriftClient, err = NewThriftClient()
	if err != nil {
		fmt.Println("Error creating Thrift client:", err)
	}
}

func Post(c *gin.Context) {
	var user example.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("Failed to parse request body: %s", err.Error()),
        })
        return
    }

    dataResult, err := thriftClient.client.PostUser(context.Background(), &user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Failed to post user: %s", err.Error()),
        })
        return
    }

    c.JSON(http.StatusOK, dataResult)
}
func Put(c *gin.Context) {
	id := c.Param("id")
	var update example.User
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	data, err := thriftClient.client.PutUser(context.Background(), id, &update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetAll(c *gin.Context) {
	listData, err := thriftClient.client.GetListUser(context.Background(), []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, listData)
  }

func Delete(c *gin.Context) {
	id := c.Param("id")
	_, err := thriftClient.client.RemoveUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	router := gin.Default()
	rou :=router.Group("/api")
    {
	rou.POST("/get", Post)
	rou.GET("/getAll", GetAll)
	rou.PUT("/:id",Put)
	rou.DELETE("/:id",Delete)
   }
	

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
