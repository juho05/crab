package interpreter

import (
	"fmt"
	"strings"
)

type list []any

func (l list) String() string {
	text := "["

	for _, v := range l {
		text = fmt.Sprintf("%s%v", text, v)
		text = fmt.Sprintf("%s,", text)
	}

	text = strings.TrimSuffix(text, ",")

	text = text + "]"
	return text
}

func (l list) equals(other list) bool {
	if len(l) != len(other) {
		return false
	}

	for i, v := range l {
		if !areEqual(v, other[i]) {
			return false
		}
	}

	return true
}
