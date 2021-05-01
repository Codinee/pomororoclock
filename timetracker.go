package pomodoroclock

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type timer struct {
	start            time.Time
	end              time.Time
	worktime         int
	breaktime        int
	laststartTime    time.Time
	nextprocesstime  time.Time
	nextState        state
	pause            bool
	pauseTime        time.Time
	nonpauseduration time.Duration
	pausepool        []chan interface{}
}

var defaultTime time.Time

type state string

const (
	workstate  state = "work"
	breakstate       = "break"
)

//Set start time to 8:30 AM, end time to 6:30 PM, work time to 25 mins and break time to 5 mins
func getnewtimer() *timer {
	//Set start time to 8:30 AM
	hr := 8
	min := 30
	now := time.Now()
	year := now.Year()
	mon := now.Month()
	day := now.Day()
	sttime := time.Date(year, mon, day, hr, min, 0, 0, time.Local)

	//end time
	hr = 18
	edtime := time.Date(year, mon, day, hr, min, 0, 0, time.Local)

	var tmpstate state
	tmpstate = breakstate

	//processtime
	var processtime time.Time
	if now.After(sttime) {
		processtime = now
	} else {
		processtime = sttime
	}

	return &timer{start: sttime, end: edtime, worktime: 25, breaktime: 5, nextState: tmpstate, nextprocesstime: processtime}
}

//Update timer with  user given start time
func (timetracker *timer) Start(starttime string) {
	s := strings.Split(starttime, ":")
	hr, err := strconv.Atoi(s[0])
	if err != nil {
		fmt.Println("Time given is not valid")
		return
	}
	var min int
	if len(s) == 1 {
		min = 0
	} else {
		min, err = strconv.Atoi(s[1])
		if err != nil {
			fmt.Println("Time given is not valid")
			return
		}
	}
	now := time.Now()
	year := now.Year()
	mon := now.Month()
	day := now.Day()

	sttime := time.Date(year, mon, day, hr, min, 0, 0, time.Local)

	timetracker.start = sttime

	if now.After(sttime) {
		timetracker.nextprocesstime = now
	} else {
		timetracker.nextprocesstime = sttime
	}
	fmt.Println("Start time set to " + s[0] + ":" + s[1])
}

//Update timer with user given end time
func (timetracker *timer) End(endtime string) {
	e := strings.Split(endtime, ":")
	hr, err := strconv.Atoi(e[0])
	if err != nil {
		fmt.Println("Time given is not valid")
		return
	}
	var min int
	if len(e) == 1 {
		min = 0
	} else {
		min, err = strconv.Atoi(e[1])
		if err != nil {
			fmt.Println("Time given is not valid")
			return
		}
	}

	now := time.Now()
	year := now.Year()
	mon := now.Month()
	day := now.Day()
	edtime := time.Date(year, mon, day, hr, min, 0, 0, time.Local)
	timetracker.end = edtime
	fmt.Println("End time set to " + e[0] + ":" + e[1])
}

//Update time with user given work duration
func (timetracker *timer) Work(interval string) {
	wt, err := strconv.Atoi(interval)
	if err != nil {
		fmt.Println("Duration given is not valid")
		return
	}
	timetracker.worktime = wt
	fmt.Println("work duration is set to ", timetracker.worktime)
}

//Update time with user given break duration
func (timetracker *timer) Break(breaktime string) {
	bt, err := strconv.Atoi(breaktime)
	if err != nil {
		fmt.Println("Duration given is not valid")
		return
	}
	timetracker.breaktime = bt
	fmt.Println("break time set to ", timetracker.breaktime)
}

func (timetracker *timer) Pause() {

	if timetracker.laststartTime == defaultTime {
		fmt.Println("No process is running that can be paused")
		return
	}
	timetracker.pause = true
	timetracker.pauseTime = time.Now()

	for _, pauseflag := range timetracker.pausepool {
		close(pauseflag)
	}
	timetracker.pausepool = nil
	fmt.Printf("The %s is paused", timetracker.nextState)
	fmt.Println()
}

