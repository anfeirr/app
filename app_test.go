package app_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	app.Import(&tests.Foo{})

	defer func() { recover() }()
	app.Import(tests.NoPointerCompo{})
}

func TestRunPanic(t *testing.T) {
	assert.Panics(t, func() {
		app.Run()
	})
}

func TestApp(t *testing.T) {
	app.Logger = func(format string, a ...interface{}) {
		log := fmt.Sprintf(format, a...)
		t.Log(log)
	}

	app.Import(&tests.Foo{})
	app.Import(&tests.Bar{})

	onRun := func() {
		d := app.CurrentDriver()
		require.NotNil(t, d)

		assert.NotEmpty(t, app.Name())
		assert.Equal(t, filepath.Join("resources", "hello", "world"), app.Resources("hello", "world"))
		assert.Equal(t, filepath.Join("storage", "hello", "world"), app.Storage("hello", "world"))

		assert.NotNil(t, app.New(tests.UnsupportedConfig{}))
		assert.NotNil(t, app.New(app.WindowConfig{}))
		assert.NotNil(t, app.ElemByCompo(&tests.Hello{}))
		app.Render(&tests.Hello{})

		app.Emit("test")

		app.UI(func() {
			app.Logf("hello")
			app.Close()
		})
	}

	defer app.NewSubscriber().
		Subscribe(app.Running, onRun).
		Close()

	app.Run(&tests.Driver{
		SimulatedTarget: "web",
	})
}

func TestLog(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	app.Log("hello", "world")
	assert.Equal(t, "hello world", log)

	app.Logf("%s %s", "bye", "world")
	assert.Equal(t, "bye world", log)
}

func TestPanic(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "hello world", log)
		assert.Equal(t, "hello world", err)
	}()

	app.Panic("hello", "world")
	assert.Fail(t, "no panic")
}

func TestPanicf(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "bye world", log)
		assert.Equal(t, "bye world", err)
	}()

	app.Panicf("%s %s", "bye", "world")
	assert.Fail(t, "no panic")
}

func TestPretty(t *testing.T) {
	t.Log(app.Pretty(app.WindowConfig{}))
}
