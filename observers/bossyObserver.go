package observer

/*bossyObserver is for game specific call backs that are not tied to a specific feature (like handling ball save)
 */

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type BossyObserver struct {
	//add your variables for the observer here
	firstShot bool //used to see if we need to play the shooter lane launch sound or not
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*BossyObserver)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *BossyObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("bossyObserver:Init called")

}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *BossyObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	switch sw.SwitchID {
	case SwShooterLane:
		if p.firstShot {
			goflip.PlaySound(SndShooter)
			p.firstShot = false
		}
	}

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *BossyObserver) BallDrained() {
	if goflip.GetGameState() != goflip.InProgress {
		return
	}

	log.Debugln("bossyObsv:BallDrained()")
	game := goflip.GetMachine()

	if !inWarmUpPeriod {
		log.Debugln("outhole: not in warm up period")
		if game.BallScore == 0 {
			log.Debugln("0 points by ball. ejecting ball")
			go BallLaunch()
		} else {
			goflip.ChangePlayerState(goflip.EndPlayer)
		}
	} else {
		//go ahead and eject it again
		log.Debugln("warmup period-firing solenoid")
		go BallLaunch()
	}
}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *BossyObserver) PlayerUp(playerID int) {
	log.Debugf("bossyObsv:PlayerUp() for player %d", playerID)
	p.firstShot = true
	goflip.SolenoidFire(SolDropTargets)

	BallLaunch()

	//turn off the other player up lights
	goflip.LampOff(LmpPlayer1, LmpPlayer2, LmpPlayer3, LmpPlayer4)
	//turn on appropriate Player Up Light
	goflip.LampSlowBlink(LmpPlayer1 + playerID - 1)
}

/*PlayerEnd is called after every ball for the player is over*/
func (p *BossyObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
	//turn off the player up light
	log.Debugln("bossyObsv:PlayerEnd()")
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *BossyObserver) PlayerStart(playerID int) {
	log.Debugln("bossyObsv:PlayerStart()")

}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *BossyObserver) PlayerFinish(playerID int) {
	log.Debugf("bossyObsv:PlayerFinish: %d\n", playerID)
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *BossyObserver) PlayerAdded(playerID int) {
	//turn on the additional player light
	log.Debugf("bossyObsv:PlayerAdded: %d\n", playerID)
	if playerID == 1 {
		goflip.PlaySound(SndStartGame)
	} else {
		goflip.PlaySound(SndAddPlayer)
	}
}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *BossyObserver) GameOver() {

	goflip.PlaySound(SndGameOver)
	log.Debugln("bossyObsv:GameOver()")

	game := goflip.GetMachine()

	goflip.SetCreditDisp(int8(game.Credits))
	//turn off all player up lights, and number of players
	goflip.LampOff(LmpPlayer1, LmpPlayer2, LmpPlayer3, LmpPlayer4)
	goflip.LampSlowBlink(LmpGameOver)
	goflip.LampOff(LmpPeriod)
	goflip.FlipperControl(false)
}

/*GameStart is called whenever a new game is started*/
func (p *BossyObserver) GameStart() {
	log.Debugln("bossyObserver:GameStart()")
	goflip.LampOff(LmpGameOver)
	goflip.LampOn(LmpPeriod) //Ball in play light behind the backglass
	goflip.FlipperControl(true)
}

func BallLaunch() {
	//	game.NextUp()
	time.Sleep(1 * time.Second)
	goflip.SolenoidFire(SolOuthole)
}
