package sqrlx

import (
	"database/sql"
	"fmt"
)

// IRows is the interface of *sql.Rows
type IRows interface {
	Scan(...interface{}) error
	Columns() ([]string, error)
	Next() bool
	Close() error
	Err() error
}

var _ IRows = &sql.Rows{}

type Rows struct {
	IRows
}

type Row struct {
	Rows IRows
	err  error
}

func (r Row) Scan(into ...interface{}) error {
	if r.err != nil {
		return fmt.Errorf("row held error %w", r.err)
	}
	if !r.Rows.Next() {
		return sql.ErrNoRows
	}
	defer r.Rows.Close()

	return r.Rows.Scan(into...)
}

func (r Row) ScanStruct(into interface{}) error {
	if err := ScanStruct(r, into); err != nil {
		return fmt.Errorf("scan struct: %w", err)
	}
	return nil
}

func (r Row) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.Rows.Columns()
}
