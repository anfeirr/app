package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

// DriverSetup is the definition of a function that creates a driver.
type DriverSetup func() app.Driver

// TestDriver is a test suite that ensure that all driver implementations behave
// the same.
func TestDriver(t *testing.T, setup DriverSetup) {
	ui := make(chan func(), 32)

	factory := app.NewFactory()
	factory.RegisterCompo(&Hello{})
	factory.RegisterCompo(&World{})
	factory.RegisterCompo(&Menu{})

	events := app.NewEventRegistry(ui)
	sub := &app.Subscriber{Events: events}

	driver := setup()

	defer sub.Subscribe(app.Running, func() {
		assert.NotEmpty(t, driver.AppName())
		assert.NotEmpty(t, driver.Resources())
		assert.NotEmpty(t, driver.Storage())

		driver.New(app.WindowConfig{}).
			WhenWindow(func(w app.Window) {
				testWindow(t, w)
			})

		driver.UI(func() {
			driver.Close()
		})
	}).Close()

	driver.Run(app.DriverConfig{
		UI:      ui,
		Factory: factory,
		Events:  events,
	})
}

func testWindow(t *testing.T, w app.Window) {
	assertElem(t, w)
	assert.NotEmpty(t, w.ID())

	isMenu := false
	w.WhenMenu(func(app.Menu) {
		isMenu = true
	})
	assert.False(t, isMenu)

	isDockTile := false
	w.WhenDockTile(func(app.DockTile) {
		isDockTile = true
	})
	assert.False(t, isDockTile)

	isStatusMenu := false
	w.WhenStatusMenu(func(app.StatusMenu) {
		isStatusMenu = true
	})
	assert.False(t, isStatusMenu)

	w.Reload()
	assert.Error(t, w.Err())

	w.Load("tests.hello")
	assertElem(t, w)
	assert.NotNil(t, w.Compo())

	w.Render(w.Compo())
	assertElem(t, w)

	w.Reload()
	assertElem(t, w)

	assert.False(t, w.CanPrevious())
	assert.False(t, w.CanNext())

	w.Previous()
	assert.Error(t, w.Err())

	w.Load("tests.world")
	assertElem(t, w)

	w.Next()
	assert.Error(t, w.Err())

	w.Load("tests.unknown")
	assert.Error(t, w.Err())

	assert.NotNil(t, w.Compo())
	assert.True(t, w.Contains(w.Compo()))

	w.EvalJS(nil, "alert(%s)", "test window")
	assertElem(t, w)

	w.Position()
	assertElem(t, w)

	w.Move(42, 42)
	assertElem(t, w)

	w.Center()
	assertElem(t, w)

	w.Size()
	assertElem(t, w)

	w.Resize(42, 42)
	assertElem(t, w)

	w.Focus()
	assertElem(t, w)

	w.FullScreen()
	assertElem(t, w)

	w.ExitFullScreen()
	assertElem(t, w)

	w.Minimize()
	assertElem(t, w)

	w.Deminimize()
	assertElem(t, w)

	w.Close()
}

func assertElem(t *testing.T, e app.Elem) {
	if e.Err() == app.ErrNotSupported {
		return
	}

	assert.NoError(t, e.Err())
}
