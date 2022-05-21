package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

var game goflip.GoFlip
var settings *config

func init() {
	var err error
	settings, err = loadConfiguration("config.toml")

	if err != nil {
		os.Exit(1)
	}

	if lvl, err := log.ParseLevel(settings.LogLevel); err != nil {
		fmt.Errorf("Error with log level in config: %v", err)
	} else {
		log.SetLevel(lvl)
	}

	log.SetOutput(os.Stdout)
	log.Debugln("bossy init complete")
	log.Printf("TotalBalls: %v\n", settings.TotalBalls)
	log.Printf("BallTimeSeconds: %v\n", settings.BallTimeSeconds)
	log.Printf("Start position: %v\n", settings.Goalie.StartPosition)
	log.Printf("limitLeft position: %v\n", settings.Goalie.LimitLeft)

}

func main() {
	game.Observers = []goflip.Observer{
		new(bossyObserver),
		new(goalObserver),
		//new(warmUpPeriodObserver),
		new(shotObserver),
		new(endOfBallBonus),
		new(attractMode),
		//new(collectOvertime),
		//new(overTimeObserver),
		new(goalieObserver),
	}

	game.DiagObserver = new(diagObserver)

	inWarmUpPeriod = false

	//Set the game limitations here
	game.TotalBalls = settings.TotalBalls
	game.MaxPlayers = settings.MaxPlayers
	game.Credits = 0

	//Set Goalie Servo control params
	game.PWMPortConfig.ArcRange = settings.Goalie.ArcRange
	game.PWMPortConfig.PulseMax = settings.Goalie.PulseMax
	game.PWMPortConfig.PulseMin = settings.Goalie.PulseMin
	game.PWMPortConfig.DeviceAddress = settings.Goalie.DeviceAddress

	log.Debugf("Main, warmupseconds is: %v\n", settings.WarmupPeriodTimeSeconds)

	//set the game initalizations here
	game.BallInPlay = 0
	game.CurrentPlayer = 0
	initStats()
	game.Init(switchHandler)

	//go ahead and go to GameOver by default
	game.ChangeGameState(goflip.GameOver)

	for {
		time.Sleep(1000 * time.Millisecond) //just keep sleeping
		game.SendStats()
	}
}

func switchHandler(sw goflip.SwitchEvent) {

	if !sw.Pressed {
		return
	}

	log.Debugf("Bossy switchHandler. Receivied SwitchID=%d Pressed=%v\n", sw.SwitchID, sw.Pressed)

	if game.GetGameState() == goflip.GameOver {
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
		log.Debugln("outhole: switch pressed")
		game.BallDrained()

	case swCredit:
		go creditControl()
	case swShooterLane:
		//		game.PlaySound(sndShooter) //play elsewhere
	case swTest:
	case swCoin:
		game.Credits++
		game.SetCreditDisp(int8(game.Credits))
	case swInnerRightLane:
		addThousands(1000)
	case swMiddleRightLane:
		addThousands(1000)
	case SwOuterRightLane:
		addThousands(5000)
	case swOuterLeftLane:
		addThousands(5000)
	case swMiddleLeftLane:
		addThousands(1000)
	case swInnerLeftLane:
		addThousands(1000)
	case swRightSlingshot:
		game.SolenoidFire(solRightSlingshot)
		game.AddScore(100)
		game.PlaySound(sndSlingshot)
	case swLeftSlingshot:
		game.SolenoidFire(solLeftSlingshot)
		game.AddScore(100)
		game.PlaySound(sndSlingshot)
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
		addThousands(1000)
	case swLeftTarget:
		game.AddScore(300)
		game.PlaySound(sndTargets)
	case swLeftBumper:
		game.SolenoidOnDuration(solLeftBumper, 4)
		game.AddScore(100)
		game.PlaySound(sndBumper)
	case swRightBumper:
		game.SolenoidOnDuration(solRightBumper, 4)
		game.AddScore(100)
		game.PlaySound(sndBumper)
	case swBehindGoalLane:
		addThousands(1000)
	case swGoalie:
		//game.AddScore(1000) handled by shotObserver
	case swTopLeftLane:
		//game.LampOn(lmpTopLeftLane)
		game.AddScore(300)
		game.PlaySound(sndTopLane)
	case swTopMiddleLane:
		game.AddScore(300)
		//game.LampOn(lmpTopMiddleLane)
		game.PlaySound(sndTopLane)
	case swTopRightLane:
		game.AddScore(300)
		//game.LampOn(lmpTopRightLane)
		game.PlaySound(sndTopLane)
	case swTargetG:
	case swTargetO:
	case swTargetA:
	case swTargetL:
	}
}

func saucerControl() {
	//JAF TODO: If BossyBonusFromGoal is set, then start the timer to hit a goal and then collect the bonus
	go func() {
		game.PlaySound(sndSaucer)
		time.Sleep(2 * time.Second)
		totalLetters := getPlayerStat(game.CurrentPlayer, bipShotCount)

		addThousands(totalLetters * 1000)

		if totalLetters == 0 {
			time.Sleep(1 * time.Second)
		}

		game.SolenoidFire(solSaucer)
	}()
}

func creditControl() {
	if game.GetGameState() == goflip.GameOver {
		game.ChangeGameState(goflip.GameStart)
		game.AddPlayer() // go ahead and add player 1
		game.ChangePlayerState(goflip.PlayerUp)
	} else {
		if game.BallInPlay == 1 {
			game.AddPlayer()
		}
	}
}

func ballLaunch() {
	//	game.NextUp()
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
}

//Incremental scoring with sounds..
func addHundreds(points int) {
	go func() {
		for i := 1; i <= points/100; i++ {
			game.AddScore(100)
			game.PlaySound(snd100Points)
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func addThousands(points int) {
	go func() {
		for i := 1; i <= points/1000; i++ {
			game.AddScore(1000)
			game.PlaySound(snd1000Points)
			time.Sleep(250 * time.Millisecond)
		}
	}()
}
