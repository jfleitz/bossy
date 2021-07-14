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

type goalObserver struct {
	//add your variables for the observer here
	goalBonusLights []int
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*goalObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *goalObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("goalObserver:Init called")

	p.goalBonusLights = []int{
		lmp5000Bonus1,
		lmp5000Bonus2,
		lmp5000Bonus3,
		lmp5000Bonus4,
	}
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *goalObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	switch sw.SwitchID {
	case swTargetG:
		break
	case swTargetO:
		break
	case swTargetA:
		break
	case swTargetL:
		break
	case swBackGoal:
		break
	default:
		return
	}

	game.AddScore(500)
	game.PlaySound(sndGoal)
	incPlayerStat(game.CurrentPlayer, totalGoalCount)
	allTargets := game.SwitchPressed(swTargetG) && game.SwitchPressed((swTargetO)) &&
		game.SwitchPressed(swTargetA) && game.SwitchPressed(swTargetL)
	//keep track of the G O A L targets, and reset the bank afterwards. Also light the 5k bonus

	//flash the goal light and reset target bank
	go func() {
		if allTargets {
			game.SolenoidFire(solDropTargets)
		}
		game.LampFastBlink(lmpGoalLight)
		time.Sleep(3 * time.Second)
		game.LampOff(lmpGoalLight)
	}()

	if allTargets {

		cnt := incPlayerStat(game.CurrentPlayer, goalTargetCount)
		log.Infof("All targets down. goal count is now %d\n", cnt)
		if cnt >= 5 {
			game.LampOn(lmp25000Bonus)
			cnt -= 5
		}

		if cnt < 5 {
			game.LampOn(p.goalBonusLights[cnt-1])
		}
	}

	//send back a command over the switch handler channel to call on choosepuck, hat trick determination, etc
	//game.BroadcastEvent(goflip.SwitchEvent{SwitchID: goalScored, Pressed: true})
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *goalObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *goalObserver) PlayerUp(playerID int) {
	game.SolenoidFire(solDropTargets)
	game.LampOff(p.goalBonusLights...)
	game.LampOff(lmp25000Bonus)
	setPlayerStat(game.CurrentPlayer, goalTargetCount, 0)
	setPlayerStat(game.CurrentPlayer, totalGoalCount, 0)
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *goalObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *goalObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *goalObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *goalObserver) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *goalObserver) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *goalObserver) GameStart() {

}
