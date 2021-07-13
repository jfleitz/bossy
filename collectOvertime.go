/*Overtime awards. On the last ball, the overtime lights are lit.
For each time you hit one, 1 second is added to the overtime.
If there are more than one player playing, whomever has the
higher overtime value gets the timed extra ball.*/

package main //this will probably be package main in your app

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type collectOvertime struct {
	//add your variables for the observer here
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*collectOvertime)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *collectOvertime) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("collectOvertime:Init called")

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *collectOvertime) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed && game.BallInPlay != 3 { //always just a 3 ball game, and we just care about the 3rd ball for incrementing over time.
		return
	}

	switch sw.SwitchID {
	case swTopLeftLane:
		p.incAndFlashLamp(lmpTopOrangeSpot)
	case swTopMiddleLane:
		p.incAndFlashLamp(lmpTopOrangeSpot)
	case swTopRightLane:
		p.incAndFlashLamp(lmpTopOrangeSpot)
	case swUpperRightTarget:
		p.incAndFlashLamp(lmpTargetsOrangeSpot)
	case swMiddleRightTarget:
		p.incAndFlashLamp(lmpTargetsOrangeSpot)
	case swLowerRightTarget:
		p.incAndFlashLamp(lmpTargetsOrangeSpot)
	case swInnerRightLane:
		p.incAndFlashLamp(lmpRightReturnLaneOrangeSpot)
	case swInnerLeftLane:
		p.incAndFlashLamp(lmpLeftReturnLaneOrangeSpot)
	case swBehindGoalLane:
		p.incAndFlashLamp(lmpGoalOnLeftOrangeSpot)
	default:
		return
	}

}

func (p *collectOvertime) incAndFlashLamp(lmpID int) {
	go func() {
		incPlayerStat(game.CurrentPlayer, otSeconds)
		game.LampFastBlink(lmpID)
		time.After(1 * time.Second)
		p.flashOTLights()
	}()
}

func (p *collectOvertime) flashOTLights() {
	game.LampSlowBlink(
		lmpOvertimeLeftOfGoal,
		lmpTopOrangeSpot,
		lmpTargetsOrangeSpot,
		lmpRightReturnLaneOrangeSpot,
		lmpLeftReturnLaneOrangeSpot,
		lmpGoalOnLeftOrangeSpot)

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *collectOvertime) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *collectOvertime) PlayerUp(playerID int) {
	if game.BallInPlay == 3 {
		p.flashOTLights()
	}
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *collectOvertime) PlayerStart(playerID int) {

}

/*PlayerEnd is called after every ball for the player is over*/
func (p *collectOvertime) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *collectOvertime) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *collectOvertime) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *collectOvertime) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *collectOvertime) GameStart() {

}
