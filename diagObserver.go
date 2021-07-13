package main

/*diagObserver is for handling mostly the test switch
 */

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

/*
 */

const (
	notTesting    = 0
	testDisplays  = 1
	testAllLamps  = 2
	testLamps     = 3
	testSolenoids = 4
	testSounds    = 5

	maxTests = 5
)

type diagObserver struct {
	//add your variables for the observer here
	testMode     int
	currentSound int8
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*diagObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *diagObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("diagObserver:Init called")
	p.testMode = notTesting
	game.TestMode = false
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/

func (p *diagObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}
	switch sw.SwitchID {
	case swTest:
		p.currentSound = 0 //always reset the sound when test is pressed
		p.testMode++
		if p.testMode > maxTests {
			p.testMode = notTesting
			game.TestMode = false
			//JAF TODO need to reset the displays etc when exiting
			game.ChangeGameState(goflip.GameOver)
		} else {
			game.TestMode = true
			p.runTest()
		}
	case swCredit:
		if game.TestMode && (p.testMode == testSounds) {
			//increment which sound is being tested
			p.currentSound++
			if p.currentSound > 15 {
				p.currentSound = 0
			}
		}
	}

}
func (p *diagObserver) runTest() {
	switch p.testMode {
	case testDisplays:
		go p.testDisplays()
	case testAllLamps:
		go p.testAllLamps()
	case testLamps:
		go p.testLamps()
	case testSolenoids:
		go p.testSolenoids()
	case testSounds:
		go p.testSounds()
	}
}

func (p *diagObserver) testDisplays() {
	game.ClearScores()
	var score int32 = 1111111
	var ball int8 = 11

	for {
		for disp := 1; disp < 5; disp++ {
			game.SetDisplay(disp, score)
		}

		game.SetCreditDisp(ball)
		game.SetBallInPlayDisp(ball)
		score += 1111111
		ball += 11

		if score > 9999999 {
			score = 1111111
			ball = 11
		}

		time.Sleep(time.Millisecond * 500)
		if p.testMode != testDisplays {
			return
		}
	}
}

func (p *diagObserver) testAllLamps() {
	for {
		game.LampOn(1)
		game.SetBallInPlayDisp(0)
		time.Sleep(time.Second * 1)
		game.LampOff(1)
		time.Sleep(time.Second * 1)
		if p.testMode != testAllLamps {
			return
		}
	}
}

func (p *diagObserver) testSounds() {
	var playedSound int8
	playedSound = -1
	for {
		if playedSound != p.currentSound {
			playedSound = p.currentSound
			game.SetBallInPlayDisp(playedSound)
			game.PlaySound(byte(playedSound))
		}

		time.Sleep(500)

		if p.testMode != testSounds {
			return
		}
	}
}

func (p *diagObserver) testLamps() {

	for {
		for id := 0; id < 65; id++ {
			game.LampOn(id)
			game.SetBallInPlayDisp(int8(id))
			time.Sleep(time.Second * 1)
			game.LampOff(id)

			if p.testMode != testLamps {
				return
			}
		}
	}
}

func (p *diagObserver) testSolenoids() {
	for {
		for id := 0; id < 16; id++ {
			game.SolenoidFire(id)
			game.SetBallInPlayDisp(int8(id))
			time.Sleep(time.Second * 1)

			if p.testMode != testSolenoids {
				return
			}
		}
	}
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *diagObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *diagObserver) PlayerUp(playerID int) {
	log.Infoln("diagObserver:PlayerUp()")
}

/*PlayerEnd is called after every ball for the player is over*/
func (p *diagObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	log.Infoln("diagObserver:PlayerEnd()")
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *diagObserver) PlayerStart(playerID int) {
	log.Infoln("diagObserver:PlayerStart()")
}

/*PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *diagObserver) PlayerFinish(playerID int) {
	log.Infoln("diagObserver:PlayerFinish()")
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *diagObserver) PlayerAdded(playerID int) {
	log.Infoln("diagObserver:PlayerAdded()")
}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *diagObserver) GameOver() {
	log.Infoln("diagObserver:GameOver()")
}

/*GameStart is called whenever a new game is started*/
func (p *diagObserver) GameStart() {
	log.Infoln("diagObserver:GameStart()")
}
