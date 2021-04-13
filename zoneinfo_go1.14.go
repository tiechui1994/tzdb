//+build go1.14,!go1.17

package tzdb

import (
	"time"
	"unsafe"
)

func init() {
	time.Now().Zone()
}

type Location struct {
	name string
	zone []zone
	tx   []zoneTrans

	extend string

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
	end := int64(1<<63 - 1)
	lo := 0
	hi := len(tx)
	for hi-lo > 1 {
		m := lo + (hi-lo)/2
		lim := tx[m].when
		if sec < lim {
			end = lim
			hi = m
		} else {
			lo = m
		}
	}
	zone := &l.zone[tx[lo].index]
	isDST = zone.isDST

	if lo == len(tx)-1 && l.extend != "" {
		if isDST, ok := tzsetIsDST(l.extend, end, sec); ok {
			return isDST
		}
	}

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

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
)

const (
	unixToInternal     int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	internalToAbsolute int64 = (292277022399 + 1) * 365.2425 * secondsPerDay
)

type ruleKind int

type rule struct {
	kind ruleKind
	day  int
	week int
	mon  int
	time int // transition time
}

//go:linkname tzsetOffset time.tzsetOffset
func tzsetOffset(s string) (offset int, rest string, ok bool)

//go:linkname tzruleTime time.tzruleTime
func tzruleTime(year int, r rule, off int) int

//go:linkname tzsetRule time.tzsetRule
func tzsetRule(s string) (rule, string, bool)

//go:linkname absDate time.absDate
func absDate(abs uint64, full bool) (year int, month time.Month, day int, yday int)

func tzsetIsDST(s string, initEnd, sec int64) (isDST, ok bool) {
	var (
		stdOffset, dstOffset int
	)

	_, s, ok = tzsetName(s)
	if ok {
		stdOffset, s, ok = tzsetOffset(s)
	}
	if !ok {
		return false, false
	}

	// The numbers in the tzset string are added to local time to get UTC,
	// but our offsets are added to UTC to get local time,
	// so we negate the number we see here.
	stdOffset = -stdOffset

	if len(s) == 0 || s[0] == ',' {
		// No daylight savings time.
		return false, true
	}

	_, s, ok = tzsetName(s)
	if ok {
		if len(s) == 0 || s[0] == ',' {
			dstOffset = stdOffset + secondsPerHour
		} else {
			dstOffset, s, ok = tzsetOffset(s)
			dstOffset = -dstOffset // as with stdOffset, above
		}
	}
	if !ok {
		return false, false
	}

	if len(s) == 0 {
		// Default DST rules per tzcode.
		s = ",M3.2.0,M11.1.0"
	}
	// The TZ definition does not mention ';' here but tzcode accepts it.
	if s[0] != ',' && s[0] != ';' {
		return false, false
	}
	s = s[1:]

	var startRule, endRule rule
	startRule, s, ok = tzsetRule(s)
	if !ok || len(s) == 0 || s[0] != ',' {
		return false, false
	}
	s = s[1:]
	endRule, s, ok = tzsetRule(s)
	if !ok || len(s) > 0 {
		return false, false
	}

	year, _, _, yday := absDate(uint64(sec+unixToInternal+internalToAbsolute), false)

	ysec := int64(yday*secondsPerDay) + sec%secondsPerDay

	startSec := int64(tzruleTime(year, startRule, stdOffset))
	endSec := int64(tzruleTime(year, endRule, dstOffset))
	dstIsDST, stdIsDST := true, false
	if endSec < startSec {
		stdIsDST, dstIsDST = dstIsDST, stdIsDST
	}

	if ysec < startSec {
		return stdIsDST, true
	} else if ysec >= endSec {
		return stdIsDST, true
	} else {
		return dstIsDST, true
	}
}

func tzsetName(s string) (string, string, bool) {
	if len(s) == 0 {
		return "", "", false
	}
	if s[0] != '<' {
		for i, r := range s {
			switch r {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '-', '+':
				if i < 3 {
					return "", "", false
				}
				return s[:i], s[i:], true
			}
		}
		if len(s) < 3 {
			return "", "", false
		}
		return s, "", true
	} else {
		for i, r := range s {
			if r == '>' {
				return s[1:i], s[i+1:], true
			}
		}
		return "", "", false
	}
}

func IsDST(t time.Time) bool {
	l := t.Location()
	x := (*Location)(unsafe.Pointer(l))
	return x.look(t.Unix())
}
