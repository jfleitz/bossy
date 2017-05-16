package main

/*
TODO/Notes:
Warm Up period..when this is going, all puck lights should be on (flashing).
You collect each (which goes to fast flash, then back to slow flash)
...so WarmUp period should control anything puck related (not have puckChase and warmUp communicate).
This will make things neater/more maintainable in code.



*/

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

func init() {
	log.SetOutput(os.Stdout)

}

var game goflip.GoFlip

func main() {
	//	var p *puckChase
	//	var g *goalObserver
	//	p = new(puckChase)
	//	g = new(goalObserver)
	game.Observers = []goflip.Observer{
		new(puckChase),
		new(goalObserver),
		new(endOfBallBonus),
		new(hatTrick),
		new(lilcoLine),
		new(collectOvertime),
		new(overTimeObserver),
	}

	inWarmUpPeriod = false
	game.TotalBalls = 3
	game.BallInPlay = 0 //using this for GameOver for now
	game.MaxPlayers = 0
	initStats()
	game.Init(switchHandler)

	for {
		time.Sleep(1000 * time.Millisecond) //just keep sleeping
		game.SendStats()

		log.Infoln("Still Looping ", time.Now)
	}
}

func switchHandler(sw goflip.SwitchEvent) {

	if !sw.Pressed {
		return
	}

	log.Infof("Bossy switchHandler. Receivied SwitchID=%d Pressed=%v\n", sw.SwitchID, sw.Pressed)

	switch sw.SwitchID {
	case swOuthole:
		//ball over
		log.Infoln("outhole: switch pressed")
		if game.BallInPlay > 0 {
			log.Infoln("outhole: ballinplay is", game.BallInPlay)
			game.BallDrained()
			log.Infoln("outhole: ball drained called")
			if !inWarmUpPeriod {
				log.Infoln("outhole: not in warm up period")
				game.PlayerEnd()  //Calling this since we don't have a ball save?
				ballOverControl() //maybe have this controlled by player end instead? register a method?
			} else {
				//go ahead and eject it again
				game.SolenoidFire(solOuthole)
			}

		}

	case swCredit:
		//start game..make this more elegant
		go creditControl()
	case swShooterLane:
	case swTest:
	case swCoin:
	case swInnerRightLane:
	case swMiddleRightLane:
	case SwOuterRightLane:
	case swOuterLeftLane:
	case swMiddleLeftLane:
	case swInnerLeftLane:
	case swRightSlingshot:
		game.SolenoidFire(solRightKicker)
	case swLeftSlingshot:
	//which one is this??
	case swLowerRightTarget:
	case swMiddleRightTarget:
	case swUpperRightTarget:
	case swSaucer:
		go saucerControl()
	case swLeftPointLane:
	case swLeftTarget:
	case swLeftBumper:
		game.SolenoidOnDuration(solLeftBumper, 4)
	case swRightBumper:
		game.SolenoidOnDuration(solRightBumper, 4)
	case swBehindGoalLane:
	case swGoalie:
	case swTopLeftLane:
		game.LampOn(lmpTopLeftLane)
	case swTopMiddleLane:
		game.LampOn(lmpTopMiddleLane)
	case swTopRightLane:
		game.LampOn(lmpTopRightLane)
	case swTargetG:
	case swTargetO:
	case swTargetA:
	case swTargetL:

	}
}

func saucerControl() {
	time.Sleep(2 * time.Second)
	game.SolenoidFire(solSaucer)
}

func creditControl() {
	game.FlipperControl(true)
	game.MaxPlayers++ //add a player to the list

	if !game.IsGameInPlay() {
		game.Scores = make([]int32, game.MaxPlayers)
		game.NextUp() //to make it player 1
		startWarmUpPeriod()
	}
	//	log.Infoln("checking outhole switch pressed")
	//	if game.SwitchPressed(swOuthole) { //JAF BUG: if the game starts with the switch down, it is not detected
	//		log.Infoln("outhole switch is pressed")
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
	//	} else {
	//		log.Infoln("outhole switch is not pressed")
	//	}
}

func ballOverControl() {

	game.NextUp()
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
}
