/*Overtime awards. On the last ball, the overtime lights are lit.
For each time you hit one, 1 second is added to the overtime.
If there are more than one player playing, whomever has the
higher overtime value gets the timed extra ball.*/

package main //this will probably be package main in your app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type overTimeObserver struct {
	//add your variables for the observer here
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

/*PlayerEnd is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *overTimeObserver) PlayerEnd(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *overTimeObserver) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *overTimeObserver) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *overTimeObserver) GameStart() {

}
