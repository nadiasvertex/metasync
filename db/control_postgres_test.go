package control

import (
	"os"
	"testing"
)

// This test expects to find a postgres distribution in /tmp/postgres. It will create a database data cluster
// in /tmp/db.
var ctx = PostgresContext{
	"/tmp/db",
	"/tmp/postgres",
}

func TestInit(t *testing.T) {
	os.RemoveAll(ctx.data_folder)
	if err := ctx.Init(nil); err != nil {
		t.Errorf("Unable to initialize data engine with data folder='%v' because: '%v'", ctx.data_folder, err.Error())
		t.Fail()
	}
}

func TestStart(t *testing.T) {
	TestInit(t)
	if err := ctx.Start(); err != nil {
		t.Errorf("Unable to start data engine because: '%v'", err.Error())
		t.Fail()
	}

	ctx.Stop()
}

func TestStop(t *testing.T) {
	TestInit(t)

	if err := ctx.Start(); err != nil {
		t.Errorf("Unable to start data engine because: '%v'", err.Error())
		t.Fail()
	}

	if err := ctx.Stop(); err != nil {
		t.Errorf("Unable to stop data engine because: '%v'", err.Error())
		t.Fail()
	}
}
