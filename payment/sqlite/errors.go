package sqlite

type code uint8

// The list of specific error codes that the sqlite service implementation can
// return.
//
// The ErrInvalidFormat errors are related with data stored in the DB which
// isn't of the expected format; these can only happen if the data has been
// inserted by other application and it has inserted/updated the values without
// previous verification. Another cases could be that the application has been
// changed and the data in the DB hasn't been, but then, they would be related
// to a bug in the process of changing the application.
const (
	ErrDBCantOpen code = iota + 1
	ErrDBLimit

	ErrDBSchemaChanged

	ErrInvalidArgDBFname

	ErrInvalidFormatBlob
	ErrInvalidFormatID

	ErrInvalidPayment
)

func (c code) String() string {
	switch c {
	case ErrDBCantOpen:
		return "DBCantOpen"
	case ErrDBLimit:
		return "DBLimit"
	case ErrDBSchemaChanged:
		return "DBSchemaChanged"
	case ErrInvalidArgDBFname:
		return "InvalidArgDBFname"
	case ErrInvalidFormatBlob:
		return "InvalidFormatBlob"
	case ErrInvalidFormatID:
		return "InvalidFormatID"
	case ErrInvalidPayment:
		return "InvalidPayment"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrDBCantOpen:
		return "the SQLite DB cannot be open. Causes: DB is corrupted, file " +
			"isn't a valid DB file or some DB files cannot be open"
	case ErrDBLimit:
		return "the operation failed because a limitation of SQLite DB"
	case ErrDBSchemaChanged:
		return "the schema of the DB has been altered, this may happen because a " +
			"new version of the service has been released using the same DB, in any " +
			"case the advice is to close the service and in case of a new version, " +
			"run it again with the new version"
	case ErrInvalidArgDBFname:
		return "the SQLite filename isn't of a valid format"
	case ErrInvalidFormatBlob:
		return "the blob stored in the DB isn't of a valid format"
	case ErrInvalidFormatID:
		return "the ID stored in the DB isn't of a valid format"
	case ErrInvalidPayment:
		return "the payment is valid due the constrains imposed by the DB schema"
	}

	return ""
}
