package observer

// switches
const (
	SwOuthole           = 13
	SwCredit            = 6
	SwShooterLane       = 39
	SwTest              = 30
	SwCoin              = 2
	SwInnerRightLane    = 8
	SwMiddleRightLane   = 20
	SwOuterRightLane    = 27
	SwOuterLeftLane     = 38
	SwMiddleLeftLane    = 10
	SwInnerLeftLane     = 9
	SwRightSlingshot    = 32
	SwLeftSlingshot     = 37
	SwLowerRightTarget  = 24
	SwMiddleRightTarget = 14
	SwUpperRightTarget  = 25
	SwSaucer            = 16
	SwLeftPointLane     = 26
	SwLeftTarget        = 28
	SwLeftBumper        = 17
	SwRightBumper       = 31
	SwBehindGoalLane    = 19
	SwGoalie            = 18
	SwTopLeftLane       = 11
	SwTopMiddleLane     = 23
	SwTopRightLane      = 22
	SwTargetG           = 36
	SwTargetO           = 35
	SwTargetA           = 34
	SwTargetL           = 4
	SwBackGoal          = 33
)

// Lamps
const (
	LmpLetterM                   = 2
	LmpLetterI                   = 3
	LmpLetterK                   = 4
	LmpLetterE                   = 5
	LmpLetterB                   = 6
	LmpLetterO                   = 7
	LmpLetterS1                  = 8
	LmpLetterS2                  = 9
	LmpLetterY                   = 10
	LmpOvertimeLeftOfGoal        = 12
	LmpGoalieWhiteSpot           = 13
	LmpPointLaneWhiteSpot        = 14
	LmpTopTargetWhiteSpot        = 15
	LmpBottomTargetWhiteSpot     = 32
	LmpGoalOnLeftOrangeSpot      = 18
	LmpTargetsOrangeSpot         = 19
	LmpTopOrangeSpot             = 17
	LmpLeftReturnLaneOrangeSpot  = 20
	LmpRightReturnLaneOrangeSpot = 21
	LmpPointLaneSpecial          = 22
	LmpBottomLeftSpecial         = 23
	LmpBottomRightSpecial        = 24
	LmpTopLeftLane               = 25
	LmpTopMiddleLane             = 26
	LmpLeftTarget                = 28
	LmpMiddleRightTarget         = 29
	LmpTopRightLane              = 31
	Lmp5000Bonus1                = 33
	Lmp5000Bonus2                = 34
	Lmp5000Bonus3                = 35
	Lmp5000Bonus4                = 36
	LmpBonusSpotRight            = 38
	LmpBottomRightGreenSpot      = 39
	LmpGoalLight                 = 46
	LmpTopLeftStar               = 50
	LmpSamePlayerShootAgain      = 52
	LmpAllLamps                  = 1
	LmpPlayer1                   = 61
	LmpPlayer2                   = 62
	LmpPlayer3                   = 63
	LmpPlayer4                   = 64
	//JAF Check these
	Lmp25000Bonus           = 45
	LmpRightCompleteLetters = 38
	LmpLeftCompleteLetters  = 37

	LmpGameOver            = 44
	LmpBackglassSamePlayer = 52
	LmpTilt                = 40
	LmpPeriod              = 53
	LmpMatch               = 38
	LmpHighScore           = 41
)

// Solenoids
const (
	SolDropTargets    = 1
	SolOuthole        = 7
	SolRightSlingshot = 10
	SolSaucer         = 11
	SolLeftBumper     = 13
	SolRightBumper    = 12
	SolFlippers       = 18
	SolLeftSlingshot  = 8 //Not sure what this is.
)

// stats constants
const (
	BipShotCount    = "bipShotCount"    //ball in progress puck count
	TotalShotCount  = "totalShotCount"  //total puck count for a player
	BipGoalCount    = "bipGoalCount"    //ball in progress goal
	TotalGoalCount  = "totalGoalCount"  //total goal count for a player
	GoalTargetCount = "goalTargetCount" //number of times the GOAL targets were completed

	OTSeconds = "otSeconds" //total number of OT seconds collected
)

// Feature constants
const (
	PassedToLeft   = 1
	PassedToRight  = 2
	PassedToSaucer = 3
)

// observer events
const (
	GoalScored = -50
	InOvertime = -51
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
	/*	sndLitPuck     = 0
		sndAnthem      = 1 //starting game
		sndBossySaucer = 2 //up down noise
		sndFiring      = 3
		sndWhistle     = 7 //made this generic as we are going to use this in a few places
		sndOutlane     = 8
		sndTargets     = 9
		sndPuckBounce  = 11
		sndRaRa        = 12 //used adding a player
		sndTimeSeconds = 11 //for counting down the last 10 seconds of a period
		sndBallDrained = 5*/

	SndCredit           = 13
	Snd100Points        = 11
	Snd500Points        = 11
	Snd1000Points       = 11
	SndStartGame        = 1
	SndAddPlayer        = 12
	SndSaucer           = 2
	SndShooter          = 3
	SndBumper           = 0
	SndSlingshot        = 7
	SndTopLane          = 9
	SndGoalBonus        = 12
	SndGoalie           = 8
	SndGoalTarget       = 5
	SndTargets          = 9
	SndLitPuck          = 2
	SndLettersCompleted = 12
	SndLetterAdded      = 7
	SndLetterBonus      = 11
	SndGameOver         = 2
	SndWarmUp           = 11
)
