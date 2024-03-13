package response

import "time"

type C2TaskInfo struct {
	MinerId string
	Number  int
	RunTime time.Time
}
