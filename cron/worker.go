package main

import (
	_ "github.com/venkat1109/cadence-codelab/cron/activity"
	_ "github.com/venkat1109/cadence-codelab/cron/workflow"
	//"github.com/venkat1109/cadence-codelab/common"
	//"go.uber.org/cadence"
	"fmt"
)

const (
	decisionTaskList   = "cron-decider"
	hostgroup1TaskList = "hostgroup-1"
	hostgroup2TaskList = "hostgroup-2"
)

func main() {
	fmt.Println("Not implemented")
}
