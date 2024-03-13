package define

const (
	QueueStatusRun = iota + 1
	QueueStatusPause
	QueueStatusFinish
	QueueStatusStop
	QueueStatusExceed
	QueueStatusAnalyzing
	QueueStatusAnalyzingFail
)

const (
	QueueSectorStatusWait = iota + 1
	QueueSectorStatusCreate
	QueueSectorStatusFinish
	QueueSectorStatusFail
	QueueSectorStatusFitFail
	QueueSectorStatusExceed
)
