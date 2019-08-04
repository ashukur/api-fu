package validator

import (
	"fmt"

	"github.com/ccbrown/api-fu/graphql/ast"
	"github.com/ccbrown/api-fu/graphql/schema"
)

type Error struct {
	Message string

	// If a validator is unable to perform its job due to an error unrelated to its purpose, it will
	// emit a secondary error. Secondary errors are always errors that should be caught by other
	// validators, so if there are any primary errors, secondary errors are discarded as they should
	// all be duplicates. If a secondary error makes it out of validation, there's probably a
	// mistake in one of the validators.
	IsSecondary bool
}

func (err *Error) Error() string {
	return err.Message
}

func newError(message string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(message, args...),
	}
}

func newSecondaryError(message string, args ...interface{}) *Error {
	return &Error{
		Message:     fmt.Sprintf(message, args...),
		IsSecondary: true,
	}
}

func ValidateDocument(doc *ast.Document, s *schema.Schema) []*Error {
	typeInfo := newTypeInfo(doc, s)
	var errs []*Error
	for _, f := range []func(*ast.Document, *schema.Schema, *TypeInfo) []*Error{
		validateDocument,
		validateOperations,
		validateFields,
		validateArguments,
		validateFragments,
		validateDirectives,
		validateVariables,
	} {
		errs = append(errs, f(doc, s, typeInfo)...)
	}
	var primary []*Error
	for _, err := range errs {
		if !err.IsSecondary {
			primary = append(primary, err)
		}
	}
	if len(primary) > 0 {
		return primary
	}
	return errs
}

func unwrappedASTType(t ast.Type) *ast.NamedType {
	for {
		if t == nil {
			return nil
		}
		switch tt := t.(type) {
		case *ast.ListType:
			t = tt.Type
		case *ast.NonNullType:
			t = tt.Type
		case *ast.NamedType:
			return tt
		default:
			panic(fmt.Sprintf("unsupported ast type: %T", t))
		}
	}
}
