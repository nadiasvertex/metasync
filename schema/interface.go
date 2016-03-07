package schema

import "database/sql"

type Getter interface {
	Get(uuid string, db *sql.DB) error
}

type Putter interface {
	Put(db *sql.DB) error
}
