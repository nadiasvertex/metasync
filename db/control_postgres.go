package control

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type PostgresContext struct {
	data_folder   string
	engine_folder string
}

type PostgresControlError struct {
	message string
	code    int
	context *PostgresContext
}

func pg_ctl(ctx *PostgresContext, args ...string) (string, error) {

	log.Printf("running command: %v", args)

	var pg_ctl = path.Clean(path.Join(ctx.engine_folder, "bin", "pg_ctl"))
	if _, err := os.Stat(pg_ctl); os.IsNotExist(err) {
		return "", &PostgresControlError{
			fmt.Sprintf("Unable to find pg_ctl at '%v'.", pg_ctl),
			-1,
			ctx,
		}
	}

	cmd := exec.Command(pg_ctl, args...)
	raw_out, err := cmd.CombinedOutput()
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

func (e *PostgresControlError) Error() string {
	return fmt.Sprintf("attempt to control postgres database at %v failed with code %v '%v'", e.context.data_folder, e.code, e.message)
}

func (ctx *PostgresContext) Init(options map[string]string) error {

	if _, err := os.Stat(ctx.data_folder); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.data_folder, 0600); err != nil {
			return &PostgresControlError{
				fmt.Sprintf("Unable to create data folder '%v'.", ctx.data_folder),
				-1,
				ctx,
			}
		}
	}
	_, err := pg_ctl(ctx, "-D", ctx.data_folder, "init")
	return err
}

func (ctx *PostgresContext) Start() error {
	_, err := pg_ctl(ctx, "-w", "-D", ctx.data_folder, "start")
	return err
}

func (ctx *PostgresContext) Stop() error {
	_, err := pg_ctl(ctx, "-w", "-m", "fast", "-D", ctx.data_folder, "stop")
	return err
}

func (ctx *PostgresContext) Restart() error {
	if err := ctx.Stop(); err != nil {
		return err
	}

	return ctx.Start()
}
