package models

type Exists struct {
	Count int8 `db:"count"`
}

func (e Exists) Found() bool {
	return e.Count > 0
}
