/*
Bonus - For each goal you get 5000 bonus light lit. At the end
 of the ball, the total number of letters * the 5000 will be awarded
  as a bonus. Example, if you spell MIKE BOSSY twice, and get
  MIK = that is 21 letters. You have 2 of the 5000 bonus lights lit,
  so you get 21 times 10,000 = 210,000 for a bonus.*/

package main //this will probably be package main in your app

import (
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type endOfBallBonus struct {
	//add your variables for the observer here
	letters          []int //all of the MikeBossy letters
	completedLetters []int
	bonus5000Points  []int
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

	p.letters = []int{lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterB, lmpLetterO, lmpLetterS1, lmpLetterS2, lmpLetterY,
	}

	p.completedLetters = []int{lmpLeftCompleteLetters, lmpRightCompleteLetters}
	p.bonus5000Points = []int{lmp5000Bonus1, lmp5000Bonus2, lmp5000Bonus3, lmp5000Bonus4}

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

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *endOfBallBonus) PlayerStart(playerID int) {

}

/*PlayerEnd is called after every ball for the player is over*/
func (p *endOfBallBonus) PlayerEnd(playerID int) {
	//target bonuses
	//number of goals is used as the shot multiplier
	goalCount := getPlayerStat(game.CurrentPlayer, totalGoalCount)

	//total number of pucks * the goal count is the bonus
	shotCount := getPlayerStat(game.CurrentPlayer, totalShotCount)

	for i := 0; i < goalCount; i++ {
		for j := shotCount; j > 0; j-- {
			game.AddScore(1000)
			game.PlaySound(sndPuckBounce)
			time.Sleep(250 * time.Millisecond)
			switch {
			case j == 18:
				game.LampOff(lmpRightCompleteLetters)
			case j == 9:
				game.LampOff(lmpLeftCompleteLetters)
			case j > 18:
				break
			default:
				game.LampOff(p.letters[j])
			}
		}
		time.Sleep(500 * time.Millisecond)
		//All lamps back on that were lit

		if i == goalCount-1 {
			continue //exit for loop to not turn the lamps on, since already counted
		}

		//set the completed indicators based on number of times
		for i := 0; i < shotCount/9; i++ {
			game.LampOn(p.completedLetters[i])
		}

		for i := 0; i < shotCount%9; i++ {
			game.LampOn(p.letters[i])
		}
	}

	//goal bonuses
	targetCount := getPlayerStat(game.CurrentPlayer, goalTargetCount)

	for i := targetCount; i > 0; i-- {
		game.AddScore(5000)
		game.PlaySound(sndWhistle)

		if i == 4 {
			game.LampOff(lmp25000Bonus)
		} else if i < 4 {
			game.LampOff(p.bonus5000Points[i])
		}
	}

	game.AddScore(targetCount * 5000) //add number of goal targets completed X 5000
}

/*PlayerFinish is called after the very last ball for the player is over (ball 3 for example)*/
func (p *endOfBallBonus) PlayerFinish(playerID int) {

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
