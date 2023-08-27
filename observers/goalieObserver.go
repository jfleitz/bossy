/* Handles the Goalie and the 4 drop targets behind the goal
1. When a goal is scored, the Goal is worth 10,000 points
2. Plus 10,000 * each letter that the player has collected
3. Need to clear out the Letters for the Goal (but keep track of total letters
)

*/

package observer

import (
	"sync"
	"time"

	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type GoalieObserver struct {
	//add your variables for the observer here
	//goalBonusLights []int
	movePosition int
	move         bool
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*GoalieObserver)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *GoalieObserver) Init() {
	p.move = false
	p.movePosition = utils.Settings().Goalie.CenterPosition
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Traceln("GoalieObserver:Init called")
	p.moveGoalie()
}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *GoalieObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	switch sw.SwitchID {
	case SwTopLeftLane:
		fallthrough
	case SwTopMiddleLane:
		fallthrough
	case SwTopRightLane:
		p.movePosition = utils.Settings().Goalie.RightPosition
		p.move = true
	case SwLeftTarget:
		fallthrough
	case SwOuterLeftLane:
		fallthrough
	case SwMiddleLeftLane:
		fallthrough
	case SwInnerLeftLane:
		log.Traceln("GoalieObserver: Moving Left")
		p.movePosition = utils.Settings().Goalie.LeftPosition
		p.move = true
	case SwInnerRightLane:
		fallthrough
	case SwMiddleRightLane:
		fallthrough
	case SwOuterRightLane:
		fallthrough
	case SwLowerRightTarget:
		fallthrough
	case SwMiddleRightTarget:
		fallthrough
	case SwUpperRightTarget:
		log.Traceln("GoalieObserver: Moving Right")
		p.movePosition = utils.Settings().Goalie.RightPosition
		p.move = true
	default:
		return
	}

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *GoalieObserver) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *GoalieObserver) PlayerUp(playerID int) {
	p.move = true
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *GoalieObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *GoalieObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	p.movePosition = utils.Settings().Goalie.CenterPosition
	p.move = false
}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *GoalieObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *GoalieObserver) PlayerAdded(playerID int) {

}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *GoalieObserver) GameOver() {
	p.movePosition = utils.Settings().Goalie.CenterPosition
	p.move = false
}

/*GameStart is called whenever a new game is started*/
func (p *GoalieObserver) GameStart() {

}

func (p *GoalieObserver) moveGoalie() {
	go func() {
		for {
			//startpos is either left or right
			goflip.ServoAngle(p.movePosition) //set to start pos

			for {
				game := goflip.GetMachine()
				if game.Quitting {
					return
				}
				if !p.move {
					break
				}
				//wait 2 seconds, then move to center
				time.Sleep(2 * time.Second)
				goflip.ServoAngle(utils.Settings().Goalie.CenterPosition)
				//wait 2 seconds from moving, then move back to startpos
				time.Sleep(2 * time.Second)
				goflip.ServoAngle(p.movePosition)
			}

			time.Sleep(250 * time.Millisecond)
		}
	}()
}
