package gronx

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CronDateFormat is Y-m-d H:i (seconds are not significant)
const CronDateFormat = "2006-01-02 15:04"

// FullDateFormat is Y-m-d H:i:s (with seconds)
const FullDateFormat = "2006-01-02 15:04:05"

// NextTick gives next run time from now
func NextTick(expr string, inclRefTime bool) (time.Time, error) {
	return NextTickAfter(expr, time.Now(), inclRefTime)
}

// NextTickAfter gives next run time from the provided time.Time
func NextTickAfter(expr string, start time.Time, inclRefTime bool) (time.Time, error) {
	gron, next := New(), start.Truncate(time.Minute)
	due, err := gron.IsDue(expr, start)
	if err != nil || (due && inclRefTime) {
		return start, err
	}

	segments, _ := Segments(expr)
	if len(segments) > 5 && isUnreachableYear(segments[5], next, inclRefTime, false) {
		return next, fmt.Errorf("unreachable year segment: %s", segments[5])
	}

	next, err = loop(gron, segments, next, inclRefTime, false)
	// Ignore superfluous err
	if err != nil && gron.isDue(expr, next) {
		err = nil
	}
	return next, err
}

func loop(gron Gronx, segments []string, start time.Time, incl bool, reverse bool) (next time.Time, err error) {
	iter, next, bumped := 500, start, false
over:
	for iter > 0 {
		iter--
		for pos, seg := range segments {
			if seg == "*" || seg == "?" {
				continue
			}
			if next, bumped, err = bumpUntilDue(gron.C, seg, pos, next, reverse); bumped {
				goto over
			}
		}
		if !incl && next.Format(FullDateFormat) == start.Format(FullDateFormat) {
			delta := time.Minute
			if reverse {
				delta = -time.Minute
			}
			next, _, err = bumpUntilDue(gron.C, segments[0], 0, next.Add(delta), reverse)
			continue
		}
		return
	}
	return start, errors.New("tried so hard")
}

var dashRe = regexp.MustCompile(`/.*$`)

func isUnreachableYear(year string, ref time.Time, incl bool, reverse bool) bool {
	if year == "*" || year == "?" {
		return false
	}

	edge, inc := ref.Year(), 1
	if !incl {
		if reverse {
			inc = -1
		}
		edge += inc
	}
	for _, offset := range strings.Split(year, ",") {
		if strings.Index(offset, "*/") == 0 || strings.Index(offset, "0/") == 0 {
			return false
		}
		for _, part := range strings.Split(dashRe.ReplaceAllString(offset, ""), "-") {
			val, err := strconv.Atoi(part)
			if err != nil || (!reverse && val >= edge) || (reverse && val < edge) {
				return false
			}
		}
	}
	return true
}

var limit = map[int]int{0: 60, 1: 24, 2: 31, 3: 12, 4: 366, 5: 100}

func bumpUntilDue(c Checker, segment string, pos int, ref time.Time, reverse bool) (time.Time, bool, error) {
	// <minute> <hour> <day> <month> <weekday> <year>
	iter := limit[pos]
	for iter > 0 {
		c.SetRef(ref)
		if ok, _ := c.CheckDue(segment, pos); ok {
			return ref, iter != limit[pos], nil
		}
		ref = bump(ref, pos, reverse)
		iter--
	}
	return ref, false, errors.New("tried so hard")
}

func bump(ref time.Time, pos int, reverse bool) time.Time {
	factor := 1
	if reverse {
		factor = -1
	}

	switch pos {
	case 0:
		ref = ref.Add(time.Duration(factor) * time.Minute)
	case 1:
		ref = ref.Add(time.Duration(factor) * time.Hour)
	case 2, 4:
		ref = ref.AddDate(0, 0, factor)
	case 3:
		ref = ref.AddDate(0, factor, 0)
	case 5:
		ref = ref.AddDate(factor, 0, 0)
	}
	return ref
}
