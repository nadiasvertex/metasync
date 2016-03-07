package schema

import (
	"github.com/nadiasvertex/metasync/db/control"
	"testing"
)

func TestLoadArticle(t *testing.T) {
	ctx, err := control.NewPostgresContext("/tmp/postgres")
	if err != nil {
		t.Errorf("Unable to connect to database: %s", err)
	}

	c, err := ctx.Open("metasync")
	if err != nil {
		t.Errorf("Unable to connect to database: %s", err)
	}

	var a *Article = nil
	a.Get("abc-def", c)
}
