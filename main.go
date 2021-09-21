package main

import (
	"github.com/chiwon99881/gone-chat/api"
	"github.com/chiwon99881/gone-chat/auth"
	"github.com/chiwon99881/gone-chat/database"
	"github.com/chiwon99881/gone-chat/env"
	"github.com/chiwon99881/gone-chat/explorer"
)

func main() {
	env.Start()
	database.NewRepository()
	auth.Start()
	go api.Start()
	explorer.Start()
}
