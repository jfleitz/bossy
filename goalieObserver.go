/* Handles the Goalie and the 4 drop targets behind the goal
1. When a goal is scored, the Goal is worth 10,000 points
2. Plus 10,000 * each letter that the player has collected
3. Need to clear out the Letters for the Goal (but keep track of total letters
)

*/

package main

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type goalieObserver struct {
	//add your variables for the observer here
	//goalBonusLights []int
	movePosition int
	move         bool
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*goalieObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *goalieObserver) Init() {
	p.move = false
	p.movePosition = settings.Goalie.CenterPosition
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("goalieObserver:Init called")
	p.moveGoalie()

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *goalieObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	switch sw.SwitchID {
	case swLeftTarget:
		p.movePosition = settings.Goalie.TargetGLeft
		break
	case swOuterLeftLane:
		p.movePosition = settings.Goalie.TargetGLeft
		break
	case swMiddleLeftLane:
		p.movePosition = settings.Goalie.TargetGLeft
		break
	case swInnerLeftLane:
		p.movePosition = settings.Goalie.TargetGLeft
		break
	case swInnerRightLane:
		p.movePosition = settings.Goalie.TargetLRight
		break
	case swMiddleRightLane:
		p.movePosition = settings.Goalie.TargetLRight
		break
	case SwOuterRightLane:
		p.movePosition = settings.Goalie.TargetLRight
		break
	case swLowerRightTarget:
		p.movePosition = settings.Goalie.TargetLRight
		break
	case swMiddleRightTarget:
		p.movePosition = settings.Goalie.TargetLRight
		break
	case swUpperRightTarget:
		p.movePosition = settings.Goalie.TargetLRight
		break
	default:
		return
	}

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *goalieObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *goalieObserver) PlayerUp(playerID int) {
	p.move = true

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *goalieObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *goalieObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	p.move = false
}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *goalieObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *goalieObserver) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *goalieObserver) GameOver() {
	game.ServoAngle(settings.Goalie.StartPosition)

}

/*GameStart is called whenever a new game is started*/
func (p *goalieObserver) GameStart() {

}

func (p *goalieObserver) moveGoalie() {
	go func() {
		for {
			//startpos is either left or right
			game.ServoAngle(p.movePosition) //set to start pos

			for {
				if !p.move {
					break
				}
				//wait 2 seconds, then move to center
				time.Sleep(2 * time.Second)
				game.ServoAngle(settings.Goalie.CenterPosition)
				//wait 2 seconds from moving, then move back to startpos
				time.Sleep(2 * time.Second)
				game.ServoAngle(p.movePosition)

				//repeat
			}

			//JAF TODO we should sleep, or make this a channel actually...
			time.Sleep(250 * time.Millisecond)
		}
	}()
}
