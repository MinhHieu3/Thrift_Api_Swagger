package main

import (
	"Demo_Api/gen-go/example"
	"context"
	"fmt"
	"net/http"
	_ "Demo_Api/docs"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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


// @Summary Create 
// @Description Create a new user with the provided data
// @Param   user     body    example.User    true        "User data"
// @Success 200 {object} example.User "Successfully created user"
// @Success 200 {object} []example.User "List of users"
// @Failure 400 {object} []example.TErrorCode "Invalid request body"
// @Router /api/users [post]
func Create(c *gin.Context) {
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

// @Summary Update
// @Description Update an existing user with the provided data
// @Param   id     path    string     true        "User ID"
// @Param   user     body    example.User    true        "User data"
// @Success 200 {object} []example.User "Successfully updated user"
// @Failure 400 {object} []example.TErrorCode "Invalid request body"
// @Failure 404 {object} []example.TErrorCode "User not found"
// @Router /api/users/{id} [put]
func Update(c *gin.Context) {
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
// FindAll return list of all user fromm the database
// @Summary Get all 
// @Description Get a list of all users
// @Success 200 {object} []example.User "List of users"
// @Router /api/users [get]
func FindAll(c *gin.Context) {
	listData, err := thriftClient.client.GetListUser(context.Background(), []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, listData)
  }
// @Summary Delete 
// @Description Delete a user by ID
// @Param   id     path    string     true        "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 404 {object} []example.TErrorCode "User not found"
// @Router /api/users/{id} [delete]
func Delete(c *gin.Context) {
	id := c.Param("id")
	_, err := thriftClient.client.RemoveUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete user: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
// @title Documenting Api
// @version 1
// @Description Sample decription
// @contact.url https://github.com/demo
// @host localhost:8080
// @BasePath /api/users
func main() {
	r := gin.Default()
	router:=r.Group("api/users")
    {
		router.POST("/", Create)
		router.GET("/", FindAll)
		router.PUT("/:id",Update)
		router.DELETE("/:id",Delete)
   }
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

r.Run()
}
