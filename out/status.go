package out

type statusNum int

const (
	StatusSuccess statusNum = iota
	StatusInvalidToAddress
	StatusDnsErr
	StatusSmtpErr
)

type Status struct {
	ErrNo  statusNum `json:"error"`
	Email  string    `json:"email"`
	Detail string    `json:"detail"`
}

func NewStatus(num statusNum, email, detail string) *Status {
	return &Status{
		ErrNo:  num,
		Email:  email,
		Detail: detail,
	}
}
