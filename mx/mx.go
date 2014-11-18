package mx

import (
	"fmt"
	"net"
)

var errMxNotFound = fmt.Errorf("mx records not found")

func ChoseMx(name string) (*net.MX, error) {
	mxs, err := net.LookupMX(name)
	if err != nil {
		return nil, err
	}

	size := len(mxs)
	if size == 0 {
		return nil, errMxNotFound
	}

	first := mxs[0]

	if size == 1 {
		return first, nil
	}

	mxsMap := map[uint16][]*net.MX{}
	minimum := first.Pref
	mxsMap[minimum] = []*net.MX{first}

	for _, mx := range mxs[1:] {
		if mx.Pref < minimum {
			minimum = mx.Pref
		}

		if slice, has := mxsMap[mx.Pref]; has {
			mxsMap[mx.Pref] = append(slice, mx)
		} else {
			mxsMap[mx.Pref] = []*net.MX{mx}
		}
	}

	// todo: random choise or ping test
	return mxsMap[minimum][0], nil
}
