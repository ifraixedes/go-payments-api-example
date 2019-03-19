package filter

import (
	"testing"

	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringify_String(t *testing.T) {
	var newLeafField = func(calls *[]payment.FilterLeaf) LeafField {
		return func(fl payment.FilterLeaf) string {
			c := append(*calls, fl)
			*calls = c
			return ""
		}
	}

	var newCmpStr = func(calls *[]payment.FilterCmp) CmpString {
		return func(cmp payment.FilterCmp, _ interface{}) (string, bool) {
			c := append(*calls, cmp)
			*calls = c
			return "", false
		}
	}

	var newLogicalStr = func(calls *[]payment.FilterLogical) LogicalString {
		return func(l payment.FilterLogical) string {
			c := append(*calls, l)
			*calls = c
			return ""
		}
	}

	t.Run("empty", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		s.String(payment.Filter{})
		assert.Empty(t, tlc)
		assert.Empty(t, tcmpc)
		assert.Empty(t, tloc)
	})

	t.Run("1 node", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f = newFilterFromLeaf(t, filterLeaf{Name: "l", Val: "v", Cmp: payment.FilterCmpEqual})
		s.String(f)
		assert.Len(t, tlc, 1)
		assert.Len(t, tcmpc, 1)
		assert.Empty(t, tloc)
	})

	t.Run("2 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterOrEqualThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})

			var err error
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)
		}

		s.String(f)
		assert.Len(t, tlc, 2)
		assert.Len(t, tcmpc, 2)
		assert.Len(t, tloc, 1)
	})

	t.Run("3 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			var err error

			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})
			f, err = payment.NewFilter(payment.FilterLogicalOr, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f, f3)
			require.NoError(t, err)
		}

		s.String(f)
		assert.Len(t, tlc, 3)
		assert.Len(t, tcmpc, 3)
		assert.Len(t, tloc, 2)
	})

	t.Run("4 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f12, f34)
			require.NoError(t, err)
		}

		s.String(f)
		assert.Len(t, tlc, 4)
		assert.Len(t, tcmpc, 4)
		assert.Len(t, tloc, 3)
	})

	t.Run("5 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f121, err := payment.NewFilter(payment.FilterLogicalOr, f12, f1)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f121, f34)
			require.NoError(t, err)
		}

		s.String(f)
		assert.Len(t, tlc, 5)
		assert.Len(t, tcmpc, 5)
		assert.Len(t, tloc, 4)
	})

	t.Run("6 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f3412, err := payment.NewFilter(payment.FilterLogicalAnd, f34, f12)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f34, f3412)
			require.NoError(t, err)
		}

		s.String(f)
		assert.Len(t, tlc, 6)
		assert.Len(t, tcmpc, 6)
		assert.Len(t, tloc, 5)
	})
}

func TestStringify_StringValues(t *testing.T) {
	var newLeafField = func(calls *[]payment.FilterLeaf) LeafField {
		return func(fl payment.FilterLeaf) string {
			c := append(*calls, fl)
			*calls = c
			return ""
		}
	}

	var newCmpStr = func(calls *[]payment.FilterCmp) CmpString {
		return func(cmp payment.FilterCmp, _ interface{}) (string, bool) {
			c := append(*calls, cmp)
			*calls = c
			return "", false
		}
	}

	var newLogicalStr = func(calls *[]payment.FilterLogical) LogicalString {
		return func(l payment.FilterLogical) string {
			c := append(*calls, l)
			*calls = c
			return ""
		}
	}

	t.Run("empty", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		_, vals := s.StringValues(payment.Filter{}, "")
		assert.Empty(t, tlc)
		assert.Empty(t, tcmpc)
		assert.Empty(t, tloc)
		assert.Empty(t, vals)
	})

	t.Run("1 node", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f = newFilterFromLeaf(t, filterLeaf{Name: "l", Val: "v", Cmp: payment.FilterCmpEqual})
		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 1)
		assert.Len(t, tcmpc, 1)
		assert.Empty(t, tloc)
		assert.Len(t, vals, 1)
	})

	t.Run("2 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterOrEqualThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})

			var err error
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)
		}

		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 2)
		assert.Len(t, tcmpc, 2)
		assert.Len(t, tloc, 1)
		assert.Len(t, vals, 2)
	})

	t.Run("3 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			var err error

			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})
			f, err = payment.NewFilter(payment.FilterLogicalOr, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f, f3)
			require.NoError(t, err)
		}

		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 3)
		assert.Len(t, tcmpc, 3)
		assert.Len(t, tloc, 2)
		assert.Len(t, vals, 3)
	})

	t.Run("4 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f12, f34)
			require.NoError(t, err)
		}

		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 4)
		assert.Len(t, tcmpc, 4)
		assert.Len(t, tloc, 3)
		assert.Len(t, vals, 4)
	})

	t.Run("5 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f121, err := payment.NewFilter(payment.FilterLogicalOr, f12, f1)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f121, f34)
			require.NoError(t, err)
		}

		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 5)
		assert.Len(t, tcmpc, 5)
		assert.Len(t, tloc, 4)
		assert.Len(t, vals, 5)
	})

	t.Run("6 nodes", func(t *testing.T) {
		var (
			tlc   []payment.FilterLeaf
			tcmpc []payment.FilterCmp
			tloc  []payment.FilterLogical
		)

		var s = stringify{
			leafField:  newLeafField(&tlc),
			cmpStr:     newCmpStr(&tcmpc),
			logicalStr: newLogicalStr(&tloc),
		}

		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f3412, err := payment.NewFilter(payment.FilterLogicalAnd, f34, f12)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f34, f3412)
			require.NoError(t, err)
		}

		_, vals := s.StringValues(f, "")
		assert.Len(t, tlc, 6)
		assert.Len(t, tcmpc, 6)
		assert.Len(t, tloc, 5)
		assert.Len(t, vals, 6)
	})
}

