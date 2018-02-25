/*
Hat Trick - If you hit 3 goals in the same game by the same flipper.

For the lower flippers, a goal must be scored within 3 seconds after the ball goes from one of the flipper inlanes to be counted toward the "Hat Trick count" for that flipper.
For the upper flipper, a goal must be scored within 2 seconds after the ball is ejected from the saucer to be counted toward the "Hat Trick count" of the upper flipper.
For each Hat Trick (3 shots by the same flipper), 22,000 points is awarded. (Award at end of game)*/

package main //this will probably be package main in your app

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type hatTrick struct {
	passedTo chan int
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*hatTrick)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *hatTrick) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("hatTrick:Init called")
	p.passedTo = make(chan int, 1)
	p.hatTrickMonitor()
}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *hatTrick) SwitchHandler(sw goflip.SwitchEvent) {
	if !sw.Pressed {
		return
	}

	//start the appropriate timer and cancel any other timer going on (only one person can score a goal)
	switch sw.SwitchID {
	case swInnerLeftLane:
		p.passedTo <- passedToLeft
		break
	case swMiddleLeftLane:
		p.passedTo <- passedToLeft
		break
	case swInnerRightLane:
		p.passedTo <- passedToRight
		break
	case swMiddleRightLane:
		p.passedTo <- passedToRight
		break
	case swSaucer:
		p.passedTo <- passedToSaucer
		return //so that we don't get the 1000 points added below
	case goalScored:
		p.checkShotCount()
		return
	default:
		return
	}

	game.AddScore(1000)
}

func (p *hatTrick) checkShotCount() {
	if getPlayerStat(game.CurrentPlayer, leftGoalCount) >= 3 {
		setPlayerStat(game.CurrentPlayer, leftGoalCount, 0)
		incPlayerStat(game.CurrentPlayer, hatTrickCount)
	}

	if getPlayerStat(game.CurrentPlayer, rightGoalCount) >= 3 {
		setPlayerStat(game.CurrentPlayer, rightGoalCount, 0)
		incPlayerStat(game.CurrentPlayer, hatTrickCount)
	}

	if getPlayerStat(game.CurrentPlayer, saucerGoalCount) >= 3 {
		setPlayerStat(game.CurrentPlayer, saucerGoalCount, 0)
		incPlayerStat(game.CurrentPlayer, hatTrickCount)
	}
}

func (p *hatTrick) hatTrickMonitor() {
	//every 3 seconds reset the hattrick variable.
	go func() {
		var timeToReset time.Duration = 3

		for {
			select {
			case newShooter := <-p.passedTo:
				if newShooter == passedToSaucer {
					timeToReset = 2
				} else {
					timeToReset = 3
				}
				setPlayerStat(game.CurrentPlayer, hatTrickFor, newShooter) //this will set the timer back to 3 seconds as well

			case <-time.After(timeToReset * time.Second):
				setPlayerStat(game.CurrentPlayer, hatTrickFor, 0) //reset back to noone
			}
		}
	}()

}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *hatTrick) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *hatTrick) PlayerUp(playerID int) {

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *hatTrick) PlayerStart(playerID int) {
	setPlayerStat(game.CurrentPlayer, hatTrickCount, 0)
	setPlayerStat(game.CurrentPlayer, leftGoalCount, 0)
	setPlayerStat(game.CurrentPlayer, rightGoalCount, 0)
	setPlayerStat(game.CurrentPlayer, saucerGoalCount, 0)
}

/*PlayerEnd is called after every ball for the player is over*/
func (p *hatTrick) PlayerEnd(playerID int) {

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *hatTrick) PlayerFinish(playerID int) {
	//Award for any hat tricks.
	htcount := getPlayerStat(game.CurrentPlayer, hatTrickCount)
	if htcount > 0 {
		game.PlaySound(sndGoal)
		log.Infof("HatTrick count was %d", htcount)
		game.AddScore(htcount * 22000) //22000 for each hat trick.
	}
}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *hatTrick) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *hatTrick) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *hatTrick) GameStart() {

}
