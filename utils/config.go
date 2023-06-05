package utils

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	TotalBalls              int
	MaxPlayers              int
	TimedMode               bool
	BossyBonusFromGoal      bool
	BallTimeSeconds         int
	WarmupPeriodTimeSeconds int
	LogLevel                string
	KeepAliveMS             int
	ConsoleMode             bool //For testing

	Goalie GoalieConf
}

var conf *Config

type GoalieConf struct {
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
	DebugGoalie    bool    `toml:"debug_goalie"`
}

func LoadConfiguration(file string) error {
	conf = new(Config)

	if _, err := toml.DecodeFile(file, &conf); err != nil {
		fmt.Printf("Could not load conf file %v\n", file)
		fmt.Println(err)

		return err
	}

	return nil
}

func Settings() Config {
	//todo JAF lock and initialize this..
	return *conf
}
