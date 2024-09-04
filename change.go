// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

const (
	added     = '+'
	removed   = '-'
	unchanged = ' '
	nested    = '^'
)

func isChangeSet(currentSign rune, upcomingSign rune) bool {
	if currentSign == added && upcomingSign == removed {
		return true
	} else if upcomingSign == added && currentSign == removed {
		return true
	} else {
		return false
	}
}
