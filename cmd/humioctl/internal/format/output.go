package format

import (
	"encoding/json"
	"fmt"
)

type YesNo bool

func (y YesNo) String() string {
	if y {
		return "yes"
	}
	return "no"
}

func (y YesNo) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(y))
}

type Checkmark bool

func (c Checkmark) String() string {
	if c {
		return "✓"
	}
	return ""
}

func (c Checkmark) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(c))
}

type ValueOrEmpty string

func (v ValueOrEmpty) String() string {
	if string(v) == "" {
		return "-"
	}
	return string(v)
}

func (v ValueOrEmpty) MarshalJSON() ([]byte, error) {
	if string(v) == "" {
		return []byte("null"), nil
	}

	return json.Marshal(string(v))
}

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
