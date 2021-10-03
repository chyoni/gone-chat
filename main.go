package main

import (
	"github.com/chiwon99881/gone-chat/api"
	"github.com/chiwon99881/gone-chat/auth"
	"github.com/chiwon99881/gone-chat/database"
	"github.com/chiwon99881/gone-chat/env"
	"github.com/chiwon99881/gone-chat/ws"
)

func main() {
	defer database.Close()
	env.Start()
	database.NewRepository()
	auth.Start()
	go ws.Start()
	api.Start()
}
