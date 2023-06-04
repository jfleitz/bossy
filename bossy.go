package main

/*
Arguments to support:

* dumping the config file struct in a pretty json to confirm settings
forcing the goalie to a specified position

These are for command debugging from the machine:
play sound
pulse solenoid
turn off all solenoids
turn off all lamps
flash lamp
clear displays
show on display


*/

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	observer "github.com/jfleitz/bossy/observers"
	"github.com/jfleitz/bossy/utils"
	"github.com/jfleitz/goflip/pkg/goflip"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	ParseConfig bool `short:"c" long:"config" description:"Parse the config file and dump to console"`
	Goalie      int  `short:"g" long:"goalie" description:"move the goalie to the position passed" default:"-1"`
	Discover    bool `short:"d" long:"discover" description:"Scans for peripherals that are connected and writes to config"`
}

func init() {
	err := utils.LoadConfiguration("config.toml")

	if err != nil {
		os.Exit(1)
	}

	if lvl, err := log.ParseLevel(utils.Settings().LogLevel); err != nil {
		fmt.Printf("Error with log level in config: %v", err)
	} else {
		log.SetLevel(lvl)
	}

	log.SetOutput(os.Stdout)

}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	log.Debugf("Command Args passed: %v\n ", opts)

	if opts.ParseConfig {
		out, err := json.MarshalIndent(utils.Settings(), "", "   ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Parsed Config options:\n")
		fmt.Print(string(out), "\n")

		return
	}

	game := goflip.GetMachine()

	game.ConsoleMode = utils.Settings().ConsoleMode

	game.Observers = []goflip.Observer{
		new(observer.BossyObserver),
		new(observer.GoalObserver),
		//new(warmUpPeriodObserver),
		new(observer.ShotObserver),
		new(observer.EndOfBallBonus),
		new(observer.AttractMode),
		//new(collectOvertime),
		//new(overTimeObserver),
		new(observer.GoalieObserver),
	}

	game.DiagObserver = new(observer.DiagObserver)

	//Set the game limitations here
	game.TotalBalls = utils.Settings().TotalBalls
	game.MaxPlayers = utils.Settings().MaxPlayers
	game.Credits = 0

	//Set Goalie Servo control params
	game.PWMPortConfig.ArcRange = utils.Settings().Goalie.ArcRange
	game.PWMPortConfig.PulseMax = utils.Settings().Goalie.PulseMax
	game.PWMPortConfig.PulseMin = utils.Settings().Goalie.PulseMin
	game.PWMPortConfig.DeviceAddress = utils.Settings().Goalie.DeviceAddress

	log.Debugf("Main, warmupseconds is: %v\n", utils.Settings().WarmupPeriodTimeSeconds)

	//set the game initalizations here
	game.BallInPlay = 0
	game.CurrentPlayer = 0
	utils.InitStats()
	if !game.Init(switchHandler) {
		log.Errorln("Error in initializing goFlip package")
		return
	}

	//see if we are just testing the goalie
	if opts.Goalie >= 0 {
		fmt.Printf("Moving goalie to %v, and sleeping for 2 seconds", opts.Goalie)
		goflip.ServoAngle(opts.Goalie)
		time.Sleep(time.Second * 2)
		return
	}

	//go ahead and go to GameOver by default
	goflip.ChangeGameState(goflip.GameEnded)
	//reader := bufio.NewReader(os.Stdin)
	//TODO defer / send all device disconnects, cleanup etc.
	for {
		time.Sleep(1000 * time.Millisecond) //just keep sleeping
		goflip.SendStats()
	}
}

