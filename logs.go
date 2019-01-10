package app

import (
	"encoding/json"
	"fmt"
)

// DriverWithLogs decorates a driver with logs.
func DriverWithLogs(d Driver) Driver {
	return &driverWithLogs{Driver: d}
}

// Driver logs.
type driverWithLogs struct {
	Driver
}

func (d *driverWithLogs) Run(c DriverConfig) {
	Logf("running %T driver", d.Driver)
	d.Driver.Run(c)
}

func (d *driverWithLogs) Render(c Compo) {
	e := d.ElemByCompo(c)

	if view, ok := e.(View); ok {
		view.Render(c)
	}
}

func (d *driverWithLogs) New(c ElemConfig) Elem {
	WhenDebug(func() {
		Logf("creating element from %T: %s", c, c.Dump())
	})

	e := d.Driver.New(c)
	if e.Err() != nil {
		Logf("creating element failed: %s", e.Err())
	}

	return d.elemLogger(e)
}

func (d *driverWithLogs) ElemByCompo(c Compo) Elem {
	e := d.Driver.ElemByCompo(c)
	return d.elemLogger(e)
}

func (d *driverWithLogs) elemLogger(e Elem) Elem {
	switch v := e.(type) {
	case Window:
		return &windowWithLogs{Window: v}

	case DockTile:
		return &dockWithLogs{DockTile: v}

	case StatusMenu:
		return &statusMenuWithLogs{StatusMenu: v}

	case Menu:
		return &menuWithLogs{Menu: v}

	default:
		return v
	}
}

func (d *driverWithLogs) MenuBar() Menu {
	WhenDebug(func() {
		Logf("getting menubar")
	})

	m := d.Driver.MenuBar()
	if m.Err() != nil {
		Logf("getting menubar failed: %s", m.Err())
	}

	return &menuWithLogs{Menu: m}
}

func (d *driverWithLogs) DockTile() DockTile {
	WhenDebug(func() {
		Logf("getting dock tile")
	})

	dt := d.Driver.DockTile()
	if dt.Err() != nil {
		Logf("getting dock tile failed: %s", dt.Err())
	}

	return &dockWithLogs{DockTile: dt}
}

func (d *driverWithLogs) Close() {
	WhenDebug(func() {
		Logf("closing driver")
	})

	d.Driver.Close()
}

// Window logs.
type windowWithLogs struct {
	Window
}

func (w *windowWithLogs) WhenWindow(f func(Window)) {
	f(w)
}

