package pomodoroclock

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GettingReady() {

	timetracker := getnewtimer()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("pomodoroclock>")
	active := true
	var inp string

	for active {
		//Get Operation to be performed
		if scanner.Scan() {
			fmt.Println("pomodoroclock>")
			inp = scanner.Text()
		}
		if inp == "quit" {
			active = false
		} else {
			callasrequired(timetracker, inp)
		}
	}
}

func callasrequired(timetracker *timer, inp string) {
	var parm string

	if strings.Contains(inp, "start(") {
		s := strings.Split(inp, "(")
		inp = "start"
		parm = s[1][:len(s[1])-1]
	}

	if strings.Contains(inp, "end") {
		s := strings.Split(inp, "(")
		inp = "end"
		parm = s[1][:len(s[1])-1]
	}

	if strings.Contains(inp, "work") {
		s := strings.Split(inp, "(")
		inp = "work"
		parm = s[1][:len(s[1])-1]
	}

	if strings.Contains(inp, "break") {
		s := strings.Split(inp, "(")
		inp = "break"
		parm = s[1][:len(s[1])-1]
	}

	switch inp {
	case "start":
		timetracker.Start(parm)
	case "end":
		timetracker.End(parm)
	case "work":
		timetracker.Work(parm)
	case "break":
		timetracker.Break(parm)
	case "run":
		go timetracker.Run()
	case "pause":
		timetracker.Pause()
	case "continue":
		go timetracker.Continue()
	case "restart":
		go timetracker.Restart()
	case "detail":
		timetracker.Detail()
	case "help":
		timetracker.Help()
	default:
		fmt.Println("The option specified is not valid. Use 'help' to get command details")
	}
}
