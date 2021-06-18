package todoist

type Priority int

const (
	PriorityNormal Priority = 1
	PriorityHigh   Priority = 2
	PriorityHigher Priority = 3
	PriorityUrgent Priority = 4
)
