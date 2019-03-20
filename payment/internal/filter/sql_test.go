package filter

import (
	"testing"

	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQL(t *testing.T) {
	var leafField = func(fl payment.FilterLeaf) string {
		l, ok := fl.(filterLeaf)
		if ok {
			return l.Name
		}

		ln := fl.(filterLeafNil)
		return ln.Name
	}

	t.Run("empty", func(t *testing.T) {
		s, vals := SQL(payment.Filter{}, leafField)
		assert.Empty(t, s)
		assert.Empty(t, vals)
	})

	t.Run("1 node", func(t *testing.T) {
		f := newFilterFromLeaf(t, filterLeaf{Name: "l", Val: "v", Cmp: payment.FilterCmpNotEqual})
		s, vals := SQL(f, leafField)
		assert.Equal(t, "l <> ?", s)
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

		s, vals := SQL(f, leafField)
		assert.Equal(t, "l1 >= ? AND l2 < ?", s)
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

			f3 := newFilterFromLeaf(t, filterLeafNil{Name: "l3", Val: nil, Cmp: payment.FilterCmpNotEqual})
			f, err = payment.NewFilter(payment.FilterLogicalAnd, f, f3)
			require.NoError(t, err)
		}

		s, vals := SQL(f, leafField)
		assert.Equal(t, "(l1 > ? OR l2 < ?) AND l3 IS NOT NULL", s)
		assert.Equal(t, []interface{}{"v1", "v2"}, vals)
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

		s, vals := SQL(f, leafField)
		assert.Equal(t, "(l1 LIKE ? AND l2 <= ?) OR (l3 = ? AND l4 > ?)", s)
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

		s, vals := SQL(f, leafField)
		assert.Equal(t, "((l1 LIKE ? AND l2 <= ?) OR l1 LIKE ?) OR (l3 = ? AND l4 > ?)", s)
		assert.Equal(t, []interface{}{"v1", "v2", "v1", "v3", "v4"}, vals)
	})

	t.Run("6 nodes", func(t *testing.T) {
		var f payment.Filter
		{
			f1 := newFilterFromLeaf(t, filterLeaf{Name: "l1", Val: "v1", Cmp: payment.FilterCmpMatch})
			f2 := newFilterFromLeaf(t, filterLeafNil{Name: "l2", Cmp: payment.FilterCmpEqual})
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

		s, vals := SQL(f, leafField)
		assert.Equal(t, "(l3 = ? AND l4 > ?) OR ((l3 = ? AND l4 > ?) AND (l1 LIKE ? AND l2 IS NULL))", s)
		assert.Equal(t, []interface{}{"v3", "v4", "v3", "v4", "v1"}, vals)
	})
}
