package storage

type status string

const (
	pending    status = "PENDING"
	inProgress status = "IN_PROGRESS"
	completed  status = "COMPLETED"
	failed     status = "FAILED"
)
