package main

import (
	"encoding/json"
	"fmt"
)

type yesNo bool

func (y yesNo) String() string {
	if y {
		return "yes"
	}
	return "no"
}

func (y yesNo) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(y))
}

type checkmark bool

func (c checkmark) String() string {
	if c {
		return "✓"
	}
	return ""
}

func (c checkmark) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(c))
}

type valueOrEmpty string

func (v valueOrEmpty) String() string {
	if string(v) == "" {
		return "-"
	}
	return string(v)
}

func (v valueOrEmpty) MarshalJSON() ([]byte, error) {
	if string(v) == "" {
		return []byte("null"), nil
	}

	return json.Marshal(string(v))
}

//func valueOrEmpty(v string) string {
//	if v == "" {
//		return "-"
//	}
//	return v
//}

type ByteCountDecimal int64

func (b ByteCountDecimal) String() string {
	const unit = 1000
	if int64(b) < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func (b ByteCountDecimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(b))
}
