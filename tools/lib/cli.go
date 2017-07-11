package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"github.com/urfave/cli"
	factory "github.com/venkat1109/cadence-codelab/common"
	"go.uber.org/cadence"
	s "go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/common"
	"go.uber.org/cadence/common/util"
)

/**
Flags used to specify cli command line arguments
*/
const (
	FlagAddress                   = "address"
	FlagAddressWithAlias          = FlagAddress + ", ad"
	FlagDomain                    = "domain"
	FlagDomainWithAlias           = FlagDomain + ", do"
	FlagWorkflowID                = "workflow_id"
	FlagWorkflowIDWithAlias       = FlagWorkflowID + ", wid, w"
	FlagRunID                     = "run_id"
	FlagRunIDWithAlias            = FlagRunID + ", rid, r"
	FlagTaskList                  = "tasklist"
	FlagTaskListWithAlias         = FlagTaskList + ", tl"
	FlagWorkflowType              = "workflow_type"
	FlagWorkflowTypeWithAlias     = FlagWorkflowType + ", wt"
	FlagExecutionTimeout          = "execution_timeout"
	FlagExecutionTimeoutWithAlias = FlagExecutionTimeout + ", et"
	FlagDecisionTimeout           = "decision_timeout"
	FlagDecisionTimeoutWithAlias  = FlagDecisionTimeout + ", dt"
	FlagInput                     = "input"
	FlagInputWithAlias            = FlagInput + ", i"
	FlagReason                    = "reason"
	FlagReasonWithAlias           = FlagReason + ", re"
	FlagOpen                      = "open"
	FlagOpenWithAlias             = FlagOpen + ", op"
	FlagPageSize                  = "pagesize"
	FlagPageSizeWithAlias         = FlagPageSize + ", ps"
	FlagEarliestTime              = "earliest_time"
	FlagEarliestTimeWithAlias     = FlagEarliestTime + ", et"
	FlagLatestTime                = "latest_time"
	FlagLatestTimeWithAlias       = FlagLatestTime + ", lt"
	FlagPrintRawTime              = "print_raw_time"
	FlagPrintRawTimeWithAlias     = FlagPrintRawTime + ", prt"
	FlagDescription               = "description"
	FlagDescriptionWithAlias      = FlagDescription + ", desc"
	FlagOwnerEmail                = "owner_email"
	FlagOwnerEmailWithAlias       = FlagOwnerEmail + ", oe"
	FlagRetentionDays             = "retention_days"
	FlagRetentionDaysWithAlias    = FlagRetentionDays + ", rd"
	FlagEmitMetric                = "emit_metric"
	FlagEmitMetricWithAlias       = FlagEmitMetric + ", em"
	FlagName                      = "name"
	FlagNameWithAlias             = FlagName + ", n"
)

const (
	localHostPort = "127.0.0.1:7933"
)

// ExitIfError exit while err is not nil and print the calling stack also
func ExitIfError(err error) {
	const stacksEnv = `CADENCE_CLI_SHOW_STACKS`
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if os.Getenv(stacksEnv) != `` {
			debug.PrintStack()
		} else {
			fmt.Fprintf(os.Stderr, "('export %s=1' to see stack traces)\n", stacksEnv)
		}
		os.Exit(1)
	}
}

// RegisterDomain register a domain
func RegisterDomain(c *cli.Context) {
	domainClient := getDomainClient(c)
	domain := getRequiredGlobalOption(c, FlagDomain)

	description := c.String(FlagDescription)
	ownerEmail := c.String(FlagOwnerEmail)
	retentionDays := c.Int(FlagRetentionDays)
	emitMetric := c.Bool(FlagEmitMetric)
	request := &s.RegisterDomainRequest{
		Name:                                   common.StringPtr(domain),
		Description:                            common.StringPtr(description),
		OwnerEmail:                             common.StringPtr(ownerEmail),
		WorkflowExecutionRetentionPeriodInDays: common.Int32Ptr(int32(retentionDays)),
		EmitMetric:                             common.BoolPtr(emitMetric),
	}

	err := domainClient.Register(request)
	if err != nil {
		if _, ok := err.(*s.DomainAlreadyExistsError); !ok {
			fmt.Printf("Operation failed: %v.\n", err.Error())
		} else {
			fmt.Printf("Domain %s already registered.\n", domain)
		}
	} else {
		fmt.Printf("Domain %s succeesfully registered.\n", domain)
	}
}

