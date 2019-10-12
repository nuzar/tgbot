package main

import (
	"fmt"
	"strconv"
)

// UnionIntString union type of int and string
type UnionIntString struct {
	int64
	string
}

func (u UnionIntString) IsInt() bool {
	if u.int64 != 0 {
		return true
	}

	if u.string != "" {
		return false
	}

	return true
}

// MarshalJSON implement json.Marshaler interface
// marshal int first
func (u UnionIntString) MarshalJSON() ([]byte, error) {
	if !u.IsInt() {
		return []byte("\"" + u.string + "\""), nil

	}

	return []byte(strconv.FormatInt(u.int64, 10)), nil
}

// UnmarshalJSON implement  json.Unmarshaler interface
func (u *UnionIntString) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		u.string = string(data[1 : len(data)-1])
		return nil
	default:
		if n, err := strconv.ParseInt(string(data), 10, 64); err != nil {
			return fmt.Errorf("%s is not UnionIntString", string(data))
		} else {
			u.int64 = n
		}
	}

	return nil
}

func isHTTPStatusOK(code int) bool {
	return code/200 == 1
}
