package main

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type config struct {
	TotalBalls              int       `json:"totalballs"`
	MaxPlayers              int       `json:"maxplayers"`
	TimedMode               bool      `json:"timedmode"`
	BallTimeSeconds         int       `json:"balltimeseconds"`
	WarmupPeriodTimeSeconds int       `json:"warmupperiodtimeseconds"`
	LogLevel                log.Level `json:"loglevel"`
}

func loadConfiguration(file string) config {
	var c config
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
