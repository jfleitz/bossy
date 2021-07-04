package main

import (
	"os"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

func init() {
	settings = loadConfiguration("config.json")
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.Infoln("bossy init complete")
}

var game goflip.GoFlip
var settings config

func main() {
	game.Observers = []goflip.Observer{
		new(bossyObserver),
		new(goalObserver),
		new(endOfBallBonus),
		new(warmUpPeriodObserver),
		//new(collectOvertime),
		//new(overTimeObserver),
	}

	game.DiagObserver = new(diagObserver)

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
		addThousands(5000)
	case swOuterLeftLane:
		addThousands(5000)
	case swMiddleLeftLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swInnerLeftLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swRightSlingshot:
		game.SolenoidFire(solRightSlingshot)
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
	case swLeftSlingshot:
		game.SolenoidFire(solLeftSlingshot)
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
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
		addHundreds(300)
		go saucerControl()
	case swLeftPointLane:
		game.AddScore(1000)
		game.PlaySound(sndFiring)
	case swLeftTarget:
		game.AddScore(300)
		game.PlaySound(sndTargets)
	case swLeftBumper:
		game.SolenoidOnDuration(solLeftBumper, 4)
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
	case swRightBumper:
		game.SolenoidOnDuration(solRightBumper, 4)
		game.AddScore(100)
		game.PlaySound(sndPuckBounce)
	case swBehindGoalLane:
		game.AddScore(1000)
	case swGoalie:
		//game.AddScore(1000) handled by shotObserver
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
		game.AddScore(500)
		game.PlaySound(sndRaRa)
	case swTargetO:
		game.AddScore(500)
		game.PlaySound(sndRaRa)
	case swTargetA:
		game.AddScore(500)
		game.PlaySound(sndRaRa)
	case swTargetL:
		game.AddScore(500)
		game.PlaySound(sndRaRa)
	}
}

func saucerControl() {
	go func() {
		game.PlaySound(sndRaRa)
		totalLetters := getPlayerStat(game.CurrentPlayer, bipShotCount)

		for i := 0; i < totalLetters; i++ {
			game.AddScore(1000)
			time.Sleep(500 * time.Millisecond)
		}

		if totalLetters == 0 {
			time.Sleep(1 * time.Second)
		}

		game.SolenoidFire(solSaucer)
	}()
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
	///	game.PlaySound(sndWhistle)
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
}

//Incremental scoring with sounds..
func addHundreds(points int) {
	go func() {
		for i := 100; i <= points; i += 100 {
			game.AddScore(i)
			game.PlaySound(snd100Points)
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func addThousands(points int) {
	go func() {
		for i := 1000; i <= points; i += 1000 {
			game.AddScore(i)
			game.PlaySound(snd1000Points)
			time.Sleep(500 * time.Millisecond)
		}
	}()
}
