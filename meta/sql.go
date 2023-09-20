package meta

import "database/sql"

type SQLResults []sql.Result

func (r SQLResults) LastInsertId() (int64, error) {
	switch length := len(r); length {
	case 0:
		return 0, sql.ErrNoRows
	default:
		return r[length-1].LastInsertId()
	}
}

func (r SQLResults) RowsAffected() (int64, error) {
	var rowsAffected int64
	for _, result := range r {
		rowsi, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		rowsAffected = rowsAffected + rowsi
	}
	return rowsAffected, nil
}

func (r SQLResults) AddResult(results ...sql.Result) SQLResults {
	return append(r, results...)
}