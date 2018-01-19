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
var settings config

func main() {
	//	var p *puckChase
	//	var g *goalObserver
	//	p = new(puckChase)
	//	g = new(goalObserver)
	game.Observers = []goflip.Observer{
		new(bossyObserver),
		//new(puckChase),
		//new(goalObserver),
		//new(endOfBallBonus),
		//new(hatTrick),
		//new(lilcoLine),
		//new(collectOvertime),
		//new(overTimeObserver),
		new(warmUpPeriodObserver),
	}

	game.DiagObserver = new(diagObserver)

	settings = loadConfiguration("config.json")

	inWarmUpPeriod = false

	//Set the game limitations here
	game.TotalBalls = settings.TotalBalls
	game.MaxPlayers = settings.MaxPlayers

	//set the game initalizations here
	game.BallInPlay = 0
	game.CurrentPlayer = 0
	initStats()
	game.Init(switchHandler)

	for {
		time.Sleep(1000 * time.Millisecond) //just keep sleeping
		game.SendStats()
		//log.Infoln("Still Looping ", time.Now)
	}
}

func switchHandler(sw goflip.SwitchEvent) {

	if !sw.Pressed {
		return
	}

	log.Infof("Bossy switchHandler. Receivied SwitchID=%d Pressed=%v\n", sw.SwitchID, sw.Pressed)

	if !game.GameRunning {
		//only care about switches that matter when a game is not running
		switch sw.SwitchID {
		case swSaucer:
			game.SolenoidFire(solSaucer)
		case swCredit:
			//start game
			go creditControl()
		}

		return
	}

	switch sw.SwitchID {
	case swOuthole:
		//ball over
		log.Infoln("outhole: switch pressed")
		game.BallDrained()

	case swCredit:
		go creditControl()
	case swShooterLane:
		game.PlaySound(sndBallLaunch)
	case swTest:
	case swCoin:
	case swInnerRightLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swMiddleRightLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case SwOuterRightLane:
		game.AddScore(5000)
		game.PlaySound(sndFiring)
	case swOuterLeftLane:
		game.AddScore(5000)
		game.PlaySound(sndOutlane)
	case swMiddleLeftLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swInnerLeftLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swRightSlingshot:
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
		game.SolenoidFire(solRightSlingshot)
	case swLeftSlingshot:
		//which one is this??
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
		game.SolenoidFire(solLeftSlingshot)

	case swLowerRightTarget:
		game.AddScore(1000)
		game.PlaySound(sndTargets)
	case swMiddleRightTarget:
		game.AddScore(300)
		game.PlaySound(sndTargets)
	case swUpperRightTarget:
		game.AddScore(1000)
		game.PlaySound(sndTargets)
	case swSaucer:
		game.AddScore(300)
		go saucerControl()
	case swLeftPointLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swLeftTarget:
		game.AddScore(300)
		game.PlaySound(sndTargets)
	case swLeftBumper:
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
		game.SolenoidOnDuration(solLeftBumper, 4)
	case swRightBumper:
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
		game.SolenoidOnDuration(solRightBumper, 4)
	case swBehindGoalLane:
		game.AddScore(1000)
	case swGoalie:
		//game.AddScore(1000) handled by goalObserver
	case swTopLeftLane:
		//game.LampOn(lmpTopLeftLane)
		game.AddScore(300)
		game.PlaySound(sndTargets)
	case swTopMiddleLane:
		game.AddScore(300)
		//game.LampOn(lmpTopMiddleLane)
		game.PlaySound(sndTargets)
	case swTopRightLane:
		game.AddScore(300)
		//game.LampOn(lmpTopRightLane)
		game.PlaySound(sndTargets)
	case swTargetG:
		game.AddScore(1000)
		game.PlaySound(sndRaRa)
	case swTargetO:
		game.AddScore(1000)
		game.PlaySound(sndRaRa)
	case swTargetA:
		game.AddScore(1000)
		game.PlaySound(sndRaRa)
	case swTargetL:
		game.AddScore(1000)
		game.PlaySound(sndRaRa)

	}
}

func saucerControl() {
	time.Sleep(2 * time.Second)
	game.SolenoidFire(solSaucer)
}

func creditControl() {
	if !game.GameRunning {
		game.GameStart()
		game.AddPlayer() // go ahead and add player 1
		game.PlayerUp()
	} else {
		if game.BallInPlay == 1 {
			game.AddPlayer()
		}
	}
}

func ballLaunch() {
	//	game.NextUp()
	game.PlaySound(sndWhistle)
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
}
