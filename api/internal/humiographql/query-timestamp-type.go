package humiographql

type QueryTimestampType string

const (
	QueryTimestampTypeIngestTimestamp QueryTimestampType = "IngestTimestamp"
	QueryTimestampTypeEventTimestamp  QueryTimestampType = "EventTimestamp"
)
