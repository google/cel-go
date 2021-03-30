// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"math"
	"time"
)

// addInt64Checked performs addition with overflow detection of two int64, returning the result of
// the addition if no overflow occurred as the first return value and a bool indicating whether no
// overflow occurred as the second return value.
func addInt64Checked(x, y int64) (int64, bool) {
	if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
		return 0, false
	}
	return x + y, true
}

// subtractInt64Checked performs subtraction with overflow detection of two int64, returning the
// result of the subtraction if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func subtractInt64Checked(x, y int64) (int64, bool) {
	if (y < 0 && x > math.MaxInt64+y) || (y > 0 && x < math.MinInt64+y) {
		return 0, false
	}
	return x - y, true
}

// negateInt64Checked performs negation with overflow detection of an int64, returning the
// result of the negation if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func negateInt64Checked(x int64) (int64, bool) {
	// In twos complement, negating MinInt64 would result in a valid of MaxInt64+1.
	if x == math.MinInt64 {
		return 0, false
	}
	return -x, true
}

// multiplyInt64Checked performs multiplication with overflow detection of two int64, returning the
// result of the multiplication if no overflow occurred as the first return value and a bool
// indicating whether no overflow occurred as the second return value.
func multiplyInt64Checked(x, y int64) (int64, bool) {
	// Detecting multiplication overflow is more complicated than the others. The first two detect
	// attempting to negate MinInt64, which would result in MaxInt64+1. The other four detect normal
	// overflow conditions.
	if (x == -1 && y == math.MinInt64) || (y == -1 && x == math.MinInt64) ||
		// x is positive, y is positive
		(x > 0 && y > 0 && x > math.MaxInt64/y) ||
		// x is positive, y is negative
		(x > 0 && y < 0 && y < math.MinInt64/x) ||
		// x is negative, y is positive
		(x < 0 && y > 0 && x < math.MinInt64/y) ||
		// x is negative, y is negative
		(x < 0 && y < 0 && y < math.MaxInt64/x) {
		return 0, false
	}
	return x * y, true
}

// divideInt64Checked performs division with overflow detection of two int64, returning the
// result of the division if no overflow occurred as the first return value and a bool
// indicating whether no overflow occurred as the second return value.
func divideInt64Checked(x, y int64) (int64, bool) {
	// In twos complement, negating MinInt64 would result in a valid of MaxInt64+1.
	if x == math.MinInt64 && y == -1 {
		return 0, false
	}
	return x / y, true
}

// moduloInt64Checked performs modulo with overflow detection of two int64, returning the
// result of the modulo if no overflow occurred as the first return value and a bool
// indicating whether no overflow occurred as the second return value.
func moduloInt64Checked(x, y int64) (int64, bool) {
	// In twos complement, negating MinInt64 would result in a valid of MaxInt64+1.
	if x == math.MinInt64 && y == -1 {
		return 0, false
	}
	return x % y, true
}

// addUint64Checked performs addition with overflow detection of two uint64, returning the result of
// the addition if no overflow occurred as the first return value and a bool indicating whether no
// overflow occurred as the second return value.
func addUint64Checked(x, y uint64) (uint64, bool) {
	if y > 0 && x > math.MaxUint64-y {
		return 0, false
	}
	return x + y, true
}

// subtractUint64Checked performs subtraction with overflow detection of two uint64, returning the
// result of the subtraction if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func subtractUint64Checked(x, y uint64) (uint64, bool) {
	if y > x {
		return 0, false
	}
	return x - y, true
}

// multiplyUint64Checked performs multiplication with overflow detection of two uint64, returning
// the result of the multiplication if no overflow occurred as the first return value and a bool
// indicating whether no overflow occurred as the second return value.
func multiplyUint64Checked(x, y uint64) (uint64, bool) {
	if y != 0 && x > math.MaxUint64/y {
		return 0, false
	}
	return x * y, true
}

// addDurationChecked performs addition with overflow detection of two time.Duration, returning the
// result of the addition if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func addDurationChecked(x, y time.Duration) (time.Duration, bool) {
	if val, ok := addInt64Checked(int64(x), int64(y)); ok {
		return time.Duration(val), true
	}
	return time.Duration(0), false
}

// subtractDurationChecked performs subtraction with overflow detection of two time.Duration,
// returning the result of the subtraction if no overflow occurred as the first return value and a
// bool indicating whether no overflow occurred as the second return value.
func subtractDurationChecked(x, y time.Duration) (time.Duration, bool) {
	if val, ok := subtractInt64Checked(int64(x), int64(y)); ok {
		return time.Duration(val), true
	}
	return time.Duration(0), false
}

// negateDurationChecked performs negation with overflow detection of a time.Duration, returning the
// result of the negation if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func negateDurationChecked(x time.Duration) (time.Duration, bool) {
	if val, ok := negateInt64Checked(int64(x)); ok {
		return time.Duration(val), true
	}
	return time.Duration(0), false
}

