package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jfleitz/goflip"
)

type puckChase struct {
	keepChoosingPuck bool   //set to true when chooseNextPuck is called. Set to False when the first switch is hit.
	pucks            []puck //keep track of what pucks are lit on the playfield
	mikebossyLights  []int  //This is a list of the MIKE BOSSY lampIDs
	bossyLightCircle []int  //This is a list of the MIKE BOSSY lampIDs, but in the order of a circular sequence
	//puckSwitches    []int
	//puckLights      []int

}

type puck struct {
	switchID int
	lampID   int
	selected bool
}

var _ goflip.Observer = (*puckChase)(nil)

func (p *puckChase) Init() {
	log.Infoln("puckChase:Init called")
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

	p.clearPuckStates()

	p.mikebossyLights = []int{lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterB, lmpLetterO, lmpLetterS1, lmpLetterS2, lmpLetterY,
	}

	p.bossyLightCircle = []int{
		lmpLetterM,
		lmpLetterI, lmpLetterK, lmpLetterE,
		lmpLetterY, lmpLetterS2, lmpLetterS1, lmpLetterO, lmpLetterB,
	}

}

func (p *puckChase) clearPuckStates() {
	log.Infoln("clearPuckStates called")
	for _, ps := range p.pucks {
		ps.selected = false
		game.LampOff(ps.lampID)
	}
}

func (p *puckChase) SwitchHandler(sw goflip.SwitchEvent) {

	if sw.Pressed == false {
		if sw.SwitchID == swShooterLane { //If the ball just got launched in the shooter lane, then choose the next puck
			go p.chooseNextPuck()
		}
		return
	}

	//some switch was hit, so we stop puck choosing:
	p.keepChoosingPuck = false

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
	case choosePuck:
		go p.chooseNextPuck()
		return
	default:
		return
	}

	var wasHit = false
	for _, pk := range p.pucks {
		if pk.switchID == sw.SwitchID && pk.selected {
			log.Infoln("puckChase:Lit puck hit")
			game.LampOff(pk.lampID)
			wasHit = true
			break
		}
	}

	if !wasHit {
		return
	}

	p.incPuckCount()
	game.AddScore(1000)

	go p.chooseNextPuck()
}

func (p *puckChase) incPuckCount() {
	incPlayerStat(game.CurrentPlayer, bipPuckCount)
	incPlayerStat(game.CurrentPlayer, totalPuckCount)
}

func (p *puckChase) getPuckCount() int {
	return getPlayerStat(game.CurrentPlayer, bipPuckCount)
}

//Spin all of the letters, and then choose a spot on the playfield.
func (p *puckChase) chooseNextPuck() {
	p.keepChoosingPuck = true //since just called.
	//pick a random number between 1-10
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	nextPuck := rnd.Intn(len(p.pucks) - 1) //0 based

	log.Infoln("chooseNextPuck called, next puck will be ", nextPuck)

	//turn off all lit pucks
	p.clearPuckStates()

	//turn off mike bossy lights
	game.LampOff(p.mikebossyLights...)

	//circle around Mike Bossy
	for p.keepChoosingPuck {
		for i := 0; i < 2; i++ {
			for _, l := range p.bossyLightCircle {
				game.LampOn(l)
				time.Sleep(100 * time.Millisecond)
				game.LampOff(l)
			}
		}
	}

	//set the Mike Bossy lights and then light the next puck.
	for i := 0; i < p.getPuckCount() && i < 10; i++ {
		game.LampOn(p.mikebossyLights[i])
	}

	newpuck := p.pucks[nextPuck]
	newpuck.selected = true
	game.LampOn(newpuck.lampID)
}

func (p *puckChase) BallDrained() {
	p.clearPuckStates()
}

func (p *puckChase) PlayerUp(playerID int) {
	//go p.chooseNextPuck()
}

func (p *puckChase) PlayerEnd(playerID int) {

}

func (p *puckChase) GameOver() {

}

func (p *puckChase) GameStart() {

}

func (p *puckChase) PlayerAdded(playerID int) {

}
