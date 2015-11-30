package control

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
}
