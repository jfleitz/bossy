package observer

/*DiagObserver is for handling mostly the test switch
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

type DiagObserver struct {
	//add your variables for the observer here
	testMode     int
	currentSound int8
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*DiagObserver)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *DiagObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("DiagObserver:Init called")
	p.testMode = notTesting
	game := goflip.GetMachine()
	game.TestMode = false
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/

func (p *DiagObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}
	game := goflip.GetMachine()

	switch sw.SwitchID {
	case SwTest:
		p.currentSound = 0 //always reset the sound when test is pressed
		p.testMode++
		if p.testMode > maxTests {
			p.testMode = notTesting
			game.TestMode = false
			//JAF TODO need to reset the displays etc when exiting
			goflip.ChangeGameState(goflip.GameEnded)
		} else {
			game.TestMode = true
			p.runTest()
		}
	case SwCredit:
		if game.TestMode && (p.testMode == testSounds) {
			//increment which sound is being tested
			p.currentSound++
			if p.currentSound > 15 {
				p.currentSound = 0
			}
		}
	}

}
func (p *DiagObserver) runTest() {
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

func (p *DiagObserver) testDisplays() {
	goflip.ClearScores()
	var score int32 = 1111111
	var ball int8 = 11

	for {
		for disp := 1; disp < 5; disp++ {
			goflip.SetDisplay(disp, score)
		}

		goflip.SetCreditDisp(ball)
		goflip.SetBallInPlayDisp(ball)
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

func (p *DiagObserver) testAllLamps() {
	for {
		goflip.LampOn(1)
		goflip.SetBallInPlayDisp(0)
		time.Sleep(time.Second * 1)
		goflip.LampOff(1)
		time.Sleep(time.Second * 1)
		if p.testMode != testAllLamps {
			return
		}
	}
}

func (p *DiagObserver) testSounds() {
	var playedSound int8
	playedSound = -1
	for {
		if playedSound != p.currentSound {
			playedSound = p.currentSound
			goflip.SetBallInPlayDisp(playedSound)
			goflip.PlaySound(byte(playedSound))
		}

		time.Sleep(500)

		if p.testMode != testSounds {
			return
		}
	}
}

func (p *DiagObserver) testLamps() {

	for {
		for id := 0; id < 65; id++ {
			goflip.LampOn(id)
			goflip.SetBallInPlayDisp(int8(id))
			time.Sleep(time.Second * 1)
			goflip.LampOff(id)

			if p.testMode != testLamps {
				return
			}
		}
	}
}

func (p *DiagObserver) testSolenoids() {
	for {
		for id := 0; id < 16; id++ {
			goflip.SolenoidFire(id)
			goflip.SetBallInPlayDisp(int8(id))
			time.Sleep(time.Second * 1)

			if p.testMode != testSolenoids {
				return
			}
		}
	}
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *DiagObserver) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *DiagObserver) PlayerUp(playerID int) {
	log.Traceln("DiagObserver:PlayerUp()")
}

/*PlayerEnd is called after every ball for the player is over*/
func (p *DiagObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	log.Traceln("DiagObserver:PlayerEnd()")
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *DiagObserver) PlayerStart(playerID int) {
	log.Traceln("DiagObserver:PlayerStart()")
}

/*
PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *DiagObserver) PlayerFinish(playerID int) {
	log.Traceln("DiagObserver:PlayerFinish()")
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *DiagObserver) PlayerAdded(playerID int) {
	log.Traceln("DiagObserver:PlayerAdded()")
}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *DiagObserver) GameOver() {
	log.Traceln("DiagObserver:GameOver()")
}

/*GameStart is called whenever a new game is started*/
func (p *DiagObserver) GameStart() {
	log.Traceln("DiagObserver:GameStart()")
}
