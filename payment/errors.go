package payment

import (
	"fmt"

	"go.fraixed.es/errors"
)

type code uint8

// The list of error codes that the payment service can return.
const (
	ErrInvalidArgFilterCmpNotExists code = iota + 1
	ErrInvalidArgFilterCmpNotSupported
	ErrInvalidArgFilterLeafNoValSet
	ErrInvalidArgFilterLogicalOpNotExists
	ErrInvalidArgFilterNodeEmpty
	ErrInvalidArgFilterValue
)

func (c code) String() string {
	switch c {
	case ErrInvalidArgFilterCmpNotExists:
		return "InvalidArgFilterCmpNotExists"
	case ErrInvalidArgFilterCmpNotSupported:
		return "InvalidArgFilterCmpNotSupported"
	case ErrInvalidArgFilterLeafNoValSet:
		return "InvalidArgFilterLeafNoValSet"
	case ErrInvalidArgFilterLogicalOpNotExists:
		return "InvalidArgFilterLogicalOp"
	case ErrInvalidArgFilterNodeEmpty:
		return "InvalidArgFilterNodeEmpty"
	case ErrInvalidArgFilterValue:
		return "InvalidArgFilterValue"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrInvalidArgFilterCmpNotExists:
		return "The filter comparison operator doesn't exist"
	case ErrInvalidArgFilterCmpNotSupported:
		return "The filter comparison operator isn't supported for this type leaf node"
	case ErrInvalidArgFilterLeafNoValSet:
		return "The filter leaf node is invalid because its value isn't set"
	case ErrInvalidArgFilterLogicalOpNotExists:
		return "The filter logical operator doesn't exist"
	case ErrInvalidArgFilterNodeEmpty:
		return "The filter node cannot be an empty"
	case ErrInvalidArgFilterValue:
		return "The filter value isn't a valid one"
	}

	return ""
}

// ErrMDArg creates a new  MD (metadata) from an function argument.
func ErrMDArg(name string, val interface{}) errors.MD {
	return errors.MD{
		K: fmt.Sprintf("arg:%s", name),
		V: val,
	}
}