// UpdateDomain updates a domain
func UpdateDomain(c *cli.Context) {
	domainClient := getDomainClient(c)
	domain := getRequiredGlobalOption(c, FlagDomain)

	description := c.String(FlagDescription)
	ownerEmail := c.String(FlagOwnerEmail)
	retentionDays := c.Int(FlagRetentionDays)
	emitMetric := c.Bool(FlagEmitMetric)
	info := &s.UpdateDomainInfo{
		Description: common.StringPtr(description),
		OwnerEmail:  common.StringPtr(ownerEmail),
	}
	config := &s.DomainConfiguration{
		WorkflowExecutionRetentionPeriodInDays: common.Int32Ptr(int32(retentionDays)),
		EmitMetric:                             common.BoolPtr(emitMetric),
	}

	err := domainClient.Update(domain, info, config)
	if err != nil {
		if _, ok := err.(*s.EntityNotExistsError); !ok {
			fmt.Printf("Operation failed: %v.\n", err.Error())
		} else {
			fmt.Printf("Domain %s not exists.\n", domain)
		}
	} else {
		fmt.Printf("Domain %s succeesfully updated.\n", domain)
	}
}

// DescribeDomain updates a domain
func DescribeDomain(c *cli.Context) {
	domainClient := getDomainClient(c)
	domain := getRequiredGlobalOption(c, FlagDomain)

	info, config, err := domainClient.Describe(domain)
	if err != nil {
		if _, ok := err.(*s.EntityNotExistsError); !ok {
			fmt.Printf("Operation failed: %v.\n", err.Error())
		} else {
			fmt.Printf("Domain %s not exists.\n", domain)
		}
	} else {
		fmt.Printf("Name:%v, Description:%v, OwnerEmail:%v, Status:%v, RetentionInDays:%v, EmitMetrics:%v\n",
			info.GetName(),
			info.GetDescription(),
			info.GetOwnerEmail(),
			info.GetStatus(),
			config.GetWorkflowExecutionRetentionPeriodInDays(),
			config.GetEmitMetric())
	}
}

// ShowHistory shows the history of given workflow execution based on workflowID and runID.
func ShowHistory(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	wid := getRequiredOption(c, FlagWorkflowID)
	rid := c.String(FlagRunID)
	printRawTime := c.Bool(FlagPrintRawTime)

	history, err := wfClient.GetWorkflowHistory(wid, rid)
	if err != nil {
		ExitIfError(err)
	}

	for _, e := range history.GetEvents() {
		if printRawTime {
			fmt.Printf("%d, %d, %s\n", e.GetEventId(), e.GetTimestamp(), util.HistoryEventToString(e))
		} else {
			fmt.Printf("%d, %s, %s\n", e.GetEventId(), convertTime(e.GetTimestamp()), util.HistoryEventToString(e))
		}
	}
}

// StartWorkflow starts a new workflow execution
func StartWorkflow(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	tasklist := getRequiredOption(c, FlagTaskList)
	workflowType := getRequiredOption(c, FlagWorkflowType)
	et := c.Int(FlagExecutionTimeout)
	if et == 0 {
		ExitIfError(errors.New(FlagExecutionTimeout + " is required"))
	}
	dt := c.Int(FlagDecisionTimeout)
	if et == 0 {
		ExitIfError(errors.New(FlagDecisionTimeout + " is required"))
	}
	wid := c.String(FlagWorkflowID)
	if len(wid) == 0 {
		wid = uuid.New()
	}

	input := c.String(FlagInput)

	workflowOptions := cadence.StartWorkflowOptions{
		ID:                              wid,
		TaskList:                        tasklist,
		ExecutionStartToCloseTimeout:    time.Duration(et) * time.Second,
		DecisionTaskStartToCloseTimeout: time.Duration(dt) * time.Second,
	}

	var we *cadence.WorkflowExecution
	var err error
	if len(input) > 0 {
		// assume workflow takes one input of string type
		we, err = wfClient.StartWorkflow(workflowOptions, workflowType, input)
	} else {
		// assume workflow takes no input
		we, err = wfClient.StartWorkflow(workflowOptions, workflowType)
	}
	if err != nil {
		fmt.Printf("Failed to create workflow with error: %+v\n", err)
	} else {
		fmt.Printf("Started Workflow Id: %s, run Id: %s\n", we.ID, we.RunID)
	}
}

