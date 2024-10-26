package storage

const (
	_articleTable = "article"
)

type status string

const (
	_pending    status = "PENDING"
	_inProgress status = "IN_PROGRESS"
	_completed  status = "COMPLETED"
	_failed     status = "FAILED"
)
