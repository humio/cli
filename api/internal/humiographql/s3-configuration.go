package humiographql

import (
	"fmt"
	"strings"
)

type S3Configuration struct {
	Bucket   string            `graphql:"bucket"`
	Region   string            `graphql:"region"`
	Disabled bool              `graphql:"disabled"`
	Format   S3ArchivingFormat `graphql:"format"`
}

// IsEnabled - determine if S3Configuration is enabled based on values and the Disabled field
// to avoid a bool defaulting to false
func (s *S3Configuration) IsEnabled() bool {
	if !s.IsConfigured() {
		return false
	}
	return !s.Disabled
}

func (s *S3Configuration) IsConfigured() bool {
	if s.Bucket != "" && s.Region != "" && s.Format != "" {
		return true
	}
	return false
}

// S3ArchivingFormat - the format in which to store the archived data on S3, either RAW or NDJSON
type S3ArchivingFormat string

const DefaultS3ArchivingFormat = S3ArchivingFormat("NDSON")

var ValidS3ArchivingFormats = []S3ArchivingFormat{"NDJSON", "RAW"}

// NewS3ArchivingFormat - creates a S3ArchivingFormat and ensures the value is uppercase
func NewS3ArchivingFormat(format string) (S3ArchivingFormat, error) {
	f := S3ArchivingFormat(strings.ToUpper(string(format)))
	err := f.Validate()
	if err != nil {
		return "", err
	}
	return f, nil
}

func (f *S3ArchivingFormat) Validate() error {
	for _, v := range ValidS3ArchivingFormats {
		if v == *f {
			return nil
		}
	}
	return fmt.Errorf("invalid S3 archiving format. Valid formats: %s", ValidS3ArchivingFormats)
}
