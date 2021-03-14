package shigoto

// Jobs instance
type Jobs struct {
	JobName        string
	JobFunc        interface{}
	JobParams      []interface{}
	PriorityStatus string   // If you set high, then the job will be prioritize
	Cron           []string // Set run a jobs with periodic by second, minute and hour
}

// Priority to set job priority (normal, medium, high)
func (j *Jobs) Priority(priority string) *Jobs {
	switch priority {
	case "normal", "medium", "high":
		j.PriorityStatus = priority
	default:
		j.PriorityStatus = "normal"
	}

	return j
}
