package api

import (
	"fmt"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type EntityType string

const (
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeSearchDomain EntityType = "search-domain"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeRepository EntityType = "repository"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeView EntityType = "view"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeIngestToken EntityType = "ingest-token"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeParser EntityType = "parser"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeAction EntityType = "action"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeAlert EntityType = "alert"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeFilterAlert EntityType = "filter-alert"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeScheduledSearch EntityType = "scheduled-search"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	EntityTypeAggregateAlert EntityType = "aggregate-alert"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (e EntityType) String() string {
	return string(e)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type EntityNotFound struct {
	entityType EntityType
	key        string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (e EntityNotFound) EntityType() EntityType {
	return e.entityType
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (e EntityNotFound) Key() string {
	return e.key
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (e EntityNotFound) Error() string {
	return fmt.Sprintf("%s %q not found", e.entityType.String(), e.key)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func SearchDomainNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeSearchDomain,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func RepositoryNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeRepository,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func ViewNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeView,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func IngestTokenNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeIngestToken,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func ParserNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeParser,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func ActionNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAction,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func AlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAlert,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func FilterAlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeFilterAlert,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func ScheduledSearchNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeScheduledSearch,
		key:        name,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func AggregateAlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAggregateAlert,
		key:        name,
	}
}
