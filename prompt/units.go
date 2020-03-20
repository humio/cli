package prompt

func AddSISuffix(v float64, binary bool) (val float64, suffix string) {
	var unit float64
	if binary {
		unit = 1024
	} else {
		unit = 1000
	}

	div, exp := unit, 0
	for n := v / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return v / div, string("kMGTPE"[exp])
}
