package humiographql

import "time"

func ConvertDaysToMillis(days int64) int64 {
	milliseconds := int64(time.Duration(days) * 24 * time.Hour / time.Millisecond)
	return milliseconds
}

func ConvertGBToBytes(gb float64) float64 {
	// Data sizes in LogScale are expressed in SI units using decimal (Base 10).
	const bytesInGB = 1e9 // 1 GB = 1,000,000,000 bytes (SI units)
	return gb * bytesInGB
}
