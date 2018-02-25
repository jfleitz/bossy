package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type puckChase struct {
	pucks            []puck //keep track of what pucks are lit on the playfield
	mikebossyLights  []int  //This is a list of the MIKE BOSSY lampIDs
	bossyLightCircle []int  //This is a list of the MIKE BOSSY lampIDs, but in the order of a circular sequence
	puckDo           chan int
}

type puck struct {
	switchID int
	lampID   int
	selected bool
}

const stopPuck = 0
const nextPuck = 1
const endGame = 2

var _ goflip.Observer = (*puckChase)(nil)

func (p *puckChase) Init() {
	log.Infoln("puckChase:Init called")

	p.puckDo = make(chan int, 1)

	p.pucks = []puck{
		{switchID: swTopLeftLane, lampID: lmpTopLeftLane},
		{switchID: swTopMiddleLane, lampID: lmpTopMiddleLane},
		{switchID: swTopRightLane, lampID: lmpTopRightLane},
		{switchID: swGoalie, lampID: lmpGoalieWhiteSpot},
		{switchID: swLeftTarget, lampID: lmpLeftTarget},
		{switchID: swLeftPointLane, lampID: lmpPointLaneWhiteSpot},
		{switchID: swUpperRightTarget, lampID: lmpTopTargetWhiteSpot},
		{switchID: swMiddleRightTarget, lampID: lmpMiddleRightTarget},
		{switchID: swLowerRightTarget, lampID: lmpBottomTargetWhiteSpot},
		{switchID: swBehindGoalLane, lampID: lmpOvertimeLeftOfGoal},
	}

	p.mikebossyLights = []int{lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterB, lmpLetterO, lmpLetterS1, lmpLetterS2, lmpLetterY,
	}

	p.bossyLightCircle = []int{
		lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterY, lmpLetterS2, lmpLetterS1, lmpLetterO, lmpLetterB,
	}

	p.clearPuckStates()

	go p.puckHandler()

}

//clearPuckStates clears the state of the puck switches (and their lamps) around
//the playfield. This does not do anything with the MIKEBOSSY lights
func (p *puckChase) clearPuckStates() {
	log.Infoln("clearPuckStates called")
	for i := 0; i < len(p.pucks); i++ {
		p.pucks[i].selected = false
		game.LampOff(p.pucks[i].lampID)
	}
}

func (p *puckChase) SwitchHandler(sw goflip.SwitchEvent) {
	if sw.Pressed == false {
		if sw.SwitchID == swShooterLane {
			//If the ball just got launched in the shooter lane,
			// then choose the next puck (even on replunge)
			p.puckDo <- nextPuck
		}
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
	case goalScored:
		p.puckDo <- nextPuck
		return
	default:
		return
	}

	log.Infof("Checking if SwitchID was selected %d", sw.SwitchID)
	wasHit := false
	for _, pk := range p.pucks {
		if pk.switchID == sw.SwitchID && pk.selected {
			log.Infoln("puckChase:Lit puck hit")
			game.LampOff(pk.lampID)
			wasHit = true
			break
		}
	}

	if !wasHit {
		log.Infof("SwitchID was not selected %d", sw.SwitchID)
		return
	}
	log.Infof("SwitchID was  selected %d", sw.SwitchID)
	p.incPuckCount()
	game.PlaySound(sndLitPuck)
	game.AddScore(5000)
	p.puckDo <- nextPuck
}

func (p *puckChase) incPuckCount() {
	incPlayerStat(game.CurrentPlayer, bipPuckCount)
	incPlayerStat(game.CurrentPlayer, totalPuckCount)
}

func (p *puckChase) getPuckCount() int {
	return getPlayerStat(game.CurrentPlayer, bipPuckCount)
}

//Spin all of the letters, and then choose a spot on the playfield.
func (p *puckChase) puckHandler() {
	//p.keepChoosingPuck = true //since just called.
	//pick a random number between 1-10

	curLamp := 0
	lightLoops := 3 //to stop the looping
	litPuck := 0

	for {
		select {
		case newPuck := <-p.puckDo:
			switch newPuck {
			case nextPuck:
				log.Infoln("puckChase choosing next puck")
				seed := rand.NewSource(time.Now().UnixNano())
				rnd := rand.New(seed)
				litPuck = rnd.Intn(len(p.pucks) - 1) //0 based

				p.clearPuckStates()
				//turn off mike bossy lights
				game.LampOff(p.mikebossyLights...)
				lightLoops = 0
				curLamp = 0

			case stopPuck:
				log.Warningln("puckChase:stopPuck called")
			case endGame:
				game.LampOff(p.mikebossyLights...)
				lightLoops = 3
				curLamp = 0
			}

		case <-time.After(100 * time.Millisecond):
			//looping here

			if lightLoops < 2 {
				//we are lighting something
				p.turnOffPrevLamp(curLamp)
				//turn on next
				game.LampOn(p.bossyLightCircle[curLamp])
				curLamp++

				if curLamp > len(p.bossyLightCircle)-1 {
					lightLoops++
					curLamp = 0
				}

				if lightLoops == 2 {
					//we got done with our loops, so display next puck
					p.setNextPuck(litPuck)
				}
			}

		}
	}

}

func (p *puckChase) setNextPuck(nextPuck int) {
	log.Infoln("chooseNextPuck called, next puck will be ", nextPuck)
	game.LampOff(p.mikebossyLights...)

	//set the Mike Bossy lights and then light the next puck.
	for i := 0; i < p.getPuckCount() && i < 9; i++ {
		game.LampOn(p.mikebossyLights[i])
	}

	p.pucks[nextPuck].selected = true
	game.LampOn(p.pucks[nextPuck].lampID)
}

func (p *puckChase) BallDrained() {
	p.clearPuckStates()
}

func (p *puckChase) PlayerUp(playerID int) {
	//go p.chooseNextPuck()
}

/*PlayerStart is called the very first time a player is playing (their first Ball1)
 */
func (p *puckChase) PlayerStart(playerID int) {
	setPlayerStat(game.CurrentPlayer, bipPuckCount, 0)
	setPlayerStat(game.CurrentPlayer, totalPuckCount, 0)

}

/*PlayerEnd is called after every ball for the player is over*/
func (p *puckChase) PlayerEnd(playerID int) {
	p.puckDo <- stopPuck
}

/*PlayerFinish is called after the very last ball for the player is over
(after ball 3 for example)*/
func (p *puckChase) PlayerFinish(playerID int) {

}

func (p *puckChase) GameOver() {
	p.puckDo <- endGame

}

func (p *puckChase) GameStart() {

}

func (p *puckChase) PlayerAdded(playerID int) {

}

func (p *puckChase) turnOffPrevLamp(curLamp int) {
	if curLamp > 0 {
		game.LampOff(p.bossyLightCircle[curLamp-1])
	} else {
		game.LampOff(p.bossyLightCircle[8])
	}
}
