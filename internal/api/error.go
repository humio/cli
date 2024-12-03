package api

import (
	"fmt"
)

type EntityType string

const (
	EntityTypeSearchDomain    EntityType = "search-domain"
	EntityTypeRepository      EntityType = "repository"
	EntityTypeView            EntityType = "view"
	EntityTypeIngestToken     EntityType = "ingest-token"
	EntityTypeParser          EntityType = "parser"
	EntityTypeAction          EntityType = "action"
	EntityTypeAlert           EntityType = "alert"
	EntityTypeFilterAlert     EntityType = "filter-alert"
	EntityTypeScheduledSearch EntityType = "scheduled-search"
	EntityTypeAggregateAlert  EntityType = "aggregate-alert"
	EntityTypeUser            EntityType = "user"
	EntityTypeClusterNode     EntityType = "cluster-node"
)

func (e EntityType) String() string {
	return string(e)
}

type EntityNotFound struct {
	entityType EntityType
	key        string
}

func (e EntityNotFound) EntityType() EntityType {
	return e.entityType
}

func (e EntityNotFound) Key() string {
	return e.key
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("%s %q not found", e.entityType.String(), e.key)
}

func SearchDomainNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeSearchDomain,
		key:        name,
	}
}

func RepositoryNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeRepository,
		key:        name,
	}
}

func ViewNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeView,
		key:        name,
	}
}

func IngestTokenNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeIngestToken,
		key:        name,
	}
}

func ParserNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeParser,
		key:        name,
	}
}

func ActionNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAction,
		key:        name,
	}
}

func AlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAlert,
		key:        name,
	}
}

func FilterAlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeFilterAlert,
		key:        name,
	}
}

func ScheduledSearchNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeScheduledSearch,
		key:        name,
	}
}

func AggregateAlertNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeAggregateAlert,
		key:        name,
	}
}

func UserNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeUser,
		key:        name,
	}
}

func ClusterNodeNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeClusterNode,
		key:        name,
	}
}
