// Copyright 2014 Maxim Kouprianov. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

/*
Package handysort implements an alphanumeric string comparison function
in order to sort alphanumeric strings correctly.

Default sort (incorrect):
	abc1
	abc10
	abc12
	abc2

Handysort:
	abc1
	abc2
	abc10
	abc12

Please note, that handysort is about 5x-8x times slower
than a simple sort, so use it wisely.
*/
package handysort

import (
	"unicode/utf8"
)

// Strings implements the sort interface, sorts an array
// of the alphanumeric strings in decreasing order.
type Strings []string

func (a Strings) Len() int           { return len(a) }
func (a Strings) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Strings) Less(i, j int) bool { return StringLess(a[i], a[j]) }

// StringLess compares two alphanumeric strings correctly.
func StringLess(s1, s2 string) (less bool) {
	n1, n2 := make([]rune, 0, 20), make([]rune, 0, 20)

	for i, j := 0, 0; i < len(s1) || j < len(s2); {
		var r1, r2 rune
		var w1, w2 int
		var d1, d2 bool

		// read rune from former string available
		if i < len(s1) {
			r1, w1 = utf8.DecodeRuneInString(s1[i:])
			i += w1

			// if digit, accumulate
			if d1 = ('0' <= r1 && r1 <= '9'); d1 {
				n1 = append(n1, r1)
			}
		}

		// read rune from latter string if available
		if j < len(s2) {
			r2, w2 = utf8.DecodeRuneInString(s2[j:])
			j += w2

			// if digit, accumulate
			if d2 = ('0' <= r2 && r2 <= '9'); d2 {
				n2 = append(n2, r2)
			}
		}

		// if have rune and other non-digit rune
		if (!d1 || !d2) && r1 > 0 && r2 > 0 {
			if len(n1) > 0 && len(n2) > 0 {
				// compare digits in accumulators
				less, equal := compareByDigits(n1, n2)
				if !equal {
					return less
				}

				// if equal, empty accumulators and continue
				n1, n2 = n1[0:0], n2[0:0]
			}

			// detect if non-digit rune from former or latter
			if r1 != r2 {
				return r1 < r2
			}
		}
	}

	// reached both strings ends, compare numeric accumulators
	less, equal := compareByDigits(n1, n2)

	if !equal {
		return less
	}

	// last hope
	return len(s1) < len(s2)
}

// Compare two numeric fields by their digits
func compareByDigits(n1, n2 []rune) (less, equal bool) {
	offset := len(n2) - len(n1)
	n1n2 := offset < 0 // len(n1) > len(n2)
	if n1n2 {
		// if n1 longer, inverse with n2
		offset = -offset
		n1, n2 = n2, n1
	}

	var j int
	for i := range n2 {
		var r1 rune
		if offset == 0 {
			// begin actual read
			r1 = n1[j]
			j++
		} else {
			// emulate zero-padding
			r1 = '0'
			offset--
		}

		r2 := n2[i]
		if r1 != r2 {
			if n1n2 {
				return r2 < r1, false // actually r1 < r2
			}
			return r1 < r2, false
		}
	}

	// use overall length then
	if n1n2 {
		return true, false
	}
	return !n1n2, len(n1) == len(n2)
}
