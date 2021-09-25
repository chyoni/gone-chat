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
	FindUserByID(userID uint) (*entity.User, error)
	FindUserByUsername(username string) (*entity.User, error)
	GetUser(userID uint) (*entity.User, error)
	DeleteUser(userID uint) error
	UpdateUserAlias(userID uint, alias string) (*entity.User, error)
	UpdatePassword(userID uint, password string) error
	CheckUserPassword(userID uint, hashedPassword string) bool
	CreateRoom(participant uint) (*entity.Room, error)
}

type RepoOperator struct{}

func (RepoOperator) CreateUser(username, password, alias string) {
	createUser(username, password, alias)
}
func (RepoOperator) FindUserByID(userID uint) (*entity.User, error) {
	return findUserByID(userID)
}
func (RepoOperator) CreateRoom(participant uint) (*entity.Room, error) {
	return createRoom(participant)
}
func (RepoOperator) FindUserByUsername(username string) (*entity.User, error) {
	return findUserByUsername(username)
}
func (RepoOperator) GetUser(userID uint) (*entity.User, error) {
	return getUser(userID)
}
func (RepoOperator) UpdateUserAlias(userID uint, alias string) (*entity.User, error) {
	return updateUserAlias(userID, alias)
}
func (RepoOperator) UpdatePassword(userID uint, password string) error {
	return updatePassword(userID, password)
}
func (RepoOperator) CheckUserPassword(userID uint, hashedPassword string) bool {
	return checkUserPassword(userID, hashedPassword)
}
func (RepoOperator) DeleteUser(userID uint) error {
	return deleteUser(userID)
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

func findUserByID(userID uint) (*entity.User, error) {
	var user entity.User
	result := NewRepository().Select("id", "username", "alias", "avatar").Find(&user, "id =?", userID)
	if result.RowsAffected != 1 {
		return nil, errors.New("can't find user with this id")
	}
	return &user, nil
}

func findUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	result := NewRepository().Select("id", "username", "password").Find(&user, "username = ?", username)
	if result.RowsAffected != 1 {
		return nil, errors.New("can't find user with this username")
	}
	return &user, nil
}

func checkUserPassword(userID uint, hashedPassword string) bool {
	var user entity.User
	result := NewRepository().Select("password").Find(&user, "id = ?", userID)
	if result.RowsAffected != 1 {
		return false
	}
	if user.Password != hashedPassword {
		return false
	}
	return true
}

func updateUserAlias(userID uint, alias string) (*entity.User, error) {
	result := NewRepository().Model(&entity.User{}).Where("id = ?", userID).Update("alias", alias)
	if result.RowsAffected != 1 {
		return nil, errors.New("can't find user with this id")
	}
	updatedUser, err := findUserByID(userID)
	if err != nil {
		return nil, errors.New("something wrong in database")
	}
	return updatedUser, nil
}

func updatePassword(userID uint, hashedPassword string) error {
	var user entity.User
	result := NewRepository().Select("password").Find(&user, "id = ?", userID).Update("password", hashedPassword)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func getUser(userID uint) (*entity.User, error) {
	var user entity.User
	result := NewRepository().Select("id", "username", "alias", "avatar", "created_at", "updated_at").Find(&user, "id = ?", userID)
	if result.RowsAffected != 1 {
		return nil, errors.New("can't find user with this id")
	}
	return &user, nil
}

func deleteUser(userID uint) error {
	result := NewRepository().Delete(&entity.User{}, userID)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func createRoom(participant uint) (*entity.Room, error) {
	var users []*entity.User
	user, err := findUserByID(participant)
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
