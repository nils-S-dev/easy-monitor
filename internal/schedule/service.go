package schedule

import (
	"easy-monitor/internal/config"
	"easy-monitor/internal/monitor"
	"easy-monitor/internal/notification"

	"github.com/robfig/cron/v3"
)

func Init() {

	c := cron.New()
	conf := config.GetConfig()
	monitors := conf.Monitors
	globalCron := conf.Cron

	var monitorsWithoutOwnCron []config.Monitor
	var monitorsWithOwnCron []config.Monitor

	// sort monitors in global and local cron
	for _, monitor := range monitors {
		if monitor.Cron != "" {
			monitorsWithOwnCron = append(monitorsWithOwnCron, monitor)
		} else {
			monitorsWithoutOwnCron = append(monitorsWithoutOwnCron, monitor)
		}
	}

	// add one global cron if defined
	if globalCron != "" {
		c.AddFunc(globalCron, func() {
			notifyForFailedMonitors(monitorsWithoutOwnCron)
		})
	}

	// group by custom cron
	monitorsWithOwnCronGrouped := make(map[string][]config.Monitor)
	for _, m := range monitorsWithOwnCron {
		monitorsWithOwnCronGrouped[m.Cron] = append(monitorsWithOwnCronGrouped[m.Cron], m)
	}
	for cron, mons := range monitorsWithOwnCronGrouped {
		c.AddFunc(cron, func() {
			notifyForFailedMonitors(mons)
		})
	}

	c.Start()
}

func notifyForFailedMonitors(monitors []config.Monitor) {
	results := monitor.GetMonitorResults(monitors)
	var failedMonitorResults []monitor.MonitorResult
	for _, mr := range results {
		if mr.Status == monitor.StatusFail {
			failedMonitorResults = append(failedMonitorResults, mr)
		}
	}
	notification.Notify(failedMonitorResults)
}
