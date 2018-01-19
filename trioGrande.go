/*
Trio Grande Line - This is to pay homage to the best hokcey line in NHL History; Bryan Trottier, Clark Gillies, and Mike Bossy forming the Long Island Lighting Company (4 consecutive Stanley Cup wins). Similar to Hat Trick, this is based on timing from the flipper inlane switches and the saucer switch. You must complete each within 5 seconds of each other. Once you do so, the Goal light will blink to signify you have LILCO Line ready. Hitting the goal will award you with 220,000 points.

In any order, you must hit within 5 seconds of each other:
Either Left Flipper Inlane switch
Either Right Flipper Inlane switch
The Upper Right saucer
This will flash the Goal Light and make the next goal 220,000 points*/

package main //this will probably be package main in your app

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type trioGrande struct {
	//add your variables for the observer here
	trioLeft   bool
	trioRight  bool
	trioSaucer bool
	passedTo   chan int
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*trioGrande)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *trioGrande) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("trioGrande:Init called")
	p.trioLeft = false
	p.trioRight = false
	p.trioSaucer = false
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *trioGrande) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	//start the appropriate timer and cancel any other timer going on (only one person can score a goal)
	switch sw.SwitchID {
	case swInnerLeftLane:
		if !p.trioLeft {
			p.passedTo <- passedToLeft
		}
		break
	case swMiddleLeftLane:
		if !p.trioLeft {
			p.passedTo <- passedToLeft
		}
		break
	case swInnerRightLane:
		if !p.trioRight {
			p.passedTo <- passedToRight
		}
		break
	case swMiddleRightLane:
		if !p.trioRight {
			p.passedTo <- passedToRight
		}
		break
	case swSaucer:
		if !p.trioSaucer {
			p.passedTo <- passedToSaucer
		}
		return
	case goalScored:
		p.checkCount()
		return
	default:
		return
	}
}

func (p *trioGrande) checkCount() {
	if p.trioLeft && p.trioRight && p.trioSaucer {
		//award
		game.AddScore(220000)

		go func() {
			//TODO flash all of the lights for the trio grande lanes/saucer with the goal light
			//game.LampFlastBlink(lmp)
			time.After(3 * time.Second)
			game.LampOff(lmpGoalLight)
		}()
	}

	p.trioLeft = false
	p.trioRight = false
	p.trioSaucer = false
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *trioGrande) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *trioGrande) PlayerUp(playerID int) {
	p.trioLeft = false
	p.trioRight = false
	p.trioSaucer = false

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *trioGrande) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *trioGrande) PlayerEnd(playerID int) {

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *trioGrande) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *trioGrande) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *trioGrande) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *trioGrande) GameStart() {

}

func (p *trioGrande) lineTimeout() {
	//every 3 seconds reset the hattrick variable.
	go func() {
		var timeToReset time.Duration = 5

		for {
			select {
			case newShooter := <-p.passedTo:
				switch newShooter {
				case passedToLeft:
					p.trioLeft = true
					break
				case passedToRight:
					p.trioRight = true
					break
				case passedToSaucer:
					p.trioSaucer = true
					break
				}
				break
			case <-time.After(timeToReset * time.Second):
				p.trioLeft = false
				p.trioRight = false
				p.trioSaucer = false
				break
			}
		}
	}()

}
