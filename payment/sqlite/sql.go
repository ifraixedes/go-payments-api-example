package sqlite

import (
	"strings"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"go.fraixed.es/errors"
)

// dbScanPymt is the function which is used to scan the list of columns of a row
// a select statement for payments.
type dbScanPymt func(*sqlite3.Stmt) (payment.Pymt, error)

// selectPymtColumns is a convenient function which returns a string which the
// list of columns to be use in a payments select statement and the dbScanPymt
// function based on s.
func selectPymtColumns(s payment.Selection) (string, dbScanPymt) {
	var sf = make([]string, 1, 5)

	sf[0] = "id"

	if s.Version {
		sf = append(sf, "version")
	}

	if s.OrgID {
		sf = append(sf, "organisation_id")
	}

	if s.Type {
		sf = append(sf, "json_extract(data, '$.type') as type")
	}

	if s.Attributes {
		sf = append(sf, "data")
	}

	return strings.Join(sf, ", "), func(stmt *sqlite3.Stmt) (payment.Pymt, error) {
		return dbScanPymtFromSelection(s, stmt)
	}
}

//  dbScanPymtFromSelection scans the columns of a row of payments select
// statement indicated by sl and using stmt.
func dbScanPymtFromSelection(sl payment.Selection, stmt *sqlite3.Stmt) (payment.Pymt, error) {
	var (
		p    payment.Pymt
		cidx = 0
	)

	s, _, err := stmt.ColumnText(cidx)
	if err != nil {
		return p, errors.Wrap(
			err, payment.ErrUnexpectedStoreError, payment.ErrMDFnCall("sqlite3.Stmt.ColumnText", cidx),
		)
	}

	cidx++
	p.ID, err = uuid.FromString(s)
	if err != nil {
		return p, errors.Wrap(err, ErrInvalidFormatID, payment.ErrMDVar("id", s))
	}

	if sl.Version {
		i, _, err := stmt.ColumnInt64(cidx)
		if err != nil {
			return p, errors.Wrap(
				err, payment.ErrUnexpectedStoreError, payment.ErrMDFnCall("sqlite3.Stmt.ColumnInt64", cidx),
			)
		}

		cidx++
		p.Version = uint32(i)
	}

	if sl.OrgID {
		s, _, err := stmt.ColumnText(cidx)
		if err != nil {
			return p, errors.Wrap(
				err, payment.ErrUnexpectedStoreError, payment.ErrMDFnCall("sqlite3.Stmt.ColumnText", cidx),
			)
		}

		cidx++
		p.OrgID, err = uuid.FromString(s)
		if err != nil {
			return p, errors.Wrap(err, ErrInvalidFormatID, payment.ErrMDVar("organisation_id", s))
		}
	}

	if sl.Type {
		s, _, err := stmt.ColumnText(cidx)
		if err != nil {
			return p, errors.Wrap(
				err, payment.ErrUnexpectedStoreError, payment.ErrMDFnCall("sqlite3.Stmt.ColumnText", cidx),
			)
		}

		cidx++
		p.Type = s
	}

	if sl.Attributes {
		b, err := stmt.ColumnBlob(cidx)
		if err != nil {
			return p, errors.Wrap(
				err, payment.ErrUnexpectedStoreError, payment.ErrMDFnCall("sqlite3.Stmt.ColumnBlob", cidx),
			)
		}

		var pd pymtData
		if err = pd.Deserialize(b); err != nil {
			return p, err
		}

		p.Attributes = pd.Attrs
		p.Type = pd.Type
	}

	return p, nil
}

func leafField(fl payment.FilterLeaf) string {
	switch fl.(type) {
	case payment.FilterLeafAmount:
		return "json_extract(data, '$.amount')"
	case payment.FilterLeafType:
		return "json_extract(data, '$.type')"
	case payment.FilterLeafID:
		return "id"
	}

	// This happens is that new filters have been added and this function has not
	// been updated
	return ""
}

func orderByColumns(s payment.Sort) string {
	var cols []string

	if s.ID.Valid() {
		cols = append(cols, "id "+orderDir(s.ID))
	}

	if s.Type.Valid() {
		cols = append(cols, "json_extract(data, '$.type') "+orderDir(s.Type))
	}

	if s.Version.Valid() {
		cols = append(cols, "version "+orderDir(s.Version))
	}

	if s.OrgID.Valid() {
		cols = append(cols, "organisation_id "+orderDir(s.OrgID))
	}

	if s.Attributes.Amount.Valid() {
		cols = append(cols, "json_extract(data, '$.amount') "+orderDir(s.Attributes.Amount))
	}

	return strings.Join(cols, ",")
}

func orderDir(s payment.SortDir) string {
	switch s {
	case payment.SortAscending:
		return "ASC"
	case payment.SortDescending:
		return "DESC"
	}

	return ""
}

func limitOffset(c payment.Chunk) []interface{} {
	if c.Limit == 0 {
		return nil
	}

	return []interface{}{c.Limit, c.Offset}
}

func adaptArgsToSQL(args []interface{}) []interface{} {
	for i, a := range args {
		switch v := a.(type) {
		case uuid.UUID:
			args[i] = v.String()
		case uint32:
			args[i] = int64(v)
		case uint64:
			args[i] = int64(v)
		}
	}

	return args
}
