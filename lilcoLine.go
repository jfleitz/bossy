/*
LILCO Line - This is to pay homage to the best hokcey line in NHL History; Bryan Trottier, Clark Gillies, and Mike Bossy forming the Long Island Lighting Company (4 consecutive Stanley Cup wins). Similar to Hat Trick, this is based on timing from the flipper inlane switches and the saucer switch. You must complete each within 5 seconds of each other. Once you do so, the Goal light will blink to signify you have LILCO Line ready. Hitting the goal will award you with 220,000 points.

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

type lilcoLine struct {
	//add your variables for the observer here
	lilcoLeft   bool
	lilcoRight  bool
	lilcoSaucer bool
	passedTo    chan int
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*lilcoLine)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *lilcoLine) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("lilcoLine:Init called")
	p.lilcoLeft = false
	p.lilcoRight = false
	p.lilcoSaucer = false
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *lilcoLine) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	//start the appropriate timer and cancel any other timer going on (only one person can score a goal)
	switch sw.SwitchID {
	case swInnerLeftLane:
		if !p.lilcoLeft {
			p.passedTo <- passedToLeft
		}
		break
	case swMiddleLeftLane:
		if !p.lilcoLeft {
			p.passedTo <- passedToLeft
		}
		break
	case swInnerRightLane:
		if !p.lilcoRight {
			p.passedTo <- passedToRight
		}
		break
	case swMiddleRightLane:
		if !p.lilcoRight {
			p.passedTo <- passedToRight
		}
		break
	case swSaucer:
		if !p.lilcoSaucer {
			p.passedTo <- passedToSaucer
		}
		return //so that we don't get the 1000 points added below
	case goalScored:
		//p.checkShotCount()
		return
	default:
		return
	}
}

func (p *lilcoLine) checkLilcoCount() {
	if p.lilcoLeft && p.lilcoRight && p.lilcoSaucer {
		//award
		game.AddScore(220000)

		go func() {
			game.LampFlastBlink(lmpGoalLight)
			time.After(3 * time.Second)
			game.LampOff(lmpGoalLight)
		}()
	}

	p.lilcoLeft = false
	p.lilcoRight = false
	p.lilcoSaucer = false
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *lilcoLine) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *lilcoLine) PlayerUp(playerID int) {
	p.lilcoLeft = false
	p.lilcoRight = false
	p.lilcoSaucer = false

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *lilcoLine) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *lilcoLine) PlayerEnd(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *lilcoLine) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *lilcoLine) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *lilcoLine) GameStart() {

}

func (p *lilcoLine) lineTimeout() {
	//every 3 seconds reset the hattrick variable.
	go func() {
		var timeToReset time.Duration = 5

		for {
			select {
			case newShooter := <-p.passedTo:
				switch newShooter {
				case passedToLeft:
					p.lilcoLeft = true
					break
				case passedToRight:
					p.lilcoRight = true
					break
				case passedToSaucer:
					p.lilcoSaucer = true
					break
				}
				break
			case <-time.After(timeToReset * time.Second):
				p.lilcoLeft = false
				p.lilcoRight = false
				p.lilcoSaucer = false
				break
			}
		}
	}()

}
