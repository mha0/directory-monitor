package domain

type Status int

const (
	OPERATIONAL Status = iota
	INITIALIZED
	WARNING
)

func (status Status) String() string {
	names := [...]string{
		"OPERATIONAL",
		"INITIALIZED",
		"WARNING"}

	if status < OPERATIONAL || status > WARNING {
		return "Unknown"
	}

	return names[status]
}

