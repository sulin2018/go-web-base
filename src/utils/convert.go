package utils

import "strconv"

type StrTo string

func (f StrTo) Exist() bool {
	// 0x1E = 30  RS  (record separator)
	return string(f) != string(0x1E)
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}

func (f StrTo) Uint() uint {
	v, _ := strconv.ParseUint(f.String(), 10, 64)
	return uint(v)
}
