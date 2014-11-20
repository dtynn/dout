package out

import (
	"fmt"
)

type InvalidEmail string

func (this InvalidEmail) Error() string {
	return fmt.Sprintf("<%s> is an invalid email address", this)
}

var (
	noRecipients = fmt.Errorf("no recipients")
	someFailed   = fmt.Errorf("some failed")
)

func SomeFailed(err error) bool {
	return err == someFailed
}
