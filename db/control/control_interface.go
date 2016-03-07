package control

import "database/sql"

type status int

const (
	Starting status = iota
	Running
	Stopping
	Stopped
	Unknown
)

type DatabaseControl interface {
	Init(options map[string]string) error
	Start() error
	Stop() error
	Restart() error
	Status() status
	Open() (*sql.DB, error)
}

type DatabaseInit struct {
	Article func(*sql.DB)
	Model   func(*sql.DB)
}
