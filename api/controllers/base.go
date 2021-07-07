package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saranrajv123/levelupblog/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %s database ", DbDriver)
		log.Fatal("this is the error ", err)
	} else {
		fmt.Printf("We are connecting to the %s database", DbDriver)
	}

	server.DB.Debug().AutoMigrate(
		&models.User{},
	)

	server.Router = gin.Default()
	server.InitializeRouter()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port --- 8080 ", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
