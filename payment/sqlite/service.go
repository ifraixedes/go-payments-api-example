package sqlite

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/ifraixedes/go-payments-api-example/payment/internal/filter"
	"go.fraixed.es/errors"
)

// New creates an instance of the SQLite implementation of the payment Service.
//
// fname is any of the values that SQLite filename can take (see
// https://www.sqlite.org/c3ref/open.html, for more information), which are:
//
// * Path to a file, which is created if it doesn't exists
//
// * An URI as described in https://www.sqlite.org/uri.html
//
// * The string ":memory:", which creates an in-memory database
//
// When not using the URI fname, the SQLite database is always opened for
// read/write operations and with shared cache mode enabled (see
// https://www.sqlite.org/sharedcache.html for more information), otherwise URI
// parameters must specify the connection and cache mode through the query
// parameters.
// All the connections are opened with WAL method (see
// https://www.sqlite.org/wal.html for more information).
//
// The following error codes can be returned:
//
// * ErrInvalidArgDBFname
//
// * ErrDBCantOpen
//
// * payment.ErrUnexpectedStoreError
//
// * payment.ErrUnexpectedOSError - this error happens if there is an error when
//   resolving the absolute path of the fname is a path to a file.
func New(fname string) (payment.Service, error) {
	if fname == "" {
		return nil, errors.New(ErrInvalidArgDBFname, payment.ErrMDArg("fname", fname))
	}

	var (
		svc = service{
			fname: fname,
		}
		isURI bool
	)

	switch {
	case fname == ":memory:":
		// When in-memory, the URI format must be used for being able to enabled the
		// shared cache
		svc.fname = "file::memory:?cache=shared"
	case !strings.HasPrefix(fname, "file:"):
		var err error
		svc.fname, err = filepath.Abs(fname)
		if err != nil {
			return nil, errors.Wrap(err, payment.ErrUnexpectedOSError, payment.ErrMDArg("fname", fname))
		}
		svc.openFlags = sqlite3.OPEN_READWRITE | sqlite3.OPEN_CREATE | sqlite3.OPEN_SHAREDCACHE
	default:
		isURI = true
	}

	var c, pc, err = svc.openConn(context.Background())
	if err != nil {
		if pc == sqlite3.ERROR && isURI {
			// Invalid URI format returns this error code
			return nil, errors.Wrap(err, ErrInvalidArgDBFname, payment.ErrMDArg("fname", fname))
		}

		return nil, err
	}

	if err := c.Close(); err != nil {
		return nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	return &svc, nil
}

type service struct {
	fname     string
	openFlags int
}

// Create stores p in the database.
//
// The function will return all the errors that payment.Service documents plus
// the following ones:
//
// * ErrDBCantOpen
//
// * ErrDBLimit
//
// * ErrDBSchemaChanged
//
// * ErrInvalidPayment
func (s *service) Create(ctx context.Context, p payment.PymtUpsert) (uuid.UUID, error) {
	if err := p.Validate(); err != nil {
		return uuid.Nil, err
	}

	var id, err = uuid.NewV4()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	var pd []byte
	{
		d := &pymtData{}
		d.Init(p)
		pd, err = d.Serialize()
		if err != nil {
			return uuid.Nil, err
		}
	}

	conn, _, err := s.openConn(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	err = conn.Exec(
		"INSERT INTO payments(id, organisation_id, data) VALUES (?, ?, ?)",
		id.String(), p.OrgID.String(), pd,
	)
	if err != nil {
		if cerr := handleSQLiteErrCommon(err); cerr != nil {
			return uuid.Nil, cerr
		}

		var pc, _, serr = isSQLiteErr(err)
		if serr != nil {
			if pc == sqlite3.CONSTRAINT {
				return uuid.Nil, errors.Wrap(serr, ErrInvalidPayment)
			}
		}

		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	return id, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}

	conn, _, err := s.openConn(ctx)
	if err != nil {
		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	err = conn.Exec("DELETE FROM payments WHERE id = ?", id.String())
	if err != nil {
		return handleSQLiteErr(err)
	}

	if conn.TotalChanges() == 0 {
		return errors.New(payment.ErrNotFound, payment.ErrMDVar("id", id))
	}

	return nil
}

func (s *service) Find(
	ctx context.Context, pf payment.Filter, sl payment.Selection, st payment.Sort, pc payment.Chunk,
) ([]payment.Pymt, error) {
	var (
		stmtargs         []interface{}
		ordby            = orderByColumns(st)
		sel, scanPymt    = selectPymtColumns(sl)
		limitargs        = limitOffset(pc)
		where, whereargs = filter.SQL(pf, leafField)
		//nolint:gosec
		query = fmt.Sprintf("SELECT %s FROM payments", sel)
	)

	if where != "" {
		//nolint:gosec
		query = fmt.Sprintf("%s WHERE %s", query, where)
		stmtargs = whereargs
	}

	if ordby != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, ordby)
	}

	if len(limitargs) > 0 {
		query = fmt.Sprintf("%s LIMIT ? OFFSET ?", query)
		stmtargs = append(stmtargs, limitargs...)
	}

	stmtargs = adaptArgsToSQL(stmtargs)

	var conn, _, err = s.openConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	stmt, err := conn.Prepare(query, stmtargs...)
	if err != nil {
		return nil, handleSQLiteErr(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var plist []payment.Pymt
	for {
		ok, err := stmt.Step()
		if err != nil {
			return nil, handleSQLiteErr(err)
		}
		if !ok {
			break
		}

		p, err := scanPymt(stmt)
		if err != nil {
			return nil, err
		}

		plist = append(plist, p)
	}

	return plist, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID, sl payment.Selection) (payment.Pymt, error) {
	if id == uuid.Nil {
		return payment.Pymt{}, errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}

	var conn, _, err = s.openConn(ctx)
	if err != nil {
		return payment.Pymt{}, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	var sq, scanPymt = selectPymtColumns(sl)
	//nolint:gosec
	stmt, err := conn.Prepare(fmt.Sprintf("SELECT %s FROM payments WHERE id = ?", sq), id.String())
	if err != nil {
		return payment.Pymt{}, handleSQLiteErr(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	ok, err := stmt.Step()
	if err != nil {
		return payment.Pymt{}, handleSQLiteErr(err)
	}
	if !ok {
		return payment.Pymt{}, errors.New(payment.ErrNotFound, payment.ErrMDVar("id", id))
	}

	return scanPymt(stmt)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, ver uint32, p payment.PymtUpsert) error {
	if id == uuid.Nil {
		return errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}

	if err := p.Validate(); err != nil {
		return err
	}

	var pd []byte
	{
		var (
			err error
			d   = &pymtData{}
		)
		d.Init(p)
		pd, err = d.Serialize()
		if err != nil {
			return err
		}
	}

	var conn, _, err = s.openConn(ctx)
	if err != nil {
		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	// The error is map to a new var because there is a bug in WithTx functions
	// which hides the error when rollback or commit succeeds, so we use the err
	// var of the parent scope to get the proper error and only errtx for the
	// rollback or commit errors
	// See https://github.com/bvinc/go-sqlite-lite/pull/20
	var errtx = conn.WithTx(func() error {
		err = conn.Exec(
			"UPDATE payments SET version = version + 1, organisation_id = ?, data = ? WHERE id = ? AND version = ?",
			p.OrgID.String(), pd, id.String(), int64(ver),
		)
		if err != nil {
			if cerr := handleSQLiteErrCommon(err); cerr != nil {
				err = cerr
				return cerr
			}

			var pc, _, serr = isSQLiteErr(err)
			if serr != nil {
				if pc == sqlite3.CONSTRAINT {
					err = errors.Wrap(serr, ErrInvalidPayment)
					return err
				}
			}

			err = errors.Wrap(err, payment.ErrUnexpectedStoreError)
			return err
		}

		if conn.TotalChanges() != 1 {
			var stmt *sqlite3.Stmt
			stmt, err = conn.Prepare("SELECT version FROM payments WHERE id = ?", id.String())
			if err != nil {
				err = handleSQLiteErr(err)
				return err
			}

			defer func() {
				_ = stmt.Close()
			}()

			var ok bool
			ok, err = stmt.Step()
			if err != nil {
				err = handleSQLiteErr(err)
				return err
			}
			if !ok {
				err = errors.New(payment.ErrNotFound, payment.ErrMDVar("id", id))
				return err
			}

			var v int64
			v, _, err = stmt.ColumnInt64(0)
			if err != nil {
				err = handleSQLiteErr(err)
				return err
			}

			if uint32(v) != ver {
				err = errors.New(payment.ErrInvalidArgVersionMismatch,
					payment.ErrMDArg("version", ver), payment.ErrMDFact("current_version", v),
				)

				return err
			}

			// This should never happen, but if it happens then return this error,than
			// returning silently
			err = errors.New(payment.ErrUnexpectedStoreError)
			return err
		}

		return nil
	})

	if err == nil && errtx != nil {
		return errors.Wrap(errtx, payment.ErrUnexpectedStoreError)
	}

	return err
}

// openConn create a new sqlite3 connection.
// It returns error the connection creation fails or the WAL journal model cannot
// be set. When an error is returned, the sqlite3 error primary code is also
// returned (see https://www.sqlite.org/rescode.html).
//
// If ctx has a deadline, it's used to set the BusyTimeout to the connection,
// see https://godoc.org/github.com/bvinc/go-sqlite-lite/sqlite3#Conn.BusyFunc
//
// The following error codes can be returned:
//
// * ErrDBCantOpen
//
// * payment.ErrUnexpectedStoreError
func (s *service) openConn(ctx context.Context) (*sqlite3.Conn, uint8, error) {
	var (
		err  error
		conn *sqlite3.Conn
	)
	if s.openFlags == 0 {
		conn, err = sqlite3.Open(s.fname)
	} else {
		conn, err = sqlite3.Open(s.fname, s.openFlags)
	}

	if err != nil {
		var pc, _, serr = isSQLiteErr(err)
		if serr == nil {
			return nil, 0, errors.Wrap(err, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", s.fname),
				payment.ErrMDVar("sqlite_open_flags", s.openFlags),
			)
		}

		switch pc {
		case sqlite3.CANTOPEN, sqlite3.NOTADB:
			return nil, pc, errors.Wrap(serr, ErrDBCantOpen, payment.ErrMDArg("fname", s.fname))
		default:
			return nil, pc, errors.Wrap(serr, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", s.fname),
				payment.ErrMDVar("sqlite_open_flags", s.openFlags),
			)
		}
	}

	if err = conn.Exec("PRAGMA journal_mode=wal"); err != nil {
		pc, _, _ := isSQLiteErr(err)

		return nil, pc, errors.Wrap(err, payment.ErrUnexpectedStoreError,
			payment.ErrMDFnCall("sqlite3.Conn.Exec", "PRAGMA journal_mode=wal"),
		)
	}

	if t, ok := ctx.Deadline(); ok {
		conn.BusyTimeout(time.Until(t))
	}

	return conn, 0, nil
}

// isSQLiteErr returns the primary error code, the extended error code and the
// specific sqlite3 error when it's of such type or 0, 0 and nil when not.
// See https://www.sqlite.org/rescode.html
func isSQLiteErr(err error) (uint8, int, *sqlite3.Error) {
	if serr, ok := err.(*sqlite3.Error); ok {
		c := serr.Code()
		return uint8(c & 0xFF), c, serr
	}

	return 0, 0, nil
}

// handleSQLiteErrCommon is a convenient function to map SQLite error codes
// which are common to the most of the SQL statements executed by the service
// methods. It returns an error is the err is successful mapped otherwise nil,
// which means that the caller should deal with it because err isn't a common
// SQLite error.
func handleSQLiteErrCommon(err error) error {
	var pc, _, serr = isSQLiteErr(err)
	if serr == nil {
		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	switch pc {
	case sqlite3.BUSY, sqlite3.INTERRUPT, sqlite3.LOCKED:
		return errors.Wrap(serr, payment.ErrAbortedOperation)
	case sqlite3.SCHEMA:
		return errors.Wrap(serr, ErrDBSchemaChanged)
	case sqlite3.CORRUPT, sqlite3.IOERR, sqlite3.FULL, sqlite3.INTERNAL, sqlite3.PROTOCOL:
		return errors.Wrap(serr, payment.ErrUnexpectedStoreError)
	case sqlite3.NOMEM:
		return errors.Wrap(serr, payment.ErrUnexpectedSysError)
	case sqlite3.TOOBIG:
		return errors.Wrap(err, ErrDBLimit)
	}

	return nil
}

// handleSQLiteErr is a convenient function to handle the error returned by the
// SQLite operations when the caller cannot make any differentiation when
// handleSQLiteErrCommon returns nil than reutrning
// payment.ErrUnexpectedStoreError error code.
func handleSQLiteErr(err error) error {
	if serr := handleSQLiteErrCommon(err); serr != nil {
		return err
	}

	return errors.Wrap(err, payment.ErrUnexpectedStoreError)
}