// TerminateWorkflow terminates a workflow execution
func TerminateWorkflow(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	wid := getRequiredOption(c, FlagWorkflowID)
	rid := c.String(FlagRunID)
	reason := c.String(FlagReason)

	err := wfClient.TerminateWorkflow(wid, rid, reason, nil)

	if err != nil {
		fmt.Printf("Terminate workflow failed: %v\n", err)
	} else {
		fmt.Println("Terminate workflow succeed.")
	}
}

// CancelWorkflow cancels a workflow execution
func CancelWorkflow(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	wid := getRequiredOption(c, FlagWorkflowID)
	rid := c.String(FlagRunID)

	err := wfClient.CancelWorkflow(wid, rid)

	if err != nil {
		fmt.Printf("Cancel workflow failed: %v\n", err)
	} else {
		fmt.Println("Cancel workflow succeed.")
	}
}

// SignalWorkflow signals a workflow execution
func SignalWorkflow(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	wid := getRequiredOption(c, FlagWorkflowID)
	rid := c.String(FlagRunID)
	name := getRequiredOption(c, FlagName)
	input := c.String(FlagInput)

	var err error
	if len(input) > 0 {
		err = wfClient.SignalWorkflow(wid, rid, name, input)
	} else {
		err = wfClient.SignalWorkflow(wid, rid, name, nil)
	}

	if err != nil {
		fmt.Printf("Signal workflow failed: %v\n", err)
	} else {
		fmt.Println("Signal workflow succeed.")
	}
}

// QueryWorkflow list workflow executions based on query filters
func QueryWorkflow(c *cli.Context) {
	wfClient := getWorkflowClient(c)

	queryOpen := c.Bool(FlagOpen)
	pageSize := c.Int(FlagPageSize)
	earliestTime := parseTime(c.String(FlagEarliestTime), 0)
	latestTime := parseTime(c.String(FlagLatestTime), time.Now().UnixNano())
	workflowID := c.String(FlagWorkflowID)
	workflowType := c.String(FlagWorkflowType)
	printRawTime := c.Bool(FlagPrintRawTime)

	if len(workflowID) > 0 && len(workflowType) > 0 {
		ExitIfError(errors.New("you can filter on workflow_id or workflow_type, but not on both"))
	}

	reader := bufio.NewReader(os.Stdin)
	var result []*s.WorkflowExecutionInfo
	var nextPageToken []byte
	for {
		if queryOpen {
			result, nextPageToken = queryOpenWorkflow(wfClient, pageSize, earliestTime, latestTime, workflowID, workflowType, nextPageToken)
		} else {
			result, nextPageToken = queryClosedWorkflow(wfClient, pageSize, earliestTime, latestTime, workflowID, workflowType, nextPageToken)
		}

		for _, e := range result {
			fmt.Printf("%s, -w %s -r %s", e.GetType().GetName(), e.GetExecution().GetWorkflowId(), e.GetExecution().GetRunId())
			if printRawTime {
				fmt.Printf(" [%d, %d]\n", e.GetStartTime(), e.GetCloseTime())
			} else {
				fmt.Printf(" [%s, %s]\n", convertTime(e.GetStartTime()), convertTime(e.GetCloseTime()))
			}
		}

		if len(result) < pageSize {
			break
		}

		fmt.Println("Press C then Enter to show more result, press any other key then Enter to quit: ")
		input, _ := reader.ReadString('\n')
		c := []byte(input)[0]
		if c == 'C' || c == 'c' {
			continue
		} else {
			break
		}
	}
}

