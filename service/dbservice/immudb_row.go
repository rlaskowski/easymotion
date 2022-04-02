package dbservice

import (
	"fmt"
	"time"

	"github.com/codenotary/immudb/embedded/sql"
)

type ImmuDBRows struct {
	rowr sql.RowReader
	row  sql.Row
	//rwmutex sync.RWMutex
}

func (i *ImmuDBRows) Next() bool {

	r, err := i.rowr.Read()

	if err == sql.ErrNoMoreRows {
		return false
	}

	i.row = *r

	return true
}

func (i *ImmuDBRows) Columns() ([]sql.ColDescriptor, error) {
	return i.rowr.Columns()
}

func (i *ImmuDBRows) Scan(params ...interface{}) error {
	if len(i.row.Values) != len(params) {
		return fmt.Errorf("different number of columns in row, expected %d got %d", len(params), len(i.row.Values))
	}

	cols, err := i.rowr.Columns()
	if err != nil {
		return fmt.Errorf("bad columns definition due to: %s", err.Error())
	}

	index := 0
	for _, c := range cols {
		key := fmt.Sprintf("(%s.%s.%s)", c.Database, c.Table, c.Column)

		valt, ok := i.row.Values[key]

		if ok {
			if err := i.parseType(params[index], valt.Value()); err != nil {
				return err
			}
		}

		index++
	}

	return nil
}

func (im *ImmuDBRows) parseType(dst, src interface{}) error {
	switch s := src.(type) {
	case string:
		switch d := dst.(type) {
		case *string:
			*d = s
			return nil
		}
	case time.Time:
		switch d := dst.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		}
	}

	return nil
}
