package apifu

import (
	"context"
	"reflect"

	"github.com/ccbrown/api-fu/graphql"
)

type Node struct {
	// TypeId should be an integer that uniquely identifies the node type. Once set, it should never
	// change and no other nodes should ever use the same id.
	TypeId int

	Model reflect.Type

	// GetByIds should be a function that accepts a context and slice of ids and returns a slice of
	// models.
	GetByIds func(ctx context.Context, ids interface{}) (models interface{}, err error)

	Fields map[string]*NodeField
}

type NodeField struct {
	Resolve func(graphql.FieldContext) (interface{}, error)
}

func NonNullString(modelField string) *NodeField {
	return &NodeField{}
}
