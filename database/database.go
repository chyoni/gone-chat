package database

import (
	"fmt"
	"os"

	"github.com/chiwon99881/gone-chat/entity"
	"github.com/chiwon99881/gone-chat/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	CreateUser(username, password string)
}

type RepoOperator struct{}

func (RepoOperator) CreateUser(username, password string) {
	createUser(username, password)
}

func NewRepository() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		os.Getenv("HOSTNAME"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default})
	if err != nil {
		utils.HandleError(err)
	}

	db.AutoMigrate(&entity.User{}, &entity.Room{})
	return db
}

func createUser(username, password string) {
	user := entity.User{Username: username, Password: password}
	NewRepository().Create(&user)
}