package main

/*bossyObserver is for game specific call backs that are not tied to a specific feature (like handling ball save)
 */

import (
	"sync"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type bossyObserver struct {
	//add your variables for the observer here

	firstShot bool //used to see if we need to play the shooter lane launch sound or not
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*bossyObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *bossyObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("bossyObserver:Init called")

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *bossyObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	switch sw.SwitchID {
	case swShooterLane:
		if p.firstShot {
			game.PlaySound(sndShooter)
			p.firstShot = false
		}
	}

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *bossyObserver) BallDrained() {
	if game.GetGameState() != goflip.GameStart {
		return
	}

	log.Debugln("bossyObsv:BallDrained()")

	if !inWarmUpPeriod {
		log.Debugln("outhole: not in warm up period")
		if game.BallScore == 0 {
			log.Debugln("0 points by ball. ejecting ball")
			go ballLaunch()
		} else {
			game.ChangePlayerState(goflip.PlayerEnd)
		}
	} else {
		//go ahead and eject it again
		log.Debugln("warmup period-firing solenoid")
		go ballLaunch()
	}
}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *bossyObserver) PlayerUp(playerID int) {
	log.Debugf("bossyObsv:PlayerUp() for player %d", playerID)
	p.firstShot = true
	game.SolenoidFire(solOuthole)

	//turn off the other player up lights
	game.LampOff(lmpPlayer1, lmpPlayer2, lmpPlayer3, lmpPlayer4)
	//turn on appropriate Player Up Light
	game.LampSlowBlink(lmpPlayer1 + playerID - 1)
}

/*PlayerEnd is called after every ball for the player is over*/
func (p *bossyObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	//turn off the player up light
	log.Debugln("bossyObsv:PlayerEnd()")
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *bossyObserver) PlayerStart(playerID int) {
	log.Debugln("bossyObsv:PlayerStart()")

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *bossyObserver) PlayerFinish(playerID int) {
	log.Debugf("bossyObsv:PlayerFinish: %d\n", playerID)
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *bossyObserver) PlayerAdded(playerID int) {
	//turn on the additional player light
	log.Debugf("bossyObsv:PlayerAdded: %d\n", playerID)
	if playerID == 1 {
		game.PlaySound(sndStartGame)
	} else {
		game.PlaySound(sndAddPlayer)
	}
}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *bossyObserver) GameOver() {

	game.PlaySound(sndGameOver)
	log.Debugln("bossyObsv:GameOver()")

	game.SetCreditDisp(int8(game.Credits))
	//turn off all player up lights, and number of players
	game.LampOff(lmpPlayer1, lmpPlayer2, lmpPlayer3, lmpPlayer4)
	game.LampSlowBlink(lmpGameOver)
	game.LampOff(lmpPeriod)
	game.FlipperControl(false)
}

/*GameStart is called whenever a new game is started*/
func (p *bossyObserver) GameStart() {
	log.Debugln("bossyObserver:GameStart()")
	game.LampOff(lmpGameOver)
	game.LampOn(lmpPeriod) //Ball in play light behind the backglass
	game.FlipperControl(true)
}
