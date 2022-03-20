package main

/*attractMode is for handling mostly the test switch
 */

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

const eachSetTimeOn = 100

/*
 */

//in order from the bottom of the playfield to the top
type attractMode struct {
	bottom            []int
	bonus5000Points   []int
	bottom25kOvertime []int
	completedLetters  []int
	bossyLetters      []int //all of the MikeBossy letters
	mikeLetters       []int //all of the MikeBossy letters
	next1             []int //2 lights, lower right target and target lane on left
	next2             []int //3 lights, special on left lane, ovetime on middle right target and middle right target
	next3             []int //2 lights, target on left, top right target
	next4             []int //2 lights, left spot around goal, and saucer overtime
	next5             []int //2 lights, overtime for around goal, and goalie
	topLanes          []int //top lanes spot, 3 lights
	goalAndOvetime    []int //2, top overtime, and goal light

}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*attractMode)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (a *attractMode) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("attractMode:Init called")

	a.bottom = []int{lmpBottomLeftSpecial, lmpBottomRightGreenSpot, lmpBottomRightSpecial, lmpSamePlayerShootAgain} //jaf todo need bottom left green spot
	a.bonus5000Points = []int{lmp5000Bonus1, lmp5000Bonus2, lmp5000Bonus3, lmp5000Bonus4}
	a.bottom25kOvertime = []int{lmp25000Bonus, lmpLeftReturnLaneOrangeSpot, lmpRightReturnLaneOrangeSpot}
	a.completedLetters = []int{lmpLeftCompleteLetters, lmpRightCompleteLetters}
	a.bossyLetters = []int{lmpLetterB, lmpLetterO, lmpLetterS1, lmpLetterS2, lmpLetterY}
	a.mikeLetters = []int{lmpLetterM, lmpLetterI, lmpLetterK, lmpLetterE}
	a.next1 = []int{lmpPointLaneWhiteSpot, lmpBottomTargetWhiteSpot}
	a.next2 = []int{lmpPointLaneSpecial, lmpMiddleRightTarget, lmpTargetsOrangeSpot}
	a.next3 = []int{lmpLeftTarget, lmpTopTargetWhiteSpot}
	a.next4 = []int{lmpOvertimeLeftOfGoal} //need saucer
	a.next5 = []int{lmpGoalOnLeftOrangeSpot, lmpGoalieWhiteSpot}
	a.topLanes = []int{lmpTopLeftLane, lmpTopMiddleLane, lmpTopRightLane}
	a.goalAndOvetime = []int{lmpGoalLight, lmpTopOrangeSpot}

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/

func (a *attractMode) SwitchHandler(sw goflip.SwitchEvent) {

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (a *attractMode) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (a *attractMode) PlayerUp(playerID int) {
}

/*PlayerEnd is called after every ball for the player is over*/
func (a *attractMode) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (a *attractMode) PlayerStart(playerID int) {
}

/*PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)*/
func (a *attractMode) PlayerFinish(playerID int) {
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (a *attractMode) PlayerAdded(playerID int) {
}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (a *attractMode) GameOver() {
	log.Debugln("attractMode:GameOver()")

	go func() {
		for {
			if game.GetGameState() != goflip.GameOver {
				return
			}

			a.controlLights(a.bottom, false, eachSetTimeOn, false)
			a.controlLights(a.bonus5000Points, false, eachSetTimeOn, false)
			a.controlLights(a.bottom25kOvertime, false, eachSetTimeOn, false)
			a.controlLights(a.completedLetters, false, eachSetTimeOn, false)
			a.controlLights(a.bossyLetters, false, eachSetTimeOn, false)
			a.controlLights(a.mikeLetters, false, eachSetTimeOn, false)
			a.controlLights(a.next1, false, eachSetTimeOn, false)
			a.controlLights(a.next2, false, eachSetTimeOn, false)
			a.controlLights(a.next3, false, eachSetTimeOn, false)
			a.controlLights(a.next4, false, eachSetTimeOn, false)
			a.controlLights(a.next5, false, eachSetTimeOn, false)
			a.controlLights(a.topLanes, false, eachSetTimeOn, false)
			a.controlLights(a.goalAndOvetime, false, eachSetTimeOn, false)

			time.Sleep(250 * time.Millisecond)
			a.offLights(a.goalAndOvetime, eachSetTimeOn)
			a.offLights(a.topLanes, eachSetTimeOn)
			a.offLights(a.next5, eachSetTimeOn)
			a.offLights(a.next4, eachSetTimeOn)
			a.offLights(a.next3, eachSetTimeOn)
			a.offLights(a.next2, eachSetTimeOn)
			a.offLights(a.next1, eachSetTimeOn)
			a.offLights(a.mikeLetters, eachSetTimeOn)
			a.offLights(a.bossyLetters, eachSetTimeOn)
			a.offLights(a.completedLetters, eachSetTimeOn)
			a.offLights(a.bottom25kOvertime, eachSetTimeOn)
			a.offLights(a.bonus5000Points, eachSetTimeOn)
			a.offLights(a.bottom, eachSetTimeOn)
			time.Sleep(500 * time.Millisecond)

		}
	}()

}

/*GameStart is called whenever a new game is started*/
func (a *attractMode) GameStart() {
	log.Debugln("attractMode:GameStart()")
}

func (a *attractMode) controlLights(lights []int, blink bool, duration int, offAfter bool) {
	if blink {
		game.LampFastBlink(lights...)
	} else {
		game.LampOn(lights...)
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)
	if offAfter {
		game.LampOff(lights...)
	}
}

func (a *attractMode) offLights(lights []int, wait int) {
	game.LampOff(lights...)
	time.Sleep(time.Duration(wait) * time.Millisecond)
}
