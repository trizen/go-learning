package substr

/*

# Author: Trizen
# License: GPLv3
# Date: 01 January 2013
# Website; http://trizen.googlecode.com

# Extracts a substring out of a string.

*/

import (
	"strings"
)

// It substracts from string 's' anything specified by positions.
// There are two available positions:
// start   - where to begin to extract;
// length  - how many characters to extract;
// If the length is omited, it will return the rest of the string.
//
// Examples:
//      strings.Substr(s, -5)    # return last five characters
//      strings.Substr(s, 3, 5)  # from pos 3, return the next five characters
//      strings.Substr(s, 2, -4) # from pos 2, return till the last four chars
func Substr(s string, pos ...int) string {
	chars := strings.Split(s, "")
	str_len := len(chars)

	if len(pos) == 0 { // when no pos specified, return s
		return s
	}

	if pos[0] < 0 {
		pos[0] = str_len + pos[0]
	}

	if len(pos) > 1 {
		if pos[1] < 0 {
			pos[1] = str_len + pos[1]
		} else {
			pos[1] = pos[0] + pos[1]
		}
	} else {
		pos = append(pos, str_len)
	}

	return strings.Join(chars[pos[0]:pos[1]], "")
}
