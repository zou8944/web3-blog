package types

import (
	"fmt"
	"time"
)

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	unixTime := fmt.Sprintf("%d", time.Time(t).UnixMilli())
	return []byte(unixTime), nil
}
