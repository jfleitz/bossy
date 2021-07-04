//JAF TODO: This can go away.

package main

import (
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type shotObserver struct {
	shots           []shot //keep track of what shots are lit on the playfield
	mikebossyLights []int  //This is a list of the MIKE BOSSY lampIDs
	shotDo          chan int
}

type shot struct {
	switchID int
	lampID   int
	wasHit   bool
}

var _ goflip.Observer = (*shotObserver)(nil)

func (s *shotObserver) Init() {
	log.Infoln("shotObserver:Init called")

	s.shotDo = make(chan int, 1)

	s.shots = []shot{
		{switchID: swTopLeftLane, lampID: lmpTopLeftLane},
		{switchID: swTopMiddleLane, lampID: lmpTopMiddleLane},
		{switchID: swTopRightLane, lampID: lmpTopRightLane},
		//{switchID: swGoalie, lampID: lmpGoalieWhiteSpot},	//using this to flash that you hit the goalie
		{switchID: swLeftTarget, lampID: lmpLeftTarget},
		{switchID: swLeftPointLane, lampID: lmpPointLaneWhiteSpot},
		{switchID: swUpperRightTarget, lampID: lmpTopTargetWhiteSpot},
		{switchID: swMiddleRightTarget, lampID: lmpMiddleRightTarget},
		{switchID: swLowerRightTarget, lampID: lmpBottomTargetWhiteSpot},
		{switchID: swBehindGoalLane, lampID: lmpOvertimeLeftOfGoal},
	}

	s.mikebossyLights = []int{lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterB, lmpLetterO, lmpLetterS1, lmpLetterS2, lmpLetterY,
	}

	s.clearShotStates()
}

//clearShotStates clears the state of the shot switches (and turns on the lamps) around
//the playfield. This does not do anything with the MIKEBOSSY lights
func (s *shotObserver) clearShotStates() {
	log.Infoln("clearShotStates called")
	for i := 0; i < len(s.shots); i++ {
		s.shots[i].wasHit = false
		game.LampOn(s.shots[i].lampID)
	}
}

func (s *shotObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if sw.Pressed == false {
		return
	}

	//hard coding the switch statement here to be faster.
	switch sw.SwitchID {
	case swTopLeftLane:
		break
	case swTopMiddleLane:
		break
	case swTopRightLane:
		break
	case swGoalie:
		s.decShotCount()                                  //decreaseShotCount
		game.LampOff(s.mikebossyLights[s.getShotCount()]) //lower the number of MIKEBOSSY letters too (to match)
		s.lightNextShot()                                 //Light one of the shots
		game.AddScore(1000)                               //give 1000 though for SOG
		game.PlaySound(sndPuckBounce)
		return
		break
	case swLeftTarget:
		break
	case swLeftPointLane:
		break
	case swUpperRightTarget:
		break
	case swMiddleRightTarget:
		break
	case swLowerRightTarget:
		break
	case swBehindGoalLane:
		break
	default:
		return
	}

	log.Infof("Checking if SwitchID was hit %d", sw.SwitchID)
	wasHit := false
	for _, shot := range s.shots {
		if shot.switchID == sw.SwitchID {
			if !shot.wasHit {
				log.Infoln("shotObserver:Lit shot was hit")
				game.LampOff(shot.lampID)
				shot.wasHit = true
				wasHit = true
			}
			break
		}
	}

	if !wasHit {
		log.Infof("SwitchID was not hit %d", sw.SwitchID)
		return
	}
	s.incShotCount()
	game.LampOn(s.mikebossyLights[s.getShotCount()-1])

	game.PlaySound(sndLitPuck)
	//	game.AddScore(5000) (score should have been counted already)

	//s.setMikeBossyLetters()
}

func (s *shotObserver) incShotCount() {
	incPlayerStat(game.CurrentPlayer, bipShotCount)
	incPlayerStat(game.CurrentPlayer, totalShotCount)
}

func (s *shotObserver) decShotCount() {
	decPlayerStat(game.CurrentPlayer, bipShotCount)
	decPlayerStat(game.CurrentPlayer, totalShotCount)
}

func (s *shotObserver) getShotCount() int {
	return getPlayerStat(game.CurrentPlayer, bipShotCount)
}

func (s *shotObserver) setMikeBossyLetters() {
	game.LampOff(s.mikebossyLights...)

	for i := 0; i < s.getShotCount() && i < 9; i++ {
		game.LampOn(s.mikebossyLights[i]) //JAF TODO.. Here. this is where we light up the mike bossy lights in sequence
	}
}

func (s *shotObserver) lightNextShot() {
	for _, shot := range s.shots {
		if shot.wasHit {
			//probably should randomize this, but for now this is ok
			game.LampOn(shot.lampID)
			shot.wasHit = false
			break
		}
	}
}

func (s *shotObserver) BallDrained() {
	s.clearShotStates()
}

func (s *shotObserver) PlayerUp(playerID int) {
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (s *shotObserver) PlayerStart(playerID int) {
	setPlayerStat(game.CurrentPlayer, bipShotCount, 0)
	setPlayerStat(game.CurrentPlayer, totalShotCount, 0)

	game.LampOff(s.mikebossyLights...)
}

/*PlayerEnd is called after every ball for the player is over*/
func (s *shotObserver) PlayerEnd(playerID int) {

}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (s *shotObserver) PlayerFinish(playerID int) {

}

func (s *shotObserver) GameOver() {

}

func (s *shotObserver) GameStart() {

}

func (s *shotObserver) PlayerAdded(playerID int) {

}