func (w *windowWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("window %s is loading %s",
			w.ID(),
			parsedURL,
		)
	})

	w.Window.Load(url, v...)
	if w.Err() != nil {
		Logf("window %s failed to load %s: %s",
			w.ID(),
			parsedURL,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Reload() {
	WhenDebug(func() {
		Logf("window %s is reloading", w.ID())
	})

	w.Window.Reload()
	if w.Err() != nil {
		Logf("window %s failed to reload: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Previous() {
	WhenDebug(func() {
		Logf("window %s is loading previous", w.ID())
	})

	w.Window.Previous()
	if w.Err() != nil {
		Logf("window %s failed to load previous: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Next() {
	WhenDebug(func() {
		Logf("window %s is loading next", w.ID())
	})

	w.Window.Next()
	if w.Err() != nil {
		Logf("window %s failed to load next: %s",
			w.ID(),
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("window %s is rendering %T",
			w.ID(),
			c,
		)
	})

	w.Window.Render(c)
	if w.Err() != nil {
		Logf("window %s failed to render %T: %s",
			w.ID(),
			c,
			w.Err(),
		)
	}
}

func (w *windowWithLogs) Move(x, y float64) {
	WhenDebug(func() {
		Logf("window %s is moving to x:%.2f y:%.2f",
			w.ID(),
			x,
			y,
		)
	})

	w.Window.Move(x, y)
}

func (w *windowWithLogs) Center() {
	WhenDebug(func() {
		Logf("window %s is moving to center", w.ID())
	})

	w.Window.Center()
}

func (w *windowWithLogs) Resize(width, height float64) {
	WhenDebug(func() {
		Logf("window %s is resizing to width:%.2f height:%.2f",
			w.ID(),
			width,
			height,
		)
	})

	w.Window.Resize(width, height)
}

func (w *windowWithLogs) Focus() {
	WhenDebug(func() {
		Logf("window %s is getting focus", w.ID())
	})

	w.Window.Focus()
}

func (w *windowWithLogs) FullScreen() {
	WhenDebug(func() {
		Logf("window %s is entering full screen", w.ID())
	})

	w.Window.FullScreen()
}

func (w *windowWithLogs) ExitFullScreen() {
	WhenDebug(func() {
		Logf("window %s is exiting full screen", w.ID())
	})

	w.Window.ExitFullScreen()
}

func (w *windowWithLogs) Minimize() {
	WhenDebug(func() {
		Logf("window %s is minimizing", w.ID())
	})

	w.Window.Minimize()
}

func (w *windowWithLogs) Deminimize() {
	WhenDebug(func() {
		Logf("window %s is deminimizing", w.ID())
	})

	w.Window.Deminimize()
}

func (w *windowWithLogs) Close() {
	WhenDebug(func() {
		Logf("window %s is closing", w.ID())
	})

	w.Window.Close()
	if w.Err() != nil {
		Logf("window %s failed to close: %s",
			w.ID(),
			w.Err(),
		)
	}
}

// Controller logs.
type controllerWithLogs struct {
	Controller
}

func (c *controllerWithLogs) Close() {
	WhenDebug(func() {
		Logf("controller %s is closing", c.ID())
	})

	c.Controller.Close()
	if c.Err() != nil {
		Logf("controller %s failed to close: %s",
			c.ID(),
			c.Err(),
		)
	}
}

// Menu logs.
type menuWithLogs struct {
	Menu
}

func (m *menuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("%s %s is loading %s",
			m.Type(),
			m.ID(),
			parsedURL,
		)
	})

	m.Menu.Load(url, v...)
	if m.Err() != nil {
		Logf("%s %s failed to load %s: %s",
			m.Type(),
			m.ID(),
			parsedURL,
			m.Err(),
		)
	}
}

func (m *menuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("%s %s is rendering %T",
			m.Type(),
			m.ID(),
			c,
		)
	})

	m.Menu.Render(c)
	if m.Err() != nil {
		Logf("%s %s failed to render %T: %s",
			m.Type(),
			m.ID(),
			c,
			m.Err(),
		)
	}
}

// Dock tile logs.
type dockWithLogs struct {
	DockTile
}

func (d *dockWithLogs) WhenDockTile(f func(DockTile)) {
	f(d)
}

func (d *dockWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("dock tile is loading %s", parsedURL)
	})

	d.DockTile.Load(url, v...)
	if d.Err() != nil {
		Logf("dock tile failed to load %s: %s",
			parsedURL,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("dock tile is rendering %T", c)
	})

	d.DockTile.Render(c)
	if d.Err() != nil {
		Logf("dock tile failed to render %T: %s",
			c,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) SetIcon(name string) {
	WhenDebug(func() {
		Logf("dock tile is setting its icon to %s", name)
	})

	d.DockTile.SetIcon(name)
	if d.Err() != nil {
		Logf("dock tile failed to set its icon: %s", d.Err())
	}
}

func (d *dockWithLogs) SetBadge(v interface{}) {
	WhenDebug(func() {
		Logf("dock tile is setting its badge to %v", v)
	})

	d.DockTile.SetBadge(v)
	if d.Err() != nil {
		Logf("dock tile failed to set its badge: %s", d.Err())
	}
}

// Status menu logs.
type statusMenuWithLogs struct {
	StatusMenu
}

func (s *statusMenuWithLogs) WhenStatusMenu(f func(StatusMenu)) {
	f(s)
}

func (s *statusMenuWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Logf("status menu %s is loading %s",
			s.ID(),
			parsedURL,
		)
	})

	s.StatusMenu.Load(url, v...)
	if s.Err() != nil {
		Logf("status menu %T failed to load %s: %s",
			s.ID(),
			parsedURL,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Logf("status menu %s is rendering %T",
			s.ID(),
			c,
		)
	})

	s.StatusMenu.Render(c)
	if s.Err() != nil {
		Logf("status menu %s failed to render %T: %s",
			s.ID(),
			c,
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetIcon(name string) {
	WhenDebug(func() {
		Logf("status menu %s is setting icon to %s",
			s.ID(),
			name,
		)
	})

	s.StatusMenu.SetIcon(name)
	if s.Err() != nil {
		Logf("status menu %s failed to set icon: %s",
			s.ID(),
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) SetText(text string) {
	WhenDebug(func() {
		Logf("status menu %s is setting text to %s",
			s.ID(),
			text,
		)
	})

	s.StatusMenu.SetText(text)
	if s.Err() != nil {
		Logf("status menu %s failed to set text: %s",
			s.ID(),
			s.Err(),
		)
	}
}

func (s *statusMenuWithLogs) Close() {
	WhenDebug(func() {
		Logf("status menu %s is closing", s.ID())
	})

	s.StatusMenu.Close()
	if s.Err() != nil {
		Logf("status menu %s failed to close: %s",
			s.ID(),
			s.Err(),
		)
	}
}

func prettyConf(c interface{}) string {
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}
