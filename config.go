package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type config struct {
	TotalBalls              int
	MaxPlayers              int
	TimedMode               bool
	BossyBonusFromGoal      bool
	BallTimeSeconds         int
	WarmupPeriodTimeSeconds int
	LogLevel                log.Level
	KeepAliveMS             int //time in milliseconds on whow to send a msg to each arduino
}

func loadConfiguration(file string) config {
	var c config

	if _, err := toml.DecodeFile(file, &c); err != nil {
		fmt.Printf("Could not load conf file %v\n", file)
		fmt.Println(err)
		return c
	}

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Errorln(err.Error())
		//fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	log.Infoln("Loaded config.json")
	return c
}
