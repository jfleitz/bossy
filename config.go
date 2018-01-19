package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	TotalBalls              int  `json:"totalballs"`
	MaxPlayers              int  `json:"maxplayers"`
	TimedMode               bool `json:"timedmode"`
	BallTimeSeconds         int  `json:"balltimeseconds"`
	WarmupPeriodTimeSeconds int  `json:"warmupperiodtimeseconds"`
}

func loadConfiguration(file string) config {
	var c config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	return c
}
