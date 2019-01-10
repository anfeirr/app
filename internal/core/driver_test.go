package core_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestDriver(t *testing.T) {
	tests.TestDriver(t, func() app.Driver {
		return &tests.Driver{
			SimulatedTarget: "test",
		}
	})
}
