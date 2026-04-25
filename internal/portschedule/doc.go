// Package portschedule implements time-window based scheduling for port scans.
//
// A Schedule is configured with one or more Windows, each defining an inclusive
// start hour and exclusive end hour (24-hour clock). Overnight windows are
// supported by setting StartHour > EndHour (e.g. 22–06).
//
// When no windows are configured the schedule is always considered active,
// preserving backward-compatible behaviour for users who do not require
// time-restricted scanning.
//
// Example:
//
//	win := []portschedule.Window{{StartHour: 8, EndHour: 18}}
//	sched, err := portschedule.New(win, time.Local)
//	if err != nil { ... }
//	if sched.Active(time.Now()) {
//	    // run scan
//	}
package portschedule
