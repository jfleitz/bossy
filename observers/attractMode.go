package observer

/*AttractMode is for handling mostly the test switch
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

// in order from the bottom of the playfield to the top
type AttractMode struct {
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

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*AttractMode)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (a *AttractMode) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("AttractMode:Init called")

	a.bottom = []int{LmpBottomLeftSpecial, LmpBottomRightGreenSpot, LmpBottomRightSpecial, LmpSamePlayerShootAgain} //jaf todo need bottom left green spot
	a.bonus5000Points = []int{Lmp5000Bonus1, Lmp5000Bonus2, Lmp5000Bonus3, Lmp5000Bonus4}
	a.bottom25kOvertime = []int{Lmp25000Bonus, LmpLeftReturnLaneOrangeSpot, LmpRightReturnLaneOrangeSpot}
	a.completedLetters = []int{LmpLeftCompleteLetters, LmpRightCompleteLetters}
	a.bossyLetters = []int{LmpLetterB, LmpLetterO, LmpLetterS1, LmpLetterS2, LmpLetterY}
	a.mikeLetters = []int{LmpLetterM, LmpLetterI, LmpLetterK, LmpLetterE}
	a.next1 = []int{LmpPointLaneWhiteSpot, LmpBottomTargetWhiteSpot}
	a.next2 = []int{LmpPointLaneSpecial, LmpMiddleRightTarget, LmpTargetsOrangeSpot}
	a.next3 = []int{LmpLeftTarget, LmpTopTargetWhiteSpot}
	a.next4 = []int{LmpOvertimeLeftOfGoal} //need saucer
	a.next5 = []int{LmpGoalOnLeftOrangeSpot, LmpGoalieWhiteSpot}
	a.topLanes = []int{LmpTopLeftLane, LmpTopMiddleLane, LmpTopRightLane}
	a.goalAndOvetime = []int{LmpGoalLight, LmpTopOrangeSpot}

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/

func (a *AttractMode) SwitchHandler(sw goflip.SwitchEvent) {

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (a *AttractMode) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (a *AttractMode) PlayerUp(playerID int) {
}

/*PlayerEnd is called after every ball for the player is over*/
func (a *AttractMode) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (a *AttractMode) PlayerStart(playerID int) {
}

/*
PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (a *AttractMode) PlayerFinish(playerID int) {
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (a *AttractMode) PlayerAdded(playerID int) {
}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (a *AttractMode) GameOver() {
	log.Debugln("AttractMode:GameOver()")

	go func() {
		for {
			if goflip.GetGameState() != goflip.GameEnded {
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
func (a *AttractMode) GameStart() {
	log.Debugln("AttractMode:GameStart()")
}

func (a *AttractMode) controlLights(lights []int, blink bool, duration int, offAfter bool) {
	if blink {
		goflip.LampFastBlink(lights...)
	} else {
		goflip.LampOn(lights...)
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)
	if offAfter {
		goflip.LampOff(lights...)
	}
}

func (a *AttractMode) offLights(lights []int, wait int) {
	goflip.LampOff(lights...)
	time.Sleep(time.Duration(wait) * time.Millisecond)
}
