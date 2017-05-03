package main

import "time"
import log "github.com/Sirupsen/logrus"

var inWarmUpPeriod bool
var cancelWarmUp bool

func startWarmUpPeriod() {
	go func() {
		inWarmUpPeriod = true
		cancelWarmUp = false
		defer func() {
			game.LampOff(lmpSamePlayerShootAgain)
			inWarmUpPeriod = false
			cancelWarmUp = false
			log.Infoln("Warmup Period complete")
		}()

		log.Infoln("Warmup Period beginning")
		game.LampOn(lmpSamePlayerShootAgain)
		sleepAndCheck(4)
		if cancelWarmUp {
			return
		}
		log.Infoln("Warmup Period 5 seconds left")
		game.LampSlowBlink(lmpSamePlayerShootAgain)
		sleepAndCheck(3)
		if cancelWarmUp {
			return
		}

		log.Infoln("Warmup Period 2 seconds left")
		game.LampFlastBlink(lmpSamePlayerShootAgain)
		sleepAndCheck(2)
		if cancelWarmUp {
			return
		}

	}()
}

func sleepAndCheck(ts int) {
	for i := 0; i < ts*2; i++ { //looping every half second to give a chance to cancel
		if cancelWarmUp {
			return
		}
		time.Sleep(time.Duration(ts) * time.Millisecond * 500)
	}
}
