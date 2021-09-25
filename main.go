package main

import (
	"github.com/chiwon99881/gone-chat/api"
	"github.com/chiwon99881/gone-chat/auth"
	"github.com/chiwon99881/gone-chat/database"
	"github.com/chiwon99881/gone-chat/env"
)

func main() {
	env.Start()
	database.NewRepository()
	auth.Start()
	api.Start()
}
