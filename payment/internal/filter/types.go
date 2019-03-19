package filter

import "github.com/ifraixedes/go-payments-api-example/payment"

// ToLeaf is a function which returns a specific string representation of fl.
// type ToLeaf func(fl payment.FilterLeaf) string

// LeafField returns the specific field name for fl. Each filter consumer (e.g.
// a DB vendor) must provide a function which returns its representation of fl
// field.
type LeafField func(fl payment.FilterLeaf) string

// CmpString returns the string representation of c considering v (value which
// is going to be compared) and returns true if value is already contained in
// the string, otherwise false.
//
// Each filter consumer (e.g. a DB vendor) must provide a function which returns
// its representation of c considering v.
//
// Some syntaxes don't make a difference on the c based on v, but others do.
// For example in Go you can do "someiface == nil", and if `nil` is another
// value of the type someiface the equality operator is the same, however, in
// SQL need to use a different operator if you compare with NULL or specific
// value, e.g. "name IS NULL" and "name = 'some'".
type CmpString func(c payment.FilterCmp, v interface{}) (string, bool)

// LogicalString returns the string representation of l. Each filter consumer
// (e.g. a DB vendor) must provide a function which returns its representation
// of l.
type LogicalString func(l payment.FilterLogical) string

// CmpStringNoVal is a convenient function which converts a function which
// returns the string representation of a payment FilterCmp and it doesn't
// require to check the value to CmpString function.
func CmpStringNoVal(f func(c payment.FilterCmp) string) CmpString {
	return func(c payment.FilterCmp, _ interface{}) (string, bool) {
		return f(c), false
	}
}
