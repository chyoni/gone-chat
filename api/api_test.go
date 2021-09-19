package api

import (
	"fmt"
	"testing"
)

type fakeRepository struct{}

func (fakeRepository) CreateUser(username, password string) {
	fmt.Println("mocked")
}

func TestCreateUser(t *testing.T) {
	dbOperator = &fakeRepository{}
}
