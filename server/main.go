package main

import (
	"gitlab.com/w1572/backend/plugin"
)

func main() {
	plugin.ClientMain(&Plugin{
		commandHandlers: map[string]commandHandlerFunc{
			setLogsLimitKey:       setLogsLimit,
			setLogsStartTimeKey:   setLogsStartTime,
			getLogsKey:            getLogs,
			setReportFrequencyKey: setReportFrequency,
			subscribeKey:          subscribe,
			unsubscribeKey:        unsubscribe,
		},
	})
}