func TestString(t *testing.T) {
	var cmpStr = CmpStringNoVal(stringCmpOp)
	var leafField = func(fl payment.FilterLeaf) string {
		l := fl.(filterLeaf)
		return l.Name
	}

	t.Run("empty", func(t *testing.T) {
		s := String(payment.Filter{}, leafField, cmpStr, stringLogicalOp)
		assert.Empty(t, s)
	})

	t.Run("1 node", func(t *testing.T) {
		f := newFilterFromLeaf(t, filterLeaf{Name: "l", Val: "v", Cmp: payment.FilterCmpEqual})
		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "l == v", s)
	})

	t.Run("2 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterOrEqualThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})

			var err error
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)
		}

		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "l1 >= v1 && l2 < v2", s)
	})

	t.Run("3 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			var err error

			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})
			f, err = payment.NewFilter(payment.FilterLogicalOr, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f, f3)
			require.NoError(t, err)
		}

		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l1 > v1 || l2 < v2) && l3 == v3", s)
	})

	t.Run("4 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f12, f34)
			require.NoError(t, err)
		}

		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l1 MATCH v1 && l2 <= v2) || (l3 == v3 && l4 > v4)", s)
	})

	t.Run("5 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f121, err := payment.NewFilter(payment.FilterLogicalOr, f12, f1)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f121, f34)
			require.NoError(t, err)
		}

		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "((l1 MATCH v1 && l2 <= v2) || l1 MATCH v1) || (l3 == v3 && l4 > v4)", s)
	})

	t.Run("6 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f3412, err := payment.NewFilter(payment.FilterLogicalAnd, f34, f12)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f34, f3412)
			require.NoError(t, err)
		}

		s := String(f, leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l3 == v3 && l4 > v4) || ((l3 == v3 && l4 > v4) && (l1 MATCH v1 && l2 <= v2))", s)
	})
}

func TestStringValues(t *testing.T) {
	var cmpStr = CmpStringNoVal(stringCmpOp)
	var leafField = func(fl payment.FilterLeaf) string {
		l := fl.(filterLeaf)
		return l.Name
	}

	t.Run("empty", func(t *testing.T) {
		s, vals := StringValues(payment.Filter{}, "$", leafField, cmpStr, stringLogicalOp)
		assert.Empty(t, s)
		assert.Empty(t, vals)
	})

	t.Run("1 node", func(t *testing.T) {
		f := newFilterFromLeaf(t, filterLeaf{Name: "l", Val: "v", Cmp: payment.FilterCmpEqual})
		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "l == $", s)
		assert.Equal(t, []interface{}{"v"}, vals)
	})

	t.Run("2 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterOrEqualThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})

			var err error
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)
		}

		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "l1 >= $ && l2 < $", s)
		assert.Equal(t, []interface{}{"v1", "v2"}, vals)
	})

	t.Run("3 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			var err error

			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpGreaterThan})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessThan})
			f, err = payment.NewFilter(payment.FilterLogicalOr, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f, f3)
			require.NoError(t, err)
		}

		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l1 > $ || l2 < $) && l3 == $", s)
		assert.Equal(t, []interface{}{"v1", "v2", "v3"}, vals)
	})

	t.Run("4 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f12, f34)
			require.NoError(t, err)
		}

		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l1 MATCH $ && l2 <= $) || (l3 == $ && l4 > $)", s)
		assert.Equal(t, []interface{}{"v1", "v2", "v3", "v4"}, vals)
	})

	t.Run("5 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f121, err := payment.NewFilter(payment.FilterLogicalOr, f12, f1)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f121, f34)
			require.NoError(t, err)
		}

		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "((l1 MATCH $ && l2 <= $) || l1 MATCH $) || (l3 == $ && l4 > $)", s)
		assert.Equal(t, []interface{}{"v1", "v2", "v1", "v3", "v4"}, vals)
	})

	t.Run("6 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeaf{Name: "l2", Val: "v2", Cmp: payment.FilterCmpLessOrEqualThan})
			f12, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
			require.NoError(t, err)

			f3 := newFilterFromLeaf(t, filterLeaf{Name: "l3", Val: "v3", Cmp: payment.FilterCmpEqual})
			f4 := newFilterFromLeaf(t, filterLeaf{Name: "l4", Val: "v4", Cmp: payment.FilterCmpGreaterThan})
			f34, err := payment.NewFilter(payment.FilterLogicalAnd, f3, f4)
			require.NoError(t, err)

			f3412, err := payment.NewFilter(payment.FilterLogicalAnd, f34, f12)
			require.NoError(t, err)

			f, err = payment.NewFilter(payment.FilterLogicalOr, f34, f3412)
			require.NoError(t, err)
		}

		s, vals := StringValues(f, "$", leafField, cmpStr, stringLogicalOp)
		assert.Equal(t, "(l3 == $ && l4 > $) || ((l3 == $ && l4 > $) && (l1 MATCH $ && l2 <= $))", s)
		assert.Equal(t, []interface{}{"v3", "v4", "v3", "v4", "v1", "v2"}, vals)
	})
}
