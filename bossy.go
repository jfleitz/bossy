package main

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

	game.Init()
	game.TotalBalls = 3
	game.Start(switchHandler)

	for {
		time.Sleep(1000 * time.Millisecond) //just keep sleeping
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
		go ballOverControl()
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
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)

}

func ballOverControl() {
	time.Sleep(1 * time.Second)
	game.SolenoidFire(solOuthole)
}
