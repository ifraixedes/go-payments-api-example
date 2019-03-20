package filter

import "github.com/ifraixedes/go-payments-api-example/payment"

// SQL returns the string representation of the filter which can be used as a
// SQL standard WHERE expression and the list of values to fill each place
// holder (i.e. '?') present in the string.
func SQL(f payment.Filter, lf LeafField) (string, []interface{}) {
	var s = stringify{
		leafField:  lf,
		cmpStr:     sqlCmpOp,
		logicalStr: sqlLogicalOp,
		root:       f,
	}

	return s.StringValues(f, "?")
}

func sqlCmpOp(op payment.FilterCmp, val interface{}) (string, bool) {
	switch op {
	case payment.FilterCmpEqual:
		if val == nil {
			return "IS NULL", true
		}

		return "=", false
	case payment.FilterCmpNotEqual:
		if val == nil {
			return "IS NOT NULL", true
		}

		return "<>", false
	case payment.FilterCmpGreaterOrEqualThan:
		return ">=", false
	case payment.FilterCmpGreaterThan:
		return ">", false
	case payment.FilterCmpLessOrEqualThan:
		return "<=", false
	case payment.FilterCmpLessThan:
		return "<", false
	case payment.FilterCmpMatch:
		return "LIKE", false
	}

	return "", false
}

func sqlLogicalOp(op payment.FilterLogical) string {
	switch op {
	case payment.FilterLogicalAnd:
		return "AND"
	case payment.FilterLogicalOr:
		return "OR"
	}

	return ""
}
