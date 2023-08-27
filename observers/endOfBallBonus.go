/*
Bonus - For each goal you get 5000 bonus light lit. At the end
 of the ball, the total number of letters * the 5000 will be awarded
  as a bonus. Example, if you spell MIKE BOSSY twice, and get
  MIK = that is 21 letters. You have 2 of the 5000 bonus lights lit,
  so you get 21 times 10,000 = 210,000 for a bonus.*/

package observer //this will probably be package main in your app

import (
	"sync"
	"time"

	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type EndOfBallBonus struct {
	//add your variables for the observer here
	letters          []int //all of the MikeBossy letters
	completedLetters []int
	bonus5000Points  []int
}

/*
the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*EndOfBallBonus)(nil)

/*
Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *EndOfBallBonus) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Debugln("EndOfBallBonus:Init called")

	p.letters = []int{LmpLetterM,
		LmpLetterI, LmpLetterK, LmpLetterE,
		LmpLetterB, LmpLetterO, LmpLetterS1, LmpLetterS2, LmpLetterY,
	}

	p.completedLetters = []int{LmpLeftCompleteLetters, LmpRightCompleteLetters}
	p.bonus5000Points = []int{Lmp5000G, Lmp5000O, Lmp5000A, Lmp5000L}

}

/*
SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *EndOfBallBonus) SwitchHandler(sw goflip.SwitchEvent) {
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *EndOfBallBonus) BallDrained() {

}

/*
PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up
*/
func (p *EndOfBallBonus) PlayerUp(playerID int) {
	goflip.LampOff(Lmp5000G, Lmp5000O, Lmp5000A, Lmp5000L)
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *EndOfBallBonus) PlayerStart(playerID int) {

}

/*PlayerEnd is called after every ball for the player is over*/
func (p *EndOfBallBonus) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	go func() {
		//target bonuses
		//number of goals is used as the shot multiplier
		game := goflip.GetMachine()
		goalCount := utils.GetPlayerStat(game.CurrentPlayer, BipGoalCount)

		//total number of pucks * the goal count is the bonus
		shotCount := utils.GetPlayerStat(game.CurrentPlayer, BipShotCount)

		for i := 0; i <= goalCount; i++ {
			for j := shotCount; j > 0; j-- {
				goflip.AddScore(1000)
				goflip.PlaySound(SndLetterBonus)

				time.Sleep(250 * time.Millisecond)
				//use the set bossyletters method..
				switch {
				case j == 18:
					goflip.LampOff(LmpRightCompleteLetters)
				case j == 9:
					goflip.LampOff(LmpLeftCompleteLetters)
				case j > 8:
					break
				default:
					goflip.LampOff(p.letters[j])
				}
			}
			time.Sleep(500 * time.Millisecond)
			//All lamps back on that were lit

			if i > goalCount-1 {
				break //exit for loop to not turn the lamps on, since already counted
			}

			//set the completed indicators based on number of times
			for i := 0; i < shotCount/9; i++ {
				goflip.LampOn(p.completedLetters[i])
			}

			for i := 0; i < shotCount%9; i++ {
				goflip.LampOn(p.letters[i])
			}
		}

		//goal bonuses

		//Turn on first all of the goal bonuses based on the score.
		targetCount := utils.GetPlayerStat(game.CurrentPlayer, GoalTargetCount)
		for i := 0; i < (targetCount % 5); i++ {
			goflip.LampOn(p.bonus5000Points[i])
		}
		//now turn them off as we count down

		for i := targetCount; i > 0; i-- {
			goflip.AddScore(5000)
			goflip.PlaySound(SndGoalBonus)
			if i == 4 {
				goflip.LampOff(LmpLeft25000Bonus, LmpRight25000Bonus)
			} else if i < 4 {
				goflip.LampOff(p.bonus5000Points[i-1])
			}
			time.Sleep(500 * time.Millisecond)
		}

		wait.Done()
	}()
}

/*PlayerFinish is called after the very last ball for the player is over (ball 3 for example)*/
func (p *EndOfBallBonus) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *EndOfBallBonus) PlayerAdded(playerID int) {

}

/*
GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode
*/
func (p *EndOfBallBonus) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *EndOfBallBonus) GameStart() {

}
