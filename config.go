package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type config struct {
	TotalBalls              int
	MaxPlayers              int
	TimedMode               bool
	BossyBonusFromGoal      bool
	BallTimeSeconds         int
	WarmupPeriodTimeSeconds int
	LogLevel                string
	KeepAliveMS             int
}

func loadConfiguration(file string) config {
	var c config

	if _, err := toml.DecodeFile(file, &c); err != nil {
		fmt.Printf("Could not load conf file %v\n", file)
		fmt.Println(err)

		return c
	}

	return c
}
