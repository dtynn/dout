package status

type statusNum int

const (
	StatusSuccess statusNum = iota
	StatusDNS
	StatusSMTP
)

type Status struct {
	num    statusNum
	detail string
}
