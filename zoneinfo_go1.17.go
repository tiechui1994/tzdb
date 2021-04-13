//+build go1.17

package tzdb

import "time"

func IsDST(t time.Time) bool {
	return t.IsDST()
}
