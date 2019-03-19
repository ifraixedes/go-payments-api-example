package payment_test

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
)

func ExampleFilter() {
	var f payment.Filter
	{
		var fleft payment.Filter
		{
			var fage10, err = payment.NewFilterByAmount(payment.FilterCmpGreaterOrEqualThan, 10)
			fatalIfErr(err)

			fidne, err := payment.NewFilterByID(
				payment.FilterCmpNotEqual, uuid.FromStringOrNil("86b16b89-61c1-4f2f-963b-b542e3597d69"),
			)
			fatalIfErr(err)

			fal8_5, err := payment.NewFilterByAmount(payment.FilterCmpLessThan, 8.5)
			fatalIfErr(err)

			faid, err := payment.NewFilter(payment.FilterLogicalAnd, fage10, fidne)
			fatalIfErr(err)

			fleft, err = payment.NewFilter(payment.FilterLogicalOr, faid, fal8_5)
			fatalIfErr(err)
		}

		fright, err := payment.NewFilterByType(payment.FilterCmpEqual, "Payment")
		fatalIfErr(err)

		f, err = payment.NewFilter(payment.FilterLogicalAnd, fleft, fright)
		fatalIfErr(err)
	}

	// Given a filter this is a processor for printing out the string filter
	// expected by the example test
	var (
		cf                     = &f
		stk              stack = []payment.Filter{}
		closeParentheses       = false
		filterSQL        string
	)

	for { // Iteratively in-order traversal of the filter (which is a binary tree)
		if cf == nil {
			var pn, ok = stk.Pop()
			if !ok {
				break
			}

			if pn.NodeType() == payment.FilterNodeTypeLeaf {
				var lstr = toLeaf(pn.Leaf())
				if closeParentheses {
					filterSQL = fmt.Sprintf("%s%s)", filterSQL, lstr)
					closeParentheses = false
				} else {
					filterSQL = fmt.Sprintf("%s%s", filterSQL, lstr)
				}
			}

			if pn.NodeType() == payment.FilterNodeTypeNonLeaf {
				if len(stk) > 0 {
					filterSQL = fmt.Sprintf("(%s", filterSQL)
					closeParentheses = true
				}

				var op, _, right = pn.Nodes()
				filterSQL = fmt.Sprintf("%s %s ", filterSQL, toLogicalOp(op))
				cf = &right
			}

			continue
		}

		if cf.NodeType() == payment.FilterNodeTypeNonLeaf {
			_, left, _ := cf.Nodes()
			stk.Push(*cf)
			cf = &left
		} else {
			stk.Push(*cf)
			cf = nil
		}
	}

	fmt.Println(filterSQL)
	// Output:
	// ((amount >= 10.00 AND id != '86b16b89-61c1-4f2f-963b-b542e3597d69') OR amount < 8.50) AND type = 'Payment'
}

type stack []payment.Filter

func (s *stack) Push(f payment.Filter) {
	var ts = append(*s, f)
	*s = ts
}

func (s *stack) Pop() (payment.Filter, bool) {
	if ts := *s; len(ts) > 0 {
		var f = ts[len(ts)-1]
		ts = ts[:len(ts)-1]

		*s = ts
		return f, true
	}

	return payment.Filter{}, false
}

func toCmpOp(op payment.FilterCmp) string {
	switch op {
	case payment.FilterCmpEqual:
		return "="
	case payment.FilterCmpGreaterOrEqualThan:
		return ">="
	case payment.FilterCmpGreaterThan:
		return ">"
	case payment.FilterCmpLessOrEqualThan:
		return "<="
	case payment.FilterCmpLessThan:
		return "<"
	case payment.FilterCmpMatch:
		return "LIKE"
	case payment.FilterCmpNotEqual:
		return "!="
	}

	fatalIfErr(fmt.Errorf("Unrecognized comparison operator: %d", op))
	return ""
}

func toLogicalOp(op payment.FilterLogical) string {
	switch op {
	case payment.FilterLogicalAnd:
		return "AND"
	case payment.FilterLogicalOr:
		return "OR"
	}

	fatalIfErr(fmt.Errorf("Unrecognized logical operator: %d", op))
	return ""
}

func toLeafName(f payment.FilterLeaf) string {
	switch f.(type) {
	case payment.FilterLeafAmount:
		return "amount"
	case payment.FilterLeafType:
		return "type"
	case payment.FilterLeafID:
		return "id"
	}

	fatalIfErr(fmt.Errorf("Unrecognized leaf type: %T", f))
	return ""
}

func toLeafValue(f payment.FilterLeaf) string {
	switch f.(type) {
	case payment.FilterLeafAmount:
		_, val := f.Filter()
		return fmt.Sprintf("%.2f", val)
	case payment.FilterLeafType:
		_, val := f.Filter()
		return fmt.Sprintf("'%s'", val)
	case payment.FilterLeafID:
		_, val := f.Filter()
		id := val.(uuid.UUID)
		return fmt.Sprintf("'%s'", id.String())
	}

	fatalIfErr(fmt.Errorf("Unrecognized leaf type: %T", f))
	return ""
}

func toLeaf(f payment.FilterLeaf) string {
	var op, _ = f.Filter()
	return fmt.Sprintf("%s %s %s", toLeafName(f), toCmpOp(op), toLeafValue(f))
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
