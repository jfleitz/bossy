/*
Bonus - For each goal you get 5000 bonus light lit. At the end
 of the ball, the total number of letters * the 5000 will be awarded
  as a bonus. Example, if you spell MIKE BOSSY twice, and get
  MIK = that is 21 letters. You have 2 of the 5000 bonus lights lit,
  so you get 21 times 10,000 = 210,000 for a bonus.*/

package main //this will probably be package main in your app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type endOfBallBonus struct {
	//add your variables for the observer here
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*endOfBallBonus)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *endOfBallBonus) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("endOfBallBonus:Init called")

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *endOfBallBonus) SwitchHandler(sw goflip.SwitchEvent) {
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *endOfBallBonus) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *endOfBallBonus) PlayerUp(playerID int) {

}

/*PlayerEnd is called after BallDrained. In a multiball game, this would be after the last
BallDrained method call*/
func (p *endOfBallBonus) PlayerEnd(playerID int) {
	//number of goals is the number of 5000 point values
	goalCount := getPlayerStat(game.CurrentPlayer, totalGoalCount)

	//total number of pucks * the goal count is the bonus
	puckCount := getPlayerStat(game.CurrentPlayer, totalPuckCount)

	//This is a compounded bonus right now (Hold Bonus). Thought we would try this out first before switching or adding an option
	game.AddScore(goalCount * puckCount * 5000)
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *endOfBallBonus) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *endOfBallBonus) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *endOfBallBonus) GameStart() {

}
