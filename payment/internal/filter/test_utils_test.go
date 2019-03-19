package filter

import (
	"testing"

	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/stretchr/testify/require"
)

// filterLeaf is a FilterLeaf implementation for the purpose of the tests
type filterLeaf struct {
	Name string
	Val  string
	Cmp  payment.FilterCmp
}

func (f filterLeaf) Filter() (payment.FilterCmp, interface{}) {
	return f.Cmp, f.Val
}

func (f filterLeaf) IsSet() bool {
	return f.Name != "" && f.Val != ""
}

// filterLeafNil is a FilterLeaf implementation which can contains nil values
// for the purpose of the tests
type filterLeafNil struct {
	Name string
	Val  interface{}
	Cmp  payment.FilterCmp
}

func (f filterLeafNil) Filter() (payment.FilterCmp, interface{}) {
	return f.Cmp, f.Val
}

func (f filterLeafNil) IsSet() bool {
	return f.Name != ""
}

func stringCmpOp(op payment.FilterCmp) string {
	switch op {
	case payment.FilterCmpEqual:
		return "=="
	case payment.FilterCmpGreaterOrEqualThan:
		return ">="
	case payment.FilterCmpGreaterThan:
		return ">"
	case payment.FilterCmpLessOrEqualThan:
		return "<="
	case payment.FilterCmpLessThan:
		return "<"
	case payment.FilterCmpMatch:
		return "MATCH"
	case payment.FilterCmpNotEqual:
		return "!="
	}

	return ""
}

func stringLogicalOp(op payment.FilterLogical) string {
	switch op {
	case payment.FilterLogicalAnd:
		return "&&"
	case payment.FilterLogicalOr:
		return "||"
	}

	return ""
}

func newFilterFromLeaf(t *testing.T, l payment.FilterLeaf) payment.Filter {
	f, err := payment.NewFilterFromLeaf(l)
	require.NoError(t, err)

	return f
}
