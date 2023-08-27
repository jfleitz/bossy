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

type GoalObserver struct {
	//add your variables for the observer here
	goalBonusLights []int
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*GoalObserver)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *GoalObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Traceln("GoalObserver:Init called")
}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *GoalObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	if goflip.GetGameState() != goflip.InProgress {
		return
	}

	switch sw.SwitchID {
	case SwTargetG:
		goflip.LampOn(Lmp5000G)
	case SwTargetO:
		goflip.LampOn(Lmp5000O)
	case SwTargetA:
		goflip.LampOn(Lmp5000A)
	case SwTargetL:
		goflip.LampOn(Lmp5000L)
	case SwBackGoal:
	default:
		return
	}

	game := goflip.GetMachine()

	goflip.AddScore(500)
	goflip.PlaySound(SndGoalTarget)
	utils.IncPlayerStat(game.CurrentPlayer, TotalGoalCount)
	allTargets := goflip.SwitchPressed(SwTargetG) && goflip.SwitchPressed((SwTargetO)) &&
		goflip.SwitchPressed(SwTargetA) && goflip.SwitchPressed(SwTargetL)
	//keep track of the G O A L targets, and reset the bank afterwards. Also light the 5k bonus

	log.Debugf("GoalObserver: AllTargets is %v", allTargets)

	//flash the goal light and reset target bank
	go func() {
		if allTargets {
			goflip.SolenoidFire(SolDropTargets)
		}
		goflip.AddScore(5000)
		goflip.LampFastBlink(LmpGoalLight)
		time.Sleep(3 * time.Second)
		goflip.LampOff(LmpGoalLight)
	}()

	if allTargets {

		cnt := utils.IncPlayerStat(game.CurrentPlayer, GoalTargetCount)
		log.Tracef("All targets down. goal count is now %d\n", cnt)
		if cnt >= 5 {
			goflip.LampOn(LmpRight25000Bonus)
			cnt -= 5
		}
	}

	//send back a command over the switch handler channel to call on choosepuck, hat trick determination, etc
	//game.BroadcastEvent(goflip.SwitchEvent{SwitchID: goalScored, Pressed: true})
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *GoalObserver) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *GoalObserver) PlayerUp(playerID int) {
	game := goflip.GetMachine()

	goflip.LampOff(p.goalBonusLights...)
	goflip.LampOff(LmpLeft25000Bonus, LmpRight25000Bonus)
	utils.SetPlayerStat(game.CurrentPlayer, GoalTargetCount, 0)
	utils.SetPlayerStat(game.CurrentPlayer, TotalGoalCount, 0)
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *GoalObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *GoalObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *GoalObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *GoalObserver) PlayerAdded(playerID int) {

}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *GoalObserver) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *GoalObserver) GameStart() {

}
