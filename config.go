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

	Goalie goalieConf
}

type goalieConf struct {
	StartPosition  int     `toml:"start_position"`
	CenterPosition int     `toml:"center_position"`
	LimitLeft      int     `toml:"limit_left"`
	LimitRight     int     `toml:"limit_right"`
	TargetGLeft    int     `toml:"target_g_left"`
	TargetGRight   int     `toml:"target_g_right"`
	TargetOLeft    int     `toml:"target_o_left"`
	TargetORight   int     `toml:"target_o_right"`
	TargetALeft    int     `toml:"target_a_left"`
	TargetARight   int     `toml:"target_a_right"`
	TargetLLeft    int     `toml:"target_l_left"`
	TargetLRight   int     `toml:"target_l_right"`
	DeviceAddress  string  `toml:"device_address"`
	PulseMin       float32 `toml:"pulse_min"`
	PulseMax       float32 `toml:"pulse_max"`
	ArcRange       int     `toml:"arc_range"`
}

func loadConfiguration(file string) (*config, error) {
	c := new(config)

	if _, err := toml.DecodeFile(file, &c); err != nil {
		fmt.Printf("Could not load conf file %v\n", file)
		fmt.Println(err)

		return c, err
	}

	fmt.Printf("Config example: %v\n", c.KeepAliveMS)
	fmt.Printf("Config c.Goalie.startPosition: %v\n", c.Goalie.StartPosition)

	return c, nil
}
