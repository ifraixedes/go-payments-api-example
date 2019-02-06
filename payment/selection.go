package payment

// Selection specifies the fields of a single payment to be retrieved.
// The payment's fields which aren't present are always retrieved.
// Each file is a boolean, when it's true, the value is retrieved otherwise it
// won't be.
type Selection struct {
	Version    bool
	Type       bool
	OrgID      bool
	Attributes bool
}
