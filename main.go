package main

import (
	"github.com/chiwon99881/gone-chat/api"
	"github.com/chiwon99881/gone-chat/env"
	"github.com/chiwon99881/gone-chat/explorer"
)

func main() {
	env.Start()
	go api.Start()
	explorer.Start()
}
