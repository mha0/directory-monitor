package domain

type Status int

const (
	OPERATIONAL Status = iota + 1
	WARNING
	INITIALIZED
)

func (status Status) String() string {
	names := [...]string{
		"OPERATIONAL",
		"WARNING",
		"INITIALIZED"}

	if status < OPERATIONAL || status > INITIALIZED {
		return "Unknown"
	}

	return names[status-1]
}

