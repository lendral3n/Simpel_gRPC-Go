package config

import (
	"log"
	srv "l3ngrpc/cmd/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/grpcsimpel?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		log.Fatal("failed conection database", err.Error())
	}

	db.AutoMigrate(
		&srv.Product{},
	)
	return db
}
