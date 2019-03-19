package filter

import (
	"fmt"

	"github.com/ifraixedes/go-payments-api-example/payment"
)

// String returns a string representation of f, using lf, cs and ls but adding
// parentheses between each pair of nodes and operation for canceling any
// operator precedence, when there are more than 2 nodes.
func String(f payment.Filter, lf LeafField, cs CmpString, ls LogicalString) string {
	var s = stringify{
		leafField:  lf,
		cmpStr:     cs,
		logicalStr: ls,
		root:       f,
	}

	return s.String(f)
}

// StringValues returns a string representation of f and the list of the values
// of the filter in the order that they appear in string, replacing each node
// value by the ph as a placeholder and using lf, cs and ls but adding
// parentheses between each pair of nodes and operation for canceling any
// operator precedence, when there are more than 2 nodes.
//
// You can think to use this function when you have f represents "id = 10", then
// you will get "id = ?", []interface{}{10}, when using "?" ph, and you need
// such filter to pass it, for example to a SQL DB in a SELECT statement, hence
// if 'fs' is the returned string and 'vals' is the returned interface{} slice,
// you could use it with the standard SQL like:
//
// db.QueryRow("SELECT * FROM users WHERE " + fs, vals)
func StringValues(f payment.Filter, ph string, lf LeafField, cs CmpString, ls LogicalString) (string, []interface{}) {
	var s = stringify{
		leafField:  lf,
		cmpStr:     cs,
		logicalStr: ls,
		root:       f,
	}

	return s.StringValues(f, ph)
}

type stringify struct {
	leafField  LeafField
	cmpStr     CmpString
	logicalStr LogicalString
	root       payment.Filter
}

func (s stringify) String(f payment.Filter) string {
	switch f.NodeType() {
	case payment.FilterNodeTypeEmpty:
		return ""
	case payment.FilterNodeTypeLeaf:
		l := f.Leaf()
		op, val := l.Filter()
		cs, ok := s.cmpStr(op, val)
		if ok {
			return fmt.Sprintf("%s %s", s.leafField(l), cs)
		}

		return fmt.Sprintf("%s %s %s", s.leafField(l), cs, val)
	}

	var (
		op, l, r = f.Nodes()
		ls       = s.String(l)
		rs       = s.String(r)
	)

	if f == s.root {
		return fmt.Sprintf("%s %s %s", ls, s.logicalStr(op), rs)
	}

	return fmt.Sprintf("(%s %s %s)", ls, s.logicalStr(op), rs)
}

func (s stringify) StringValues(f payment.Filter, ph string) (string, []interface{}) {
	switch f.NodeType() {
	case payment.FilterNodeTypeEmpty:
		return "", nil
	case payment.FilterNodeTypeLeaf:
		l := f.Leaf()
		op, val := l.Filter()
		cs, ok := s.cmpStr(op, val)
		if ok {
			return fmt.Sprintf("%s %s", s.leafField(l), cs), nil
		}

		return fmt.Sprintf("%s %s %s", s.leafField(l), cs, ph), []interface{}{val}
	}

	var (
		op, l, r = f.Nodes()
		ls, lv   = s.StringValues(l, ph)
		rs, rv   = s.StringValues(r, ph)
	)

	if f == s.root {
		return fmt.Sprintf("%s %s %s", ls, s.logicalStr(op), rs), append(lv, rv...)
	}

	return fmt.Sprintf("(%s %s %s)", ls, s.logicalStr(op), rs), append(lv, rv...)
}