func switchHandler(sw goflip.SwitchEvent) {

	if !sw.Pressed {
		return
	}

	log.Debugf("Bossy switchHandler. Receivied SwitchID=%d Pressed=%v\n", sw.SwitchID, sw.Pressed)

	if goflip.GetGameState() == goflip.GameEnded {
		//only care about switches that matter when a game is not running
		switch sw.SwitchID {
		case observer.SwSaucer:
			goflip.SolenoidFire(observer.SolSaucer)
		case observer.SwCredit:
			//start game
			go creditControl()
		}

		return
	}
	game := goflip.GetMachine()

	switch sw.SwitchID {
	case observer.SwOuthole:
		//ball over
		log.Debugln("outhole: switch pressed")
		goflip.BallDrained()

	case observer.SwCredit:
		go creditControl()
	case observer.SwShooterLane:
		//		game.PlaySound(sndShooter) //play elsewhere
	case observer.SwTest:
	case observer.SwCoin:
		game.Credits++
		goflip.SetCreditDisp(int8(game.Credits))
	case observer.SwInnerRightLane:
		addThousands(1000)
	case observer.SwMiddleRightLane:
		addThousands(1000)
	case observer.SwOuterRightLane:
		addThousands(5000)
	case observer.SwOuterLeftLane:
		addThousands(5000)
	case observer.SwMiddleLeftLane:
		addThousands(1000)
	case observer.SwInnerLeftLane:
		addThousands(1000)
	case observer.SwRightSlingshot:
		goflip.SolenoidFire(observer.SolRightSlingshot)
		goflip.AddScore(100)
		goflip.PlaySound(observer.SndSlingshot)
	case observer.SwLeftSlingshot:
		goflip.SolenoidFire(observer.SolLeftSlingshot)
		goflip.AddScore(100)
		goflip.PlaySound(observer.SndSlingshot)
	case observer.SwLowerRightTarget:
		goflip.AddScore(1000)
		goflip.PlaySound(observer.SndTargets)
	case observer.SwMiddleRightTarget:
		goflip.AddScore(300)
		goflip.PlaySound(observer.SndTargets)
	case observer.SwUpperRightTarget:
		goflip.AddScore(1000)
		goflip.PlaySound(observer.SndTargets)
	case observer.SwSaucer:
		addHundreds(300)
		go saucerControl()
	case observer.SwLeftPointLane:
		addThousands(1000)
	case observer.SwLeftTarget:
		goflip.AddScore(300)
		goflip.PlaySound(observer.SndTargets)
	case observer.SwLeftBumper:
		goflip.SolenoidOnDuration(observer.SolLeftBumper, 4)
		goflip.AddScore(100)
		goflip.PlaySound(observer.SndBumper)
	case observer.SwRightBumper:
		goflip.SolenoidOnDuration(observer.SolRightBumper, 4)
		goflip.AddScore(100)
		goflip.PlaySound(observer.SndBumper)
	case observer.SwBehindGoalLane:
		addThousands(1000)
	case observer.SwGoalie:
		//game.AddScore(1000) handled by shotObserver
	case observer.SwTopLeftLane:
		//game.LampOn(lmpTopLeftLane)
		goflip.AddScore(300)
		goflip.PlaySound(observer.SndTopLane)
	case observer.SwTopMiddleLane:
		goflip.AddScore(300)
		//game.LampOn(lmpTopMiddleLane)
		goflip.PlaySound(observer.SndTopLane)
	case observer.SwTopRightLane:
		goflip.AddScore(300)
		//game.LampOn(lmpTopRightLane)
		goflip.PlaySound(observer.SndTopLane)
	case observer.SwTargetG:
	case observer.SwTargetO:
	case observer.SwTargetA:
	case observer.SwTargetL:
	}
}

func saucerControl() {
	//JAF TODO: If BossyBonusFromGoal is set, then start the timer to hit a goal and then collect the bonus
	go func() {
		goflip.PlaySound(observer.SndSaucer)
		time.Sleep(2 * time.Second)
		game := goflip.GetMachine()
		totalLetters := utils.GetPlayerStat(game.CurrentPlayer, observer.BipShotCount)

		addThousands(totalLetters * 1000)

		if totalLetters == 0 {
			time.Sleep(1 * time.Second)
		}

		goflip.SolenoidFire(observer.SolSaucer)
	}()
}

func creditControl() {
	game := goflip.GetMachine()

	if goflip.GetGameState() == goflip.GameEnded {
		goflip.ChangeGameState(goflip.InProgress)
		goflip.AddPlayer() // go ahead and add player 1
		goflip.ChangePlayerState(goflip.UpPlayer)
	} else {
		if game.BallInPlay == 1 {
			goflip.AddPlayer()
		}
	}
}

// Incremental scoring with sounds..
func addHundreds(points int) {
	go func() {
		for i := 1; i <= points/100; i++ {
			goflip.AddScore(100)
			goflip.PlaySound(observer.Snd100Points)
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func addThousands(points int) {
	go func() {
		for i := 1; i <= points/1000; i++ {
			goflip.AddScore(1000)
			goflip.PlaySound(observer.Snd1000Points)
			time.Sleep(250 * time.Millisecond)
		}
	}()
}
