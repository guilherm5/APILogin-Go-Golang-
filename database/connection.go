package database

import (
	"fmt"

	"github.com/guilherm5/GoAndReact/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("guilherme:123@tcp(127.0.0.1:3306)/dbreact?charset=utf8"), &gorm.Config{})
	if err != nil {
		panic("Falha ao conectar com o banco de dados")
	} else {
		fmt.Println("Sucesso ao realizar conexão com banco de dados")
	}

	DB = connection
	connection.AutoMigrate(&models.UserAuth{})

}
