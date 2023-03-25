/*This is a continuation of the game, and based only on time. We should count down the time too.
Everything scores 1000 points, Goals score 22,000 points

Since this is a continuation, observer needs to be added last to the list
(so that no observer is skipped in giving points
*/

package observer //this will probably be package main in your app

import (
	"sync"
	"time"

	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type OverTimeObserver struct {
	//add your variables for the observer here
	otWasPlayed bool
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*OverTimeObserver)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *OverTimeObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("OverTimeObserver:Init called")

}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *OverTimeObserver) SwitchHandler(sw goflip.SwitchEvent) {

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *OverTimeObserver) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *OverTimeObserver) PlayerUp(playerID int) {

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *OverTimeObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *OverTimeObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (p *OverTimeObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *OverTimeObserver) PlayerAdded(playerID int) {

}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *OverTimeObserver) GameOver() {

	if p.otWasPlayed {
		return //already played overtime, so just exit
	}

	p.otWasPlayed = true
	game := goflip.GetMachine()

	game.BallInPlay = 4 //GoFlip would have set this to zero, which stops the game.. overriding to Ball 4

	//Determine which player is up in overtime.
	playerUp := 0
	otSecs := 0
	for i := 0; i < game.MaxPlayers; i++ {
		if utils.GetPlayerStat(i, OTSeconds) > otSecs {
			playerUp = i
		}
	}

	game.CurrentPlayer = playerUp
	//let everyone know we are in overtime.
	goflip.BroadcastEvent(goflip.SwitchEvent{SwitchID: InOvertime, Pressed: true})

	//TODO: increment the time awarded onto the credit display here. For now, just simulate...

	go func() {
		for i := 0; i < otSecs; i++ {
			log.Debugf("OT Seconds for player %d:%d", playerUp, i)
			time.After(300 * time.Millisecond)
		}
		//some sort of light sequence here?

		//launch the ball. at this point the other observers had plenty of time to get ready.
		goflip.SolenoidFire(SolOuthole)
	}()

}

/*GameStart is called whenever a new game is started*/
func (p *OverTimeObserver) GameStart() {
	p.otWasPlayed = false
}
