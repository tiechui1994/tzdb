//+build !go1.14

package tzdb

import (
	"time"
	"unsafe"
)

type Location struct {
	name string
	zone []zone
	tx   []zoneTrans

	cacheStart int64
	cacheEnd   int64
	cacheZone  *zone
}

// A zone represents a single time zone such as CEST or CET.
type zone struct {
	name   string // abbreviated name, "CET"
	offset int    // seconds east of UTC
	isDST  bool   // is this zone Daylight Savings Time?
}

// A zoneTrans represents a single time zone transition.
type zoneTrans struct {
	when         int64 // transition time, in seconds since 1970 GMT
	index        uint8 // the index of the zone that goes into effect at that time
	isstd, isutc bool  // ignored - no idea what these mean
}

func (l *Location) look(sec int64) (isDST bool) {
	if len(l.zone) == 0 {
		isDST = false
		return
	}

	if zone := l.cacheZone; zone != nil && l.cacheStart <= sec && sec < l.cacheEnd {
		isDST = zone.isDST
		return
	}

	if len(l.tx) == 0 || sec < l.tx[0].when {
		zone := &l.zone[l.lookupFirstZone()]
		isDST = zone.isDST
		return
	}

	tx := l.tx
	lo := 0
	hi := len(tx)
	for hi-lo > 1 {
		m := lo + (hi-lo)/2
		lim := tx[m].when
		if sec < lim {
			hi = m
		} else {
			lo = m
		}
	}
	zone := &l.zone[tx[lo].index]
	isDST = zone.isDST

	return isDST
}

func (l *Location) lookupFirstZone() int {
	// Case 1.
	used := false
	for _, tx := range l.tx {
		if tx.index == 0 {
			used = true
		}
	}
	if !used {
		return 0
	}

	// Case 2.
	if len(l.tx) > 0 && l.zone[l.tx[0].index].isDST {
		for zi := int(l.tx[0].index) - 1; zi >= 0; zi-- {
			if !l.zone[zi].isDST {
				return zi
			}
		}
	}

	// Case 3.
	for zi := range l.zone {
		if !l.zone[zi].isDST {
			return zi
		}
	}

	// Case 4.
	return 0
}

func IsDST(t time.Time) bool {
	l := t.Location()
	x := (*Location)(unsafe.Pointer(l))
	return x.look(t.Unix())
}
