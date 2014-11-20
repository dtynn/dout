package out

type statusNum int

const (
	StatusSuccess statusNum = iota
	StatusInvalidToAddress
	StatusDnsErr
	StatusSmtpErr
)

type Status struct {
	Num    statusNum
	Email  string
	Detail string
}

func NewStatus(num statusNum, email, detail string) *Status {
	return &Status{
		Num:    num,
		Email:  email,
		Detail: detail,
	}
}
