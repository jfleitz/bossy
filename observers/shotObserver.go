//ShotObserver keeps track of the hockey puck shots that are completed around the playfield.
//The goalie will take away one of the lit MIKEBOSSY letters and light a spot somewhere on the game again

package observer

import (
	"sync"
	"time"

	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

type ShotObserver struct {
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

var _ goflip.Observer = (*ShotObserver)(nil)

func (s *ShotObserver) Init() {
	log.Debugln("ShotObserver:Init called")

	s.shotDo = make(chan int, 1)

	s.shots = []shot{
		{switchID: SwTopLeftLane, lampID: LmpTopLeftLane},
		{switchID: SwTopMiddleLane, lampID: LmpTopMiddleLane},
		{switchID: SwTopRightLane, lampID: LmpTopRightLane},
		//{switchID: swGoalie, lampID: lmpGoalieWhiteSpot},	//using this to flash that you hit the goalie
		{switchID: SwLeftTarget, lampID: LmpLeftTarget},
		{switchID: SwLeftPointLane, lampID: LmpPointLaneWhiteSpot},
		{switchID: SwUpperRightTarget, lampID: LmpTopTargetWhiteSpot},
		{switchID: SwMiddleRightTarget, lampID: LmpMiddleRightTarget},
		{switchID: SwLowerRightTarget, lampID: LmpBottomTargetWhiteSpot},
		{switchID: SwBehindGoalLane, lampID: LmpOvertimeLeftOfGoal},
	}

	s.mikebossyLights = []int{LmpLetterM,
		LmpLetterI, LmpLetterK, LmpLetterE,
		LmpLetterB, LmpLetterO, LmpLetterS1, LmpLetterS2, LmpLetterY,
	}

	s.completedIndicators = []int{LmpLeftCompleteLetters, LmpRightCompleteLetters}

	s.clearShotStates()
}

// clearShotStates clears the state of the shot switches (and turns on the lamps) around
// the playfield. This does not do anything with the MIKEBOSSY lights
func (s *ShotObserver) clearShotStates() {
	log.Debugln("clearShotStates called")
	for i := 0; i < len(s.shots); i++ {
		s.shots[i].wasHit = false
		goflip.LampOn(s.shots[i].lampID)
	}
}

func (s *ShotObserver) SwitchHandler(sw goflip.SwitchEvent) {
	if sw.Pressed == false {
		return
	}

	switch sw.SwitchID {
	case SwTopLeftLane:
	case SwTopMiddleLane:
	case SwTopRightLane:
	case SwGoalie:
		s.decShotCount() //decreaseShotCount
		s.setMikeBossyLetters()
		//game.LampOff(s.mikebossyLights[s.getShotCount()]) //lower the number of MIKEBOSSY letters too (to match)
		s.lightNextShot()     //Light one of the shots
		goflip.AddScore(1000) //give 1000 though for SOG
		goflip.PlaySound(SndGoalie)
		return
	case SwLeftTarget:
	case SwLeftPointLane:
	case SwUpperRightTarget:
	case SwMiddleRightTarget:
	case SwLowerRightTarget:
	case SwBehindGoalLane:
	default:
		return
	}

	log.Debugf("Checking if SwitchID was hit %d", sw.SwitchID)
	wasHit := false
	for i, shot := range s.shots {
		if shot.switchID == sw.SwitchID {
			if !shot.wasHit {
				log.Debugln("ShotObserver:Lit shot was hit")
				goflip.LampOff(shot.lampID)
				s.shots[i].wasHit = true
				wasHit = true
			}
			break
		}
	}

	if !wasHit {
		log.Debugf("SwitchID was not hit %d", sw.SwitchID)
		return
	}
	s.incShotCount()
	log.Debugf("Shot Count incremented: %d\n", s.getShotCount())

	if !s.completedLetters() {
		s.setMikeBossyLetters()
		goflip.PlaySound(SndLetterAdded)
	}
}

func (s *ShotObserver) incShotCount() {
	game := goflip.GetMachine()

	utils.IncPlayerStat(game.CurrentPlayer, BipShotCount)
	utils.IncPlayerStat(game.CurrentPlayer, TotalShotCount)
}

func (s *ShotObserver) decShotCount() {
	game := goflip.GetMachine()

	utils.DecPlayerStat(game.CurrentPlayer, BipShotCount)
	utils.DecPlayerStat(game.CurrentPlayer, TotalShotCount)
}

func (s *ShotObserver) getShotCount() int {
	game := goflip.GetMachine()

	return utils.GetPlayerStat(game.CurrentPlayer, BipShotCount)
}

func (s *ShotObserver) setMikeBossyLetters() {
	goflip.LampOff(s.mikebossyLights...)

	totalLights := s.getShotCount()

	toLight := 9 //default to lighting them all
	if totalLights < 27 {
		toLight = totalLights % 9 //since we could have lit up the MikeBossy lights up to 3 times.
	}

	for i := 0; i < toLight; i++ {
		goflip.LampOn(s.mikebossyLights[i])
	}

	//set the completed indicators based on number of times
	goflip.LampOff(s.completedIndicators...)
	for i := 0; i < totalLights/9 && i < 2; i++ {
		goflip.LampOn(s.completedIndicators[i])
	}
}

func (s *ShotObserver) lightNextShot() {
	for i, shot := range s.shots {
		if shot.wasHit {
			//probably should randomize this, but for now this is ok
			goflip.LampOn(shot.lampID)
			s.shots[i].wasHit = false
			break
		}
	}
}

func (s *ShotObserver) completedLetters() bool {
	//flash the letters, then turn off, light all shots
	totalCount := s.getShotCount()

	if totalCount%9 == 0 { //If there is a multiple of 9 targets that have been lit
		go func() {
			goflip.PlaySound(SndLettersCompleted)
			s.clearShotStates() //relight all of the targets
			goflip.LampFastBlink(s.mikebossyLights...)
			time.Sleep(3 * time.Second)
			//let the default logic handle what letters to light now
			s.setMikeBossyLetters()

		}()

		return true
	}

	return false
}

func (s *ShotObserver) BallDrained() {
	s.clearShotStates()
}

func (s *ShotObserver) PlayerUp(playerID int) {
	s.clearShotStates()
	game := goflip.GetMachine()
	utils.SetPlayerStat(game.CurrentPlayer, BipShotCount, 0)
	s.setMikeBossyLetters()
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (s *ShotObserver) PlayerStart(playerID int) {
	game := goflip.GetMachine()

	utils.SetPlayerStat(game.CurrentPlayer, BipShotCount, 0)
	utils.SetPlayerStat(game.CurrentPlayer, TotalShotCount, 0)

	goflip.LampOff(s.mikebossyLights...)
	goflip.LampOff(s.completedIndicators...)
}

/*PlayerEnd is called after every ball for the player is over*/
func (s *ShotObserver) PlayerEnd(playerID int, wait *sync.WaitGroup) {
	defer wait.Done()
}

/*
PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)
*/
func (s *ShotObserver) PlayerFinish(playerID int) {

}

func (s *ShotObserver) GameOver() {

}

func (s *ShotObserver) GameStart() {

}

func (s *ShotObserver) PlayerAdded(playerID int) {

}
