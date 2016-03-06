package control

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PostgresContext struct {
	data_folder   string
	engine_folder string
	log_folder    string
	username      string
	password      string
	port          int
}

type PostgresControlError struct {
	message string
	code    int
	context *PostgresContext
}

func pg_ctl(ctx *PostgresContext, args ...string) (string, error) {

	log.Printf("running command: %v", args)

	var pg_ctl = filepath.Clean(filepath.Join(ctx.engine_folder, "bin", "pg_ctl"))
	if _, err := os.Stat(pg_ctl); os.IsNotExist(err) {
		return "", &PostgresControlError{
			fmt.Sprintf("Unable to find pg_ctl at '%v'.", pg_ctl),
			-1,
			ctx,
		}
	}

	cmd := exec.Command(pg_ctl, args...)
	raw_out, err := cmd.Output()
	out := string(raw_out)
	log.Print(out)

	if err != nil {
		return out, &PostgresControlError{
			err.Error(),
			-1,
			ctx,
		}
	}

	return out, nil
}

func pg_status(ctx *PostgresContext) status {
	output, err := pg_ctl(ctx, "status", "-D", ctx.data_folder)
	if err != nil {
		return Unknown
	}

	if strings.Contains(output, "stopped") {
		return Stopped
	} else if strings.Contains(output, "started") {
		return Running
	} else {
		return Unknown
	}
}

func logging_flag(ctx *PostgresContext) string {
	return "--log=" + filepath.Join(ctx.log_folder, "data-engine.log")
}

func (e *PostgresControlError) Error() string {
	return fmt.Sprintf("attempt to control postgres database at %v failed with code %v: '%v'", e.context.data_folder, e.code, e.message)
}

func (ctx *PostgresContext) Init(options map[string]string) error {

	if _, err := os.Stat(ctx.data_folder); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.data_folder, 0700); err != nil {
			return &PostgresControlError{
				fmt.Sprintf("Unable to create data folder '%v'.", ctx.data_folder),
				-1,
				ctx,
			}
		}
	}

	if _, err := os.Stat(ctx.log_folder); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.log_folder, 0700); err != nil {
			return &PostgresControlError{
				fmt.Sprintf("Unable to create log folder '%v'.", ctx.log_folder),
				-1,
				ctx,
			}
		}
	}

	pwfile, err := ioutil.TempFile("", "metasync")
	if err != nil {
		return err
	}
	pwfile.Write([]byte(ctx.password))
	pwfile.Close()
	defer os.Remove(pwfile.Name())

	init_options := fmt.Sprintf(`"--username=%v --pwfile=%v"`, ctx.username, pwfile.Name())

	_, err = pg_ctl(ctx, "-D", ctx.data_folder, logging_flag(ctx), "init", "-o", init_options)

	return err
}

func (ctx *PostgresContext) Start() error {
	options := fmt.Sprintf(`"-p %v"`, ctx.port)
	_, err := pg_ctl(ctx, "start", "-w", "-D", ctx.data_folder, logging_flag(ctx), "-o", options)
	return err
}

func (ctx *PostgresContext) Stop() error {
	_, err := pg_ctl(ctx, "stop", "-w", "-m", "fast", "-D", ctx.data_folder, logging_flag(ctx))
	return err
}

func (ctx *PostgresContext) Restart() error {
	if err := ctx.Stop(); err != nil {
		return err
	}

	return ctx.Start()
}

func (ctx *PostgresContext) Open(service string) (*sql.DB, error) {
	url := fmt.Sprintf(
		"postgres://%v:%v@localhost/%v?sslmode=disable&port=%v",
		ctx.username, ctx.password, service, ctx.port,
	)
	return sql.Open("postgres", url)
}

func NewPostgresContext(engine_folder string) (PostgresContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return PostgresContext{}, err
	}

	ctx := PostgresContext{
		engine_folder: engine_folder,
		data_folder:   filepath.Join(wd, "data"),
		log_folder:    filepath.Join(wd, "log"),
		username:      "metasync",
		password:      "metasyncpass",
		port:          1914,
	}

	// If the data folder doesn't exist, we need to initialize the whole system.
	if _, err = os.Stat(ctx.data_folder); err != nil && os.IsNotExist(err) {
		err = ctx.Init(nil)
		if err != nil {
			return PostgresContext{}, err
		}

	}

	return ctx, nil
}
