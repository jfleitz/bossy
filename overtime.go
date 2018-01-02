/*This is a continuation of the game, and based only on time. We should count down the time too.
Everything scores 1000 points, Goals score 22,000 points

Since this is a continuation, observer needs to be added last to the list
(so that no observer is skipped in giving points
*/

package main //this will probably be package main in your app

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type overTimeObserver struct {
	//add your variables for the observer here
	otWasPlayed bool
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*overTimeObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *overTimeObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("overTimeObserver:Init called")

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *overTimeObserver) SwitchHandler(sw goflip.SwitchEvent) {

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *overTimeObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *overTimeObserver) PlayerUp(playerID int) {

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *overTimeObserver) PlayerStart(playerID int) {

}

/*PlayerEnd is called after the ball for the player is over)*/
func (p *overTimeObserver) PlayerEnd(playerID int) {

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *overTimeObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *overTimeObserver) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *overTimeObserver) GameOver() {

	if p.otWasPlayed {
		return //already played overtime, so just exit
	}

	p.otWasPlayed = true
	game.BallInPlay = 4 //GoFlip would have set this to zero, which stops the game.. overriding to Ball 4

	//Determine which player is up in overtime.
	playerUp := 0
	otSecs := 0
	for i := 0; i < game.MaxPlayers; i++ {
		if getPlayerStat(i, otSeconds) > otSecs {
			playerUp = i
		}
	}

	game.CurrentPlayer = playerUp
	//let everyone know we are in overtime.
	game.BroadcastEvent(goflip.SwitchEvent{SwitchID: inOvertime, Pressed: true})

	//TODO: increment the time awarded onto the credit display here. For now, just simulate...

	go func() {
		for i := 0; i < otSecs; i++ {
			log.Infof("OT Seconds for player %d:%d", playerUp, i)
			time.After(300 * time.Millisecond)
		}
		//some sort of light sequence here?

		//launch the ball. at this point the other observers had plenty of time to get ready.
		game.SolenoidFire(solOuthole)
	}()

}

/*GameStart is called whenever a new game is started*/
func (p *overTimeObserver) GameStart() {
	p.otWasPlayed = false
}
