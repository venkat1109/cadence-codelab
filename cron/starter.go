package main

import (
	"github.com/venkat1109/cadence-codelab/common"
	"github.com/venkat1109/cadence-codelab/cron/workflow"
	"go.uber.org/cadence"
	"time"
)

func main() {

	runtime := common.NewRuntime()

	workflowOptions := cadence.StartWorkflowOptions{
		TaskList:                        "cron-decider",
		ExecutionStartToCloseTimeout:    24 * time.Hour,
		DecisionTaskStartToCloseTimeout: 20 * time.Minute,
	}

	schedule := &workflow.CronSchedule{
		Count:      5,
		Frequency:  2 * time.Minute,
		Hostgroups: []string{"hostgroup-1", "hostgroup-2"},
	}

	runtime.StartWorkflow(workflowOptions, workflow.Cron, schedule)
}
