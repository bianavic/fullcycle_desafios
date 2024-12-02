package main

import (
	"github.com/bianavic/fullcycle_desafios.git/config"
)

var (
	logger config.Logger
)

func main() {

	logger = *config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("config initialization error %v ", err)
		return
	}
	//router.Initialize()
}
