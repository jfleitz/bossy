/* Handles the Goalie and the 4 drop targets behind the goal
1. When a goal is scored, the Goal is worth 10,000 points
2. Plus 10,000 * each letter that the player has collected
3. Need to clear out the Letters for the Goal (but keep track of total letters
)

*/

package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type goalObserver struct {
	//add your variables for the observer here
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
	case swGoalie:
		game.AddScore(1000)
		return
	case swTargetG:
		break
	case swTargetO:
		break
	case swTargetA:
		break
	case swTargetL:
		break
	default:
		return
	}
	//For Hat Trick counting:
	goalfor := getPlayerStat(game.CurrentPlayer, hatTrickFor)
	switch goalfor {
	case passedToLeft:
		incPlayerStat(game.CurrentPlayer, leftGoalCount)
		break
	case passedToRight:
		incPlayerStat(game.CurrentPlayer, rightGoalCount)
		break
	case passedToSaucer:
		incPlayerStat(game.CurrentPlayer, saucerGoalCount)
		break
	}

	//For goal scoring:
	addScore := 10000
	addScore += 10000 * getPlayerStat(game.CurrentPlayer, bipPuckCount)
	game.AddScore(addScore)
	setPlayerStat(game.CurrentPlayer, bipPuckCount, 0)

	incPlayerStat(game.CurrentPlayer, bipGoalCount)
	incPlayerStat(game.CurrentPlayer, totalGoalCount)

	//play a sound

	//flash the goal light and reset target bank
	go func() {
		game.LampFlastBlink(lmpGoalLight)
		time.Sleep(3 * time.Second)
		game.LampOff(lmpGoalLight)
		game.SolenoidFire(solDropTargets)
	}()

	//send back a command over the switch handler channel to call on choosepuck, hat trick determination, etc
	game.BroadcastEvent(goflip.SwitchEvent{SwitchID: goalScored, Pressed: true})
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *goalObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *goalObserver) PlayerUp(playerID int) {
	game.SolenoidFire(solDropTargets)
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *goalObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *goalObserver) PlayerEnd(playerID int) {

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
