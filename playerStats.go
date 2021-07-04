package main

import (
	"sync"
)

var _playerStats [5]playerStats

type playerStats struct {
	mu     sync.Mutex
	values map[string]int
}

func getPlayerStat(player int, key string) int {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	return _playerStats[player].values[key]
}

func setPlayerStat(player int, key string, value int) {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	_playerStats[player].values[key] = value
}

func incPlayerStat(player int, key string) {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	_playerStats[player].values[key]++
}

func decPlayerStat(player int, key string) {
	defer _playerStats[player].mu.Unlock()
	_playerStats[player].mu.Lock()
	if _playerStats[player].values[key] > 0 {
		_playerStats[player].values[key]--
	}
}

func initStats() {
	for i := range _playerStats {
		_playerStats[i].values = make(map[string]int)
	}
}
