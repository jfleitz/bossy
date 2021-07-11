package main

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

var inWarmUpPeriod bool
var cancelWarmUp bool

/*
 */

type warmUpPeriodObserver struct {
	//add your variables for the observer here
	warmUpComplete bool
	startWarmUp    bool
}

/*the following line should be called to ensure that your observer DOES
implement the goflip.Observer interface:
*/
var _ goflip.Observer = (*warmUpPeriodObserver)(nil)

/*Init is called by goflip when the application is first started (Init). This
is called only once:
*/
func (p *warmUpPeriodObserver) Init() {
	/*using logrus package for logging. Best practice to call logging when
	only necessary and not in routines that are called a lot*/
	log.Infoln("warmUpPeriodObserver:Init called")

}

/*SwitchHandler is called any time a switch event is received by goflip. This
routine must be kept as fast as possible. Make use of go routines when necessary
Any delay in this routine can cause issues with latency
*/
func (p *warmUpPeriodObserver) SwitchHandler(sw goflip.SwitchEvent) {
	//start the warm up period after the ball is launched, and only for the first time
	if sw.SwitchID == swShooterLane &&
		p.startWarmUp {
		if sw.Pressed {
			log.Infoln("warmupPeriod starting after ball launch")
		} else {
			if !inWarmUpPeriod {
				p.startWarmUp = false
				startWarmUpPeriod(10) //change this to config
			}
		}
	}
}

/*BallDrained is called whenever a ball is drained on the playfield (Before PlayerEnd)*/
func (p *warmUpPeriodObserver) BallDrained() {

}

/*PlayerUp is called after the ball is launched from the Ball Trough for the next ball up
playerID is the player that is now up*/
func (p *warmUpPeriodObserver) PlayerUp(playerID int) {

}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *warmUpPeriodObserver) PlayerStart(playerID int) {
	log.Infoln("PlayerUp: startWarmUp true")
	p.startWarmUp = true
}

/*PlayerEnd is called after thet ball for the player is over*/
func (p *warmUpPeriodObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	p.startWarmUp = false
	defer wait.Done()
}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *warmUpPeriodObserver) PlayerFinish(playerID int) {

}

/*PlayerAdded is called after a player is added by the credit button, and after the GameStart event*/
func (p *warmUpPeriodObserver) PlayerAdded(playerID int) {

}

/*GameOver is called after the last player of the last ball is drained, before the game goes
into the GameOver mode*/
func (p *warmUpPeriodObserver) GameOver() {

}

/*GameStart is called whenever a new game is started*/
func (p *warmUpPeriodObserver) GameStart() {
	p.warmUpComplete = false
}

func startWarmUpPeriod(totalSeconds int) {
	if inWarmUpPeriod {
		return
	}

	go func() {
		inWarmUpPeriod = true
		cancelWarmUp = false
		defer func() {
			game.LampOff(lmpSamePlayerShootAgain)
			inWarmUpPeriod = false
			cancelWarmUp = false
			log.Infoln("Warmup Period complete")
		}()

		for elapsedTime := 0; elapsedTime < totalSeconds; elapsedTime++ {
			if elapsedTime > totalSeconds-3 {
				game.LampFastBlink(lmpSamePlayerShootAgain)
			} else if elapsedTime > totalSeconds-6 {
				game.LampSlowBlink(lmpSamePlayerShootAgain)
			} else {
				game.LampOn(lmpSamePlayerShootAgain)
			}

			//JAF TODO, add a sound for ticking here.
			sleepAndCheck(1)
		}
	}()
}

func sleepAndCheck(ts int) {
	for i := 0; i < ts*2; i++ { //looping every half second to give a chance to cancel
		if cancelWarmUp {
			return
		}
		time.Sleep(time.Duration(ts) * time.Millisecond * 500)
	}
}
