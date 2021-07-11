//JAF TODO: This can go away.

package main

import (
	"sync"
	"time"

	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type shotObserver struct {
	shots               []shot //keep track of what shots are lit on the playfield
	mikebossyLights     []int  //This is a list of the MIKE BOSSY lampIDs
	completedIndicators []int  //the two white spots at the bottom of the playfield that symbolize how many times Mike Bossy was spelled
	shotDo              chan int
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

	s.completedIndicators = []int{lmpLeftCompleteLetters, lmpRightCompleteLetters}

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
	case swTopMiddleLane:
	case swTopRightLane:
	case swGoalie:
		s.decShotCount() //decreaseShotCount
		s.setMikeBossyLetters()
		//game.LampOff(s.mikebossyLights[s.getShotCount()]) //lower the number of MIKEBOSSY letters too (to match)
		s.lightNextShot()   //Light one of the shots
		game.AddScore(1000) //give 1000 though for SOG
		game.PlaySound(sndPuckBounce)
		return
	case swLeftTarget:
	case swLeftPointLane:
	case swUpperRightTarget:
	case swMiddleRightTarget:
	case swLowerRightTarget:
	case swBehindGoalLane:
	default:
		return
	}

	log.Infof("Checking if SwitchID was hit %d", sw.SwitchID)
	wasHit := false
	for i, shot := range s.shots {
		if shot.switchID == sw.SwitchID {
			if !shot.wasHit {
				log.Infoln("shotObserver:Lit shot was hit")
				game.LampOff(shot.lampID)
				s.shots[i].wasHit = true
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
	log.Infof("Shot Count incremented: %d\n", s.getShotCount())

	if !s.completedLetters() {
		s.setMikeBossyLetters()
		game.PlaySound(sndLitPuck)
	}
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

	totalLights := s.getShotCount()

	toLight := 9 //default to lighting them all
	if totalLights < 27 {
		toLight = totalLights % 9 //since we could have lit up the MikeBossy lights up to 3 times.
	}

	for i := 0; i < toLight; i++ {
		game.LampOn(s.mikebossyLights[i]) //JAF TODO.. Here. this is where we light up the mike bossy lights in sequence
	}

	//set the completed indicators based on number of times
	game.LampOff(s.completedIndicators...)
	for i := 0; i < totalLights/9 && i < 2; i++ {
		game.LampOn(s.completedIndicators[i])
	}
}

func (s *shotObserver) lightNextShot() {
	for i, shot := range s.shots {
		if shot.wasHit {
			//probably should randomize this, but for now this is ok
			game.LampOn(shot.lampID)
			s.shots[i].wasHit = false
			break
		}
	}
}

func (s *shotObserver) completedLetters() bool {
	//flash the letters, then turn off, lighting the next White circle
	totalCount := s.getShotCount()

	if totalCount%9 == 0 {
		go func() {
			game.PlaySound(sndRaRa)
			s.clearShotStates()                      //relight all of the targets
			game.LampFastBlink(s.mikebossyLights...) //JAF TODO trying this
			time.Sleep(3 * time.Second)
			//let the default logic handle what letters to light now
			s.setMikeBossyLetters()

		}()

		return true
	}

	return false
}

func (s *shotObserver) BallDrained() {
	s.clearShotStates()
}

func (s *shotObserver) PlayerUp(playerID int) {
	setPlayerStat(game.CurrentPlayer, bipShotCount, 0)
	s.setMikeBossyLetters()
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (s *shotObserver) PlayerStart(playerID int) {
	setPlayerStat(game.CurrentPlayer, bipShotCount, 0)
	setPlayerStat(game.CurrentPlayer, totalShotCount, 0)

	game.LampOff(s.mikebossyLights...)
	game.LampOff(s.completedIndicators...)
}

/*PlayerEnd is called after every ball for the player is over*/
func (s *shotObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
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
