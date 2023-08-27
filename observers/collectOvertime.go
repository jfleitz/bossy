/*Overtime (extra ball).

We need to determine when Extra Ball is lit. This is for overtime, which should be a tied game?


On the last ball, the overtime lights are lit.
For each time you hit one, 1 second is added to the overtime.
If there are more than one player playing, whomever has the
higher overtime value gets the timed extra ball.*/

package observer //this will probably be package main in your app

import (
	"sync"
	"time"

	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type CollectOvertime struct {
	//add your variables for the observer here
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*CollectOvertime)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *CollectOvertime) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Traceln("CollectOvertime:Init called")

}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *CollectOvertime) SwitchHandler(sw goflip.SwitchEvent) {
	game := goflip.GetMachine()
	if !sw.Pressed && game.BallInPlay != 3 { //always just a 3 ball game, and we just care about the 3rd ball for incrementing over time.
		return
	}

	switch sw.SwitchID {
	case SwTopLeftLane:
		p.incAndFlashLamp(LmpTopOrangeSpot)
	case SwTopMiddleLane:
		p.incAndFlashLamp(LmpTopOrangeSpot)
	case SwTopRightLane:
		p.incAndFlashLamp(LmpTopOrangeSpot)
	case SwUpperRightTarget:
		p.incAndFlashLamp(LmpTargetsOrangeSpot)
	case SwMiddleRightTarget:
		p.incAndFlashLamp(LmpTargetsOrangeSpot)
	case SwLowerRightTarget:
		p.incAndFlashLamp(LmpTargetsOrangeSpot)
	case SwInnerRightLane:
		p.incAndFlashLamp(LmpRightReturnLaneOrangeSpot)
	case SwInnerLeftLane:
		p.incAndFlashLamp(LmpLeftReturnLaneOrangeSpot)
	case SwBehindGoalLane:
		p.incAndFlashLamp(LmpGoalOnLeftOrangeSpot)
	default:
		return
	}

}

func (p *CollectOvertime) incAndFlashLamp(lmpID int) {
	go func() {
		game := goflip.GetMachine()
		utils.IncPlayerStat(game.CurrentPlayer, OTSeconds)
		goflip.LampFastBlink(lmpID)
		time.After(1 * time.Second)
		p.OTLightsEnable(true)
	}()
}

func (p *CollectOvertime) OTLightsEnable(on bool) {
	if on {
		goflip.LampOn(
			LmpOvertimeLeftOfGoal,
			LmpTopOrangeSpot,
			LmpTargetsOrangeSpot,
			LmpRightReturnLaneOrangeSpot,
			LmpLeftReturnLaneOrangeSpot,
			LmpGoalOnLeftOrangeSpot)
	} else {
		goflip.LampOff(
			LmpOvertimeLeftOfGoal,
			LmpTopOrangeSpot,
			LmpTargetsOrangeSpot,
			LmpRightReturnLaneOrangeSpot,
			LmpLeftReturnLaneOrangeSpot,
			LmpGoalOnLeftOrangeSpot)
	}
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *CollectOvertime) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *CollectOvertime) PlayerUp(playerID int) {
	game := goflip.GetMachine()
	if game.BallInPlay == 3 {
		p.OTLightsEnable(true)
	} else {
		p.OTLightsEnable((false))
	}
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *CollectOvertime) PlayerStart(playerID int) {

}

/*PlayerEnd is called after every ball for the player is over*/
func (p *CollectOvertime) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	p.OTLightsEnable((false))
	defer wait.Done()

}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *CollectOvertime) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *CollectOvertime) PlayerAdded(playerID int) {

}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *CollectOvertime) GameOver() {
	p.OTLightsEnable((false))
}

/*GameStart is called whenever a new game is started*/
func (p *CollectOvertime) GameStart() {

}
