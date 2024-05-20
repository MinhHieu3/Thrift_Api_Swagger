package main

import (
	"Demo_Api/db"
	"Demo_Api/gen-go/example"
	"context"
	"fmt"
	"log"
	"github.com/apache/thrift/lib/go/thrift"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserServiceHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) GetListUser(ctx context.Context, _ []string) (*example.TListDataResult_, error) {
	var users []*example.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &example.TListDataResult_{
		ErrorCode: example.TErrorCode_EGood,
		Data:      users,
	}, nil
}

func (u *UserHandler) PostUser(ctx context.Context, data *example.User) (*example.TDataResult_, error) {
	err := u.db.Save(data).Error
	if err != nil {
		return &example.TDataResult_{ErrorCode: example.TErrorCode_EUnknown}, err
	}
	return &example.TDataResult_{ErrorCode: example.TErrorCode_EGood, Data: data}, nil
}

func (u *UserHandler) PutUser(ctx context.Context, key string, data *example.User) (*example.TDataResult_, error) {
	var user example.User
	err := u.db.Where("id = ?", key).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return  &example.TDataResult_{ErrorCode: example.TErrorCode_EUnknown}, err
		}
		return  &example.TDataResult_{ErrorCode: example.TErrorCode_EUnknown}, err
	}
	user.Name = data.Name
	user.Age=data.Age
	err = u.db.Save(&user).Error
	if err != nil {
		return  &example.TDataResult_{ErrorCode: example.TErrorCode_EUnknown}, err
	}
	return  &example.TDataResult_{ErrorCode: example.TErrorCode_EGood, Data: data}, nil
}

func (u *UserHandler) RemoveUser(ctx context.Context, key string) (example.TErrorCode, error) {
	err := u.db.Where("id = ?", key).Delete(&example.User{}).Error
	if err != nil {
		return example.TErrorCode_EUnknown, err
	}
	return example.TErrorCode_EGood, nil
}

func main() {
	dbs, err := db.Connect()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	handler := NewUserServiceHandler(dbs)
	processor := example.NewUserServiceProcessor(handler)

	serverTransport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		log.Fatal("Error creating server socket:", err)
	}
	server := thrift.NewTSimpleServer2(processor, serverTransport)
	fmt.Println("Starting the server...")
	if err := server.Serve(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}