package response

type SysPatrolConfig struct {
	ID              uint  `json:"ID" gorm:"column:id"`
	PatrolType      int   `json:"patrolType" gorm:"column:patrol_type"`
	IntervalHours   int64 `json:"intervalHours" gorm:"column:interval_hours"`
	IntervalMinutes int64 `json:"intervalMinutes" gorm:"column:interval_minutes"`
	IntervalSeconds int64 `json:"intervalSeconds" gorm:"column:interval_seconds"`
}