//Run the timer
func (timetracker *timer) Run() {

	for time.Now().Before(timetracker.end) {
		if time.Now().After(timetracker.nextprocesstime) && !timetracker.pause {
			for _, pauseflag := range timetracker.pausepool {
				close(pauseflag)
			}
			timetracker.nonpauseduration = 0
			timetracker.pausepool = nil
			timetracker.changestate()
			pauseflag := make(chan interface{})

			timetracker.laststartTime = time.Now()
			timetracker.pausepool = append(timetracker.pausepool, pauseflag)

			var duration time.Duration
			if timetracker.nextState == workstate {
				duration = time.Duration(timetracker.worktime) * time.Minute
			} else {
				duration = time.Duration(timetracker.breaktime) * time.Minute
			}
			timetracker.nextprocesstime = timetracker.laststartTime.Add(duration)

			go process(duration, pauseflag)

		}
	}
	//Set start and end time to next day
	if time.Now().Weekday() == 0 || time.Now().Weekday() == 6 {
		timetracker.start = timetracker.start.Add(72 * time.Hour)
		timetracker.end = timetracker.end.Add(72 * time.Hour)
	} else {
		timetracker.start = timetracker.start.Add(24 * time.Hour)
		timetracker.end = timetracker.end.Add(24 * time.Hour)
	}
	timetracker.laststartTime = defaultTime
}

//Continue the operation break/work for remaining duration on Pause
func (timetracker *timer) Continue() {

	if !timetracker.pause {
		fmt.Println("No process is Paused to continue")
		return
	}
	for _, pauseflag := range timetracker.pausepool {
		close(pauseflag)
	}
	fmt.Println("Continued")
	timetracker.pausepool = nil

	d := timetracker.pauseTime.Sub(timetracker.laststartTime)
	timetracker.nonpauseduration += d

	timetracker.laststartTime = time.Now()

	var duration int

	if timetracker.nextState == workstate {
		duration = timetracker.worktime
	} else {
		duration = timetracker.breaktime
	}

	newduration := time.Duration(duration)*time.Minute - timetracker.nonpauseduration
	timetracker.nextprocesstime = timetracker.laststartTime.Add(newduration)
	pauseflag := make(chan interface{})

	timetracker.pausepool = append(timetracker.pausepool, pauseflag)
	timetracker.pause = false
	go process(newduration, pauseflag)

}

//Restarts the current opetion work/break from zero
func (timetracker *timer) Restart() {

	if !timetracker.pause {
		fmt.Println("No process is Paused to restart")
		return
	}
	fmt.Println("Restarted the process")
	timetracker.laststartTime = time.Now()
	timetracker.nonpauseduration = 0
	var duration time.Duration
	if timetracker.nextState == workstate {
		duration = time.Duration(timetracker.worktime) * time.Minute
	} else {
		duration = time.Duration(timetracker.breaktime) * time.Minute
	}
	timetracker.nextprocesstime = timetracker.laststartTime.Add(duration)
	timetracker.pause = false

	for _, pauseflag := range timetracker.pausepool {
		close(pauseflag)
	}
	timetracker.pausepool = nil

	pauseflag := make(chan interface{})

	timetracker.pausepool = append(timetracker.pausepool, pauseflag)
	go process(duration, pauseflag)
}

//Wait for work/break time and raise beep
func process(sleeptime time.Duration, pauseflag <-chan interface{}) {

	waittime(sleeptime)

	select {
	case <-pauseflag:
		return
	default:
		raisebeep()
	}
}

func raisebeep() {
	beep()
}

func waittime(sleeptime time.Duration) {
	d := sleeptime - time.Duration(10)*time.Second
	time.Sleep(d)

}

func (timetracker *timer) changestate() {
	if timetracker.nextState == workstate {
		timetracker.nextState = breakstate
	} else {
		timetracker.nextState = workstate
	}
}

func (timetracker *timer) Detail() {
	if timetracker.pause {
		fmt.Printf("The current state  -  %s is paused", timetracker.nextState)
		fmt.Println()
	} else {
		if timetracker.laststartTime == defaultTime {
			fmt.Println("No process is running")
		} else {
			fmt.Println("The current state  - ", timetracker.nextState)
			fmt.Println("In the current state from - ", timetracker.laststartTime)
			fmt.Println("Will be in the current state till - ", timetracker.nextprocesstime)
		}
	}
}

func (timetracker *timer) Help() {
	fmt.Println("Pomorodoclock works on these commands - ")
	fmt.Println("start - to set start time in format hh:mm eg. start(13:30) ")
	fmt.Println("end - to set end time in format hh:mm eg. end(18:30)")
	fmt.Println("work- to set work duration in minutes eg. work(30)")
	fmt.Println("break - to set break duration in minutes eg. break(10)")
	fmt.Println("run - to begin work and break tracking")
	fmt.Println("pause - to hold current state")
	fmt.Println("continue - continue tracking from current state")
	fmt.Println("restart - restart tracking from beginning of current state")
	fmt.Println("detail - get current state and time elapsed in current state")
	fmt.Println("quit - to exit the process")
}
