package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/chiwon99881/gone-chat/entity"
	"github.com/chiwon99881/gone-chat/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	CreateUser(username, password, alias string)
	FindUser(userID uint) (*entity.User, error)
	CreateRoom(participant uint) (*entity.Room, error)
}

type RepoOperator struct{}

func (RepoOperator) CreateUser(username, password, alias string) {
	createUser(username, password, alias)
}

func (RepoOperator) FindUser(userID uint) (*entity.User, error) {
	return findUser(userID)
}

func (RepoOperator) CreateRoom(participant uint) (*entity.Room, error) {
	return createRoom(participant)
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

func createUser(username, password, alias string) {
	user := entity.User{Username: username, Password: password, Alias: alias}
	NewRepository().Create(&user)
}

func findUser(userID uint) (*entity.User, error) {
	var user entity.User
	result := NewRepository().Select("id", "username", "alias").Find(&user, "id =?", userID)
	if result.RowsAffected != 1 {
		return nil, errors.New("can't find user with this id")
	}
	return &user, nil
}

func createRoom(participant uint) (*entity.Room, error) {
	var users []*entity.User
	user, err := findUser(participant)
	if err != nil {
		return nil, errors.New("can't find user")
	}
	users = append(users, user)
	room := entity.Room{Participants: users}
	result := NewRepository().Create(&room)
	if result.RowsAffected != 1 {
		return nil, result.Error
	}
	return &room, nil
}
