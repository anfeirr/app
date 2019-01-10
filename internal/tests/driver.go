package tests

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver provide a app.Driver implementation suitable for testing.
type Driver struct {
	SimulatedTarget string

	ui      chan func()
	factory *app.Factory
	events  *app.EventRegistry
	elems   *core.ElemDB
	close   func()
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return d.SimulatedTarget
}

// Call satisfies the app.Driver interface.
func (d *Driver) Call(method string, out interface{}, in interface{}) error {
	return nil
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) {
	d.ui = c.UI
	d.factory = c.Factory
	d.events = c.Events
	d.elems = core.NewElemDB()

	ctx, cancel := context.WithCancel(context.Background())
	d.close = cancel

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case fn := <-d.ui:
				fn()

			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	d.events.Emit(app.Running)
	wg.Wait()
}

// Factory satisfies the app.Driver interface.
func (d *Driver) Factory() *app.Factory {
	return d.factory
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "test"
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	p := filepath.Join(path...)
	return filepath.Join("resources", p)
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(path ...string) string {
	p := filepath.Join(path...)
	return filepath.Join("storage", p)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)

	if e.Err() == app.ErrElemNotSet {
		return
	}

	e.(app.View).Render(c)
}

// New satisfies the app.Driver interface.
func (d *Driver) New(c app.ElemConfig) app.Elem {
	switch c := c.(type) {
	case app.WindowConfig:
		return newWindow(d, c)

	default:
		e := &core.Elem{}
		e.SetErr(app.ErrNotSupported)
		return e
	}
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.elems.GetByCompo(c)
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return nil
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	return nil
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.ui <- f
}

// Close satisfies the app.Driver interface.
func (d *Driver) Close() {
	d.close()
}