// addDurationChecked performs addition with overflow detection of a time.Time and time.Duration,
// returning the result of the addition if no overflow occurred as the first return value and a bool
// indicating whether no overflow occurred as the second return value.
func addTimeDurationChecked(x time.Time, y time.Duration) (time.Time, bool) {
	// This is tricky. A time is represented as (int64, int32) where the first is seconds and second
	// is nanoseconds. A duration is int64 representing nanoseconds. We cannot normalize time to int64
	// as it could potentially overflow. The only way to proceed is to break time and duration into
	// second and nanosecond components.

	// First we break time into its components by truncating and subtracting.
	sec1 := x.Truncate(time.Second).Unix()                // Truncate to seconds.
	nsec1 := x.Sub(x.Truncate(time.Second)).Nanoseconds() // Get nanoseconds by truncating and subtracting.

	// Second we break duration into its components by dividing and modulo.
	sec2 := int64(y) / int64(time.Second)  // Truncate to seconds.
	nsec2 := int64(y) % int64(time.Second) // Get remainder.

	// Add seconds first, detecting any overflow.
	sec, ok := addInt64Checked(sec1, sec2)
	if !ok {
		return time.Time{}, false
	}

	// Nanoseconds cannot overflow as time.Time normalizes them to [0, 999999999].
	nsec := nsec1 + nsec2

	// We need to normalize nanoseconds to be positive and carry extra nanoseconds to seconds.
	// Adapted from time.Unix(int64, int64).
	if nsec < 0 || nsec >= int64(time.Second) {
		// Add seconds.
		sec, ok = addInt64Checked(sec, nsec/int64(time.Second))
		if !ok {
			return time.Time{}, false
		}
		nsec -= (nsec / int64(time.Second)) * int64(time.Second)
		if nsec < 0 {
			// Subtract an extra second
			sec, ok = addInt64Checked(sec, -1)
			if !ok {
				return time.Time{}, false
			}
			nsec += int64(time.Second)
		}
	}

	// Check if the the number of seconds from Unix epoch is within our acceptable range.
	if sec < minUnixTime || sec > maxUnixTime {
		return time.Time{}, false
	}

	// Return resulting time and propagate time zone.
	return time.Unix(sec, nsec).In(x.Location()), true
}

// subtractTimeChecked performs subtraction with overflow detection of two time.Time, returning the
// result of the subtraction if no overflow occurred as the first return value and a bool indicating
// whether no overflow occurred as the second return value.
func subtractTimeChecked(x, y time.Time) (time.Duration, bool) {
	// Similar to addTimeDurationOverflow() above.

	// First we break time into its components by truncating and subtracting.
	sec1 := x.Truncate(time.Second).Unix()                // Truncate to seconds.
	nsec1 := x.Sub(x.Truncate(time.Second)).Nanoseconds() // Get nanoseconds by truncating and subtracting.

	// Second we break duration into its components by truncating and subtracting.
	sec2 := y.Truncate(time.Second).Unix()                // Truncate to seconds.
	nsec2 := y.Sub(y.Truncate(time.Second)).Nanoseconds() // Get nanoseconds by truncating and subtracting.

	// Subtract seconds first, detecting any overflow.
	sec, ok := subtractInt64Checked(sec1, sec2)
	if !ok {
		return time.Duration(0), false
	}

	// Nanoseconds cannot overflow as time.Time normalizes them to [0, 999999999].
	nsec := nsec1 - nsec2

	// We need to normalize nanoseconds to be positive and carry extra nanoseconds to seconds.
	// Adapted from time.Unix(int64, int64).
	if nsec < 0 || nsec >= int64(time.Second) {
		// Add an extra second.
		sec, ok = addInt64Checked(sec, nsec/int64(time.Second))
		if !ok {
			return time.Duration(0), false
		}
		nsec -= (nsec / int64(time.Second)) * int64(time.Second)
		if nsec < 0 {
			// Subtract an extra second
			sec, ok = addInt64Checked(sec, -1)
			if !ok {
				return time.Duration(0), false
			}
			nsec += int64(time.Second)
		}
	}

	// Scale seconds to nanoseconds detecting overflow.
	tsec, ok := multiplyInt64Checked(sec, int64(time.Second))
	if !ok {
		return time.Duration(0), false
	}

	// Lastly we need to add the two nanoseconds together.
	val, ok := addInt64Checked(tsec, nsec)
	if !ok {
		return time.Duration(0), false
	}

	return time.Duration(val), true
}

// subtractTimeDurationChecked performs subtraction with overflow detection of a time.Time and
// time.Duration, returning the result of the subtraction if no overflow occurred as the first
// return value and a bool indicating whether no overflow occurred as the second return value.
func subtractTimeDurationChecked(x time.Time, y time.Duration) (time.Time, bool) {
	// The easiest way to implement this is to negate y and add them.
	// x - y = x + -y
	val, ok := negateDurationChecked(y)
	if !ok {
		return time.Time{}, false
	}
	return addTimeDurationChecked(x, val)
}
