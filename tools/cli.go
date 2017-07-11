package main

import (
	"github.com/urfave/cli"
	"github.com/venkat1109/cadence-codelab/tools/lib"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "cadence"
	app.Usage = "A command-line tool for cadence users"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   lib.FlagAddressWithAlias,
			Value:  "127.0.0.1:7933",
			Usage:  "host:port for cadence frontend service",
			EnvVar: "CADENCE_CLI_ADDRESS",
		},
		cli.StringFlag{
			Name:   lib.FlagDomainWithAlias,
			Usage:  "cadence workflow domain",
			EnvVar: "CADENCE_CLI_DOMAIN",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "register",
			Aliases: []string{"re"},
			Usage:   "Register workflow domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagDescriptionWithAlias,
					Usage: "Domain description",
				},
				cli.StringFlag{
					Name:  lib.FlagOwnerEmailWithAlias,
					Usage: "Owner email",
				},
				cli.StringFlag{
					Name:  lib.FlagRetentionDaysWithAlias,
					Usage: "Workflow execution retention in days",
				},
				cli.BoolFlag{
					Name:  lib.FlagEmitMetricWithAlias,
					Usage: "Flag to emit metric",
				},
			},
			Action: func(c *cli.Context) {
				lib.RegisterDomain(c)
			},
		},
		{
			Name:    "update",
			Aliases: []string{"up", "u"},
			Usage:   "Update existing workflow domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagDescriptionWithAlias,
					Usage: "Domain description",
				},
				cli.StringFlag{
					Name:  lib.FlagOwnerEmailWithAlias,
					Usage: "Owner email",
				},
				cli.StringFlag{
					Name:  lib.FlagRetentionDaysWithAlias,
					Usage: "Workflow execution retention in days",
				},
				cli.BoolFlag{
					Name:  lib.FlagEmitMetricWithAlias,
					Usage: "Flag to emit metric",
				},
			},
			Action: func(c *cli.Context) {
				lib.UpdateDomain(c)
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe existing workflow domain",
			Action: func(c *cli.Context) {
				lib.DescribeDomain(c)
			},
		},
		{
			Name:  "show",
			Usage: "show workflow history",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagRunIDWithAlias,
					Usage: "RunID",
				},
				cli.BoolFlag{
					Name:  lib.FlagPrintRawTimeWithAlias,
					Usage: "Print raw time stamp",
				},
			},
			Action: func(c *cli.Context) {
				lib.ShowHistory(c)
			},
		},
		{
			Name:  "start",
			Usage: "start a new workflow execution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagTaskListWithAlias,
					Usage: "TaskList",
				},
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagWorkflowTypeWithAlias,
					Usage: "WorkflowTypeName",
				},
				cli.IntFlag{
					Name:  lib.FlagExecutionTimeoutWithAlias,
					Usage: "Execution start to close timeout in seconds",
				},
				cli.IntFlag{
					Name:  lib.FlagDecisionTimeoutWithAlias,
					Usage: "Decision task start to close timeout in seconds",
				},
				cli.StringFlag{
					Name:  lib.FlagInputWithAlias,
					Usage: "Input data for the workflow",
				},
			},
			Action: func(c *cli.Context) {
				lib.StartWorkflow(c)
			},
		},
		{
			Name:    "cancel",
			Aliases: []string{"c"},
			Usage:   "cancel a workflow execution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagRunIDWithAlias,
					Usage: "RunID",
				},
			},
			Action: func(c *cli.Context) {
				lib.CancelWorkflow(c)
			},
		},
		{
			Name:    "signal",
			Aliases: []string{"s"},
			Usage:   "signal a workflow execution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagRunIDWithAlias,
					Usage: "RunID",
				},
				cli.StringFlag{
					Name:  lib.FlagNameWithAlias,
					Usage: "SignalName",
				},
				cli.StringFlag{
					Name:  lib.FlagInputWithAlias,
					Usage: "Input message assosciated with signal",
				},
			},
			Action: func(c *cli.Context) {
				lib.SignalWorkflow(c)
			},
		},
		{
			Name:    "terminate",
			Aliases: []string{"term"},
			Usage:   "terminate a new workflow execution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagRunIDWithAlias,
					Usage: "RunID",
				},
				cli.StringFlag{
					Name:  lib.FlagReasonWithAlias,
					Usage: "The reason you want to terminate the workflow",
				},
			},
			Action: func(c *cli.Context) {
				lib.TerminateWorkflow(c)
			},
		},
		{
			Name:  "list",
			Usage: "list open or closed workflow executions",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  lib.FlagOpenWithAlias,
					Usage: "list for open workflow executions, default is to list for closed ones",
				},
				cli.IntFlag{
					Name:  lib.FlagPageSizeWithAlias,
					Value: 10,
					Usage: "Result page size",
				},
				cli.StringFlag{
					Name:  lib.FlagEarliestTimeWithAlias,
					Usage: "EarliestTime of start time, supported formats are '2006-01-02T15:04:05Z07:00' and raw UnixNano",
				},
				cli.StringFlag{
					Name:  lib.FlagLatestTimeWithAlias,
					Usage: "LatestTime of start time, supported formats are '2006-01-02T15:04:05Z07:00' and raw UnixNano",
				},
				cli.StringFlag{
					Name:  lib.FlagWorkflowIDWithAlias,
					Usage: "WorkflowID",
				},
				cli.StringFlag{
					Name:  lib.FlagWorkflowTypeWithAlias,
					Usage: "WorkflowTypeName",
				},
				cli.BoolFlag{
					Name:  lib.FlagPrintRawTimeWithAlias,
					Usage: "Print raw time stamp",
				},
			},
			Action: func(c *cli.Context) {
				lib.QueryWorkflow(c)
			},
		},
	}

	app.Run(os.Args)
}
