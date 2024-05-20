package db

import (
	"Demo_Api/gen-go/example"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm")

func Connect() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = mysqlDB.AutoMigrate(&example.User{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Kết nối MySQL thành công")

	return mysqlDB, nil
}
