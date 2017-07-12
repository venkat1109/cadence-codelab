package workflow

import (
	//"github.com/venkat1109/cadence-codelab/cron/activity"
	"errors"
	"go.uber.org/cadence"
	"go.uber.org/zap"
	"time"
)

type (
	CronSchedule struct {
		Count      int           // total number of jobs to schedule
		Frequency  time.Duration // frequency at which to schedule jobs
		Hostgroups []string      // schedule a job for each one of these hostgroup
	}
)

const maxJobsPerLoop = 1

func init() {
	cadence.RegisterWorkflow(Cron)
}

// Cron implements the Cron workflow
func Cron(ctx cadence.Context, schedule *CronSchedule) error {

	cadence.GetLogger(ctx).Info("Cron started", zap.Int("Count", schedule.Count),
		zap.Duration("frequency", schedule.Frequency), zap.Strings("groups", schedule.Hostgroups))

	activityCtx := cadence.WithActivityOptions(ctx, cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
		HeartbeatTimeout:       time.Minute,
	})

	return runScheduler(ctx, activityCtx, schedule)
}

// runScheduler runs the cron scheduler
func runScheduler(ctx cadence.Context, activityCtx cadence.Context, schedule *CronSchedule) error {
	return errors.New("not implemented")
}
