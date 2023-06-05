package utils

import (
	"sync"
)

var _playerStats [5]playerStats

type playerStats struct {
	mu     sync.Mutex
	values map[string]int
}

func GetPlayerStat(player int, key string) int {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	return _playerStats[player].values[key]
}

func SetPlayerStat(player int, key string, value int) {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	_playerStats[player].values[key] = value
}

func IncPlayerStat(player int, key string) int {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	_playerStats[player].values[key]++
	return _playerStats[player].values[key]
}

func DecPlayerStat(player int, key string) {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	if _playerStats[player].values[key] > 0 {
		_playerStats[player].values[key]--
	}
}

func InitStats() {
	for i := range _playerStats {
		_playerStats[i].values = make(map[string]int)
	}
}
