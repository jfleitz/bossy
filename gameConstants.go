package main

//switches
const (
	swOuthole           = 13
	swCredit            = 6
	swShooterLane       = 39
	swTest              = 30
	swCoin              = 2
	swInnerRightLane    = 8
	swMiddleRightLane   = 20
	SwOuterRightLane    = 27
	swOuterLeftLane     = 38
	swMiddleLeftLane    = 10
	swInnerLeftLane     = 9
	swRightSlingshot    = 32
	swLeftSlingshot     = 37
	swLowerRightTarget  = 24
	swMiddleRightTarget = 14
	swUpperRightTarget  = 25
	swSaucer            = 16
	swLeftPointLane     = 26
	swLeftTarget        = 28
	swLeftBumper        = 17
	swRightBumper       = 31
	swBehindGoalLane    = 19
	swGoalie            = 18
	swTopLeftLane       = 11
	swTopMiddleLane     = 23
	swTopRightLane      = 22
	swTargetG           = 36
	swTargetO           = 35
	swTargetA           = 34
	swTargetL           = 4

	//TODO need tilt
)

//Lamps
const (
	lmpLetterM                   = 2
	lmpLetterI                   = 3
	lmpLetterK                   = 4
	lmpLetterE                   = 5
	lmpLetterB                   = 6
	lmpLetterO                   = 7
	lmpLetterS1                  = 8
	lmpLetterS2                  = 9
	lmpLetterY                   = 10
	lmpOvertimeLeftOfGoal        = 12
	lmpGoalieWhiteSpot           = 13
	lmpPointLaneWhiteSpot        = 14
	lmpTopTargetWhiteSpot        = 15
	lmpBottomTargetWhiteSpot     = 32
	lmpGoalOnLeftOrangeSpot      = 18
	lmpTargetsOrangeSpot         = 19
	lmpTopOrangeSpot             = 17
	lmpLeftReturnLaneOrangeSpot  = 20
	lmpRightReturnLaneOrangeSpot = 21
	lmpPointLaneSpecial          = 22
	lmpBottomLeftSpecial         = 23
	lmpBottomRightSpecial        = 24
	lmpTopLeftLane               = 25
	lmpTopMiddleLane             = 26
	lmpLeftTarget                = 28
	lmpMiddleRightTarget         = 29
	lmpTopRightLane              = 31
	lmp5000Bonus1                = 33
	lmp5000Bonus2                = 34
	lmp5000Bonus3                = 35
	lmp5000Bonus4                = 36
	lmpBonusSpotRight            = 38
	lmpBottomRightGreenSpot      = 39
	lmpGoalLight                 = 46
	lmpTopLeftStar               = 50
	lmpSamePlayerShootAgain      = 52
	lmpAllLamps                  = 1
	lmpPlayer1                   = 61
	lmpPlayer2                   = 62
	lmpPlayer3                   = 63
	lmpPlayer4                   = 64

	lmpGameOver            = 44
	lmpBackglassSamePlayer = 52
	lmpTilt                = 40
	lmpPeriod              = 53
	lmpMatch               = 38
	lmpHighScore           = 41
)

//Solenoids
const (
	solDropTargets    = 1
	solOuthole        = 7
	solRightSlingshot = 10
	solSaucer         = 11
	solLeftBumper     = 13
	solRightBumper    = 12
	solFlippers       = 18
	solLeftSlingshot  = 8 //Not sure what this is.
)

//stats constants
const (
	bipPuckCount    = "bipPuckCount"    //ball in progress puck count
	totalPuckCount  = "totalPuckCount"  //total puck count for a player
	bipGoalCount    = "bipGoalCount"    //ball in progress goal
	totalGoalCount  = "totalGoalCount"  //total goal count for a player
	leftGoalCount   = "leftGoalCount"   //total left flipper goals (hattrick)
	rightGoalCount  = "rightGoalCount"  //total right flipper goals (hattrack)
	saucerGoalCount = "saucerGoalCount" //total saucer Goal Count (hattrick)
	hatTrickFor     = "hatTrickFor"     //Some logic around this. 0=no goal counted. 1=Left Flipper, 2=Right Flipper, 3=Saucer. Set by the hattrick Observer. Used by the goalObserver.
	hatTrickCount   = "hatTrickCount"   //total number of hat tricks recieved
	otSeconds       = "otSeconds"       //total number of OT seconds collected
)

//Feature constants
const (
	passedToLeft   = 1
	passedToRight  = 2
	passedToSaucer = 3
)

//observer events
const (
	goalScored = -50
	inOvertime = -51
)

//sounds
/*
SND = Sound board control (when ready)
Sounds for Bossy:
0 = High tone (maybe a puck bouncing ? ?) - inlanes
1 = Star spangled banner - after warmup period
2 = up then down tone - like crowd doing the wave? (maybe on ball launch have this?)
3 = icing / asteroids fire -- pop bumpers and sling shots (defense)
4 = reset
5 = whistle, up and down tone, and asteroids icing (end of game)
6 = n/a
7 = whistle only - ball drain, and 2 for ball launch?
8 = low tones (bouncing ball..) - outlanes?
9 = high tones (bouncing ball) ?? - when we are moving the lit puck to a new place
10 = nothing
11 = puck bounce (10 pt shot) --make this counting down timer
12 = ra ra , ra-ra-ra, ra ra, ra-ra-ra (when you add a player)
13 = charge - (add player / credit)
14 = nothing
15 = nothing*/

const (
	sndLitPuck     = 0
	sndAnthem      = 1
	sndBallLaunch  = 2
	sndFiring      = 3
	sndGoal        = 5 //maybe use an actual horn for this
	sndWhistle     = 7 //made this generic as we are going to use this in a few places
	sndOutlane     = 8
	sndTargets     = 9
	sndPuckBounce  = 11
	sndRaRa        = 12
	sndCredit      = 13
	sndTimeSeconds = 11 //for counting down the last 10 seconds of a period
	sndBallDrained = 5
)
