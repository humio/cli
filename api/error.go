package api

import "fmt"

type EntityType string

const (
	EntityTypeParser EntityType = "parser"
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

func ParserNotFound(name string) error {
	return EntityNotFound{
		entityType: EntityTypeParser,
		key:        name,
	}
}