func queryOpenWorkflow(client cadence.Client, pageSize int, earliestTime, latestTime int64, workflowID, workflowType string, nextPageToken []byte) ([]*s.WorkflowExecutionInfo, []byte) {
	request := &s.ListOpenWorkflowExecutionsRequest{
		MaximumPageSize: common.Int32Ptr(int32(pageSize)),
		NextPageToken:   nextPageToken,
		StartTimeFilter: &s.StartTimeFilter{
			EarliestTime: common.Int64Ptr(earliestTime),
			LatestTime:   common.Int64Ptr(latestTime),
		},
	}
	if len(workflowID) > 0 {
		request.ExecutionFilter = &s.WorkflowExecutionFilter{WorkflowId: common.StringPtr(workflowID)}
	}
	if len(workflowType) > 0 {
		request.TypeFilter = &s.WorkflowTypeFilter{Name: common.StringPtr(workflowType)}
	}

	response, err := client.ListOpenWorkflow(request)
	if err != nil {
		ExitIfError(err)
	}
	return response.GetExecutions(), response.GetNextPageToken()
}

func queryClosedWorkflow(client cadence.Client, pageSize int, earliestTime, latestTime int64, workflowID, workflowType string, nextPageToken []byte) ([]*s.WorkflowExecutionInfo, []byte) {
	request := &s.ListClosedWorkflowExecutionsRequest{
		MaximumPageSize: common.Int32Ptr(int32(pageSize)),
		NextPageToken:   nextPageToken,
		StartTimeFilter: &s.StartTimeFilter{
			EarliestTime: common.Int64Ptr(earliestTime),
			LatestTime:   common.Int64Ptr(latestTime),
		},
	}
	if len(workflowID) > 0 {
		request.ExecutionFilter = &s.WorkflowExecutionFilter{WorkflowId: common.StringPtr(workflowID)}
	}
	if len(workflowType) > 0 {
		request.TypeFilter = &s.WorkflowTypeFilter{Name: common.StringPtr(workflowType)}
	}

	response, err := client.ListClosedWorkflow(request)
	if err != nil {
		ExitIfError(err)
	}
	return response.GetExecutions(), response.GetNextPageToken()
}

func getDomainClient(c *cli.Context) cadence.DomainClient {
	address := c.GlobalString(FlagAddress)

	builder := getBuilder(address)
	domainClient, err := builder.BuildCadenceDomainClient()
	if err != nil {
		ExitIfError(err)
	}
	return domainClient
}

func getWorkflowClient(c *cli.Context) cadence.Client {
	address := c.GlobalString(FlagAddress)
	domain := getRequiredGlobalOption(c, FlagDomain)

	builder := getBuilder(address).SetDomain(domain)
	wfClient, err := builder.BuildCadenceClient()
	if err != nil {
		ExitIfError(err)
	}

	return wfClient
}

func getRequiredOption(c *cli.Context, optionName string) string {
	value := c.String(optionName)
	if len(value) == 0 {
		ExitIfError(fmt.Errorf("%s is required", optionName))
	}
	return value
}

func getRequiredGlobalOption(c *cli.Context, optionName string) string {
	value := c.GlobalString(optionName)
	if len(value) == 0 {
		ExitIfError(fmt.Errorf("%s is required", optionName))
	}
	return value
}

func getBuilder(address string) *factory.WorkflowClientBuilder {
	builder := factory.NewBuilder()
	if len(address) == 0 {
		address = localHostPort
	}
	builder = builder.SetHostPort(address)
	return builder
}

func convertTime(unixNano int64) string {
	t2 := time.Unix(0, unixNano)
	return t2.Format(time.RFC3339)
}

func parseTime(timeStr string, defaultValue int64) int64 {
	if len(timeStr) == 0 {
		return defaultValue
	}

	// try to parse
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return parsedTime.UnixNano()
	}

	// treat as raw time
	resultValue, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		ExitIfError(fmt.Errorf("cannot parse time '%s', use RFC3339 format '%s' or raw UnixNano directly", timeStr, time.RFC3339))
	}

	return resultValue
}
