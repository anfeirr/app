package core

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// Window represents a window.
type Window struct {
	// The driver that manages the window.
	Driver app.Driver

	// The dom engine used by the window.
	Dom dom.Engine

	// The function to open the default browser.
	OpenDefaultBrowser func(rawurl string) error

	// The name of the javascript function to pass data to golang.
	JSToGo string

	id           string
	history      History
	compo        app.Compo
	isFocus      bool
	isFullscreen bool
	isMinimized  bool
	err          error
}

// ID satisfies the app.Window interface.
func (w *Window) ID() string {
	return w.id
}

// Err satisfies the app.Window interface.
func (w *Window) Err() error {
	return w.err
}

// WhenWindow satisfies the app.Window interface.
func (w *Window) WhenWindow(f func(w app.Window)) {
	f(w)
}

// WhenPage satisfies the app.Window interface.
func (w *Window) WhenPage(func(app.Page)) {}

// WhenWebView satisfies the app.Window interface.
func (w *Window) WhenWebView(f func(w app.WebView)) {
	f(w)
}

// WhenMenu satisfies the app.Window interface.
func (w *Window) WhenMenu(func(app.Menu)) {}

// WhenDockTile satisfies the app.Window interface.
func (w *Window) WhenDockTile(func(app.DockTile)) {}

// WhenStatusMenu satisfies the app.Window interface.
func (w *Window) WhenStatusMenu(func(app.StatusMenu)) {}

// Create creates and display the window.
func (w *Window) Create(c app.WindowConfig) {
	w.id = uuid.New().String()
	w.Dom.Sync = w.render

	if err := w.Driver.Call("windows.New", nil, struct {
		ID                string
		Title             string
		X                 float64
		Y                 float64
		Width             float64
		MinWidth          float64
		MaxWidth          float64
		Height            float64
		MinHeight         float64
		MaxHeight         float64
		BackgroundColor   string
		FrostedBackground bool
		FixedSize         bool
		CloseHidden       bool
		MinimizeHidden    bool
		TitlebarHidden    bool
	}{
		ID:                w.id,
		Title:             c.Title,
		X:                 c.X,
		Y:                 c.Y,
		Width:             c.Width,
		MinWidth:          c.MinWidth,
		MaxWidth:          c.MaxWidth,
		Height:            c.Height,
		MinHeight:         c.MinHeight,
		MaxHeight:         c.MaxHeight,
		BackgroundColor:   c.BackgroundColor,
		FrostedBackground: c.FrostedBackground,
		FixedSize:         c.FixedSize,
		CloseHidden:       c.CloseHidden,
		MinimizeHidden:    c.MinimizeHidden,
	}); err != nil {
		w.err = err
		return
	}

	w.Driver.Elems().Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}
}

// Load satisfies the app.Window interface.
func (w *Window) Load(urlFmt string, v ...interface{}) {
	u := fmt.Sprintf(urlFmt, v...)
	n := CompoNameFromURLString(u)

	// Redirect web page to default web browser.
	if !w.Driver.Compos().IsCompoRegistered(n) {
		w.err = w.OpenDefaultBrowser(u)
		return
	}

	if w.compo, w.err = w.Driver.Compos().NewCompo(n); w.err != nil {
		return
	}

	if u != w.history.Current() {
		w.history.NewEntry(u)
	}

	htmlConf := app.HTMLConfig{}
	if configurator, ok := w.compo.(app.Configurator); ok {
		htmlConf = configurator.Config()
	}

	if len(htmlConf.CSS) == 0 {
		htmlConf.CSS = file.Filenames(w.Driver.Resources("css"), ".css")
	}

	if len(htmlConf.Javascripts) == 0 {
		htmlConf.Javascripts = file.Filenames(w.Driver.Resources("js"), ".js")
	}

	page := dom.Page{
		Title:         htmlConf.Title,
		Metas:         htmlConf.Metas,
		CSS:           htmlConf.CSS,
		Javascripts:   htmlConf.Javascripts,
		GoRequest:     w.JSToGo,
		RootCompoName: n,
	}

	if w.err = w.Driver.Call("windows.Load", nil, struct {
		ID      string
		Title   string
		Page    string
		LoadURL string
		BaseURL string
	}{
		ID:      w.id,
		Title:   htmlConf.Title,
		Page:    page.String(),
		LoadURL: u,
		BaseURL: w.Driver.Resources(),
	}); w.err != nil {
		return
	}

	if w.err = w.Dom.New(w.compo); w.err != nil {
		return
	}

	if nav, ok := w.compo.(app.Navigable); ok {
		navURL, _ := url.Parse(u)
		nav.OnNavigate(navURL)
	}
}

// Compo satisfies the app.Window interface.
func (w *Window) Compo() app.Compo {
	return w.compo
}

// Contains satisfies the app.Window interface.
func (w *Window) Contains(c app.Compo) bool {
	return w.Dom.Contains(c)
}

// Render satisfies the app.Window interface.
func (w *Window) Render(c app.Compo) {
	w.err = w.Dom.Render(c)
}

func (w *Window) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encoding changes failed")
	}

	return w.Driver.Call("windows.Render", nil, struct {
		ID      string
		Changes string
	}{
		ID:      w.id,
		Changes: string(b),
	})
}

// Reload satisfies the app.Window interface.
func (w *Window) Reload() {
	u := w.history.Current()

	if len(u) == 0 {
		w.err = errors.New("no component loaded")
		return
	}

	w.Load(u)
}

// CanPrevious satisfies the app.Window interface.
func (w *Window) CanPrevious() bool {
	return w.history.CanPrevious()
}

// Previous satisfies the app.Window interface.
func (w *Window) Previous() {
	u := w.history.Previous()

	if len(u) == 0 {
		w.err = errors.New("no previous component")
		return
	}

	w.Load(u)
}

// CanNext satisfies the app.Window interface.
func (w *Window) CanNext() bool {
	return w.history.CanNext()
}

// Next satisfies the app.Window interface.
func (w *Window) Next() {
	u := w.history.Next()

	if len(u) == 0 {
		w.err = errors.New("no next component")
		return
	}

	w.Load(u)
}

// EvalJS satisfies the app.Window interface.
func (w *Window) EvalJS(result interface{}, eval string, args ...interface{}) error {
	ev, err := FormatJS(eval, args...)
	if err != nil {
		return err
	}

	out := struct {
		Result interface{}
	}{
		Result: result,
	}

	err = w.Driver.Call("windows.EvalJS", &out, struct {
		ID   string
		Eval string
	}{
		ID:   w.id,
		Eval: ev,
	})

	return err
}

// Position satisfies the app.Window interface.
func (w *Window) Position() (x, y float64) {
	out := struct {
		X float64
		Y float64
	}{}

	w.err = w.Driver.Call("windows.Position", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	return out.X, out.Y
}

// Move satisfies the app.Window interface.
func (w *Window) Move(x, y float64) {
	w.err = w.Driver.Call("windows.Move", nil, struct {
		ID string
		X  float64
		Y  float64
	}{
		ID: w.id,
		X:  x,
		Y:  y,
	})
}

// Center satisfies the app.Window interface.
func (w *Window) Center() {
	w.err = w.Driver.Call("windows.Center", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	out := struct {
		Width  float64
		Heigth float64
	}{}

	w.err = w.Driver.Call("windows.Size", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	return out.Width, out.Heigth
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	w.err = w.Driver.Call("windows.Resize", nil, struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.id,
		Width:  width,
		Height: height,
	})
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	w.err = w.Driver.Call("windows.Focus", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// IsFocus satisfies the app.Window interface.
func (w *Window) IsFocus() bool {
	return w.isFocus
}

// SetIsFocus set the focus status with the given value.
func (w *Window) SetIsFocus(v bool) {
	w.isFocus = v
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	if w.isFullscreen {
		w.err = nil
		return
	}

	w.err = w.Driver.Call("windows.ToggleFullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	if !w.isFullscreen {
		w.err = nil
		return
	}

	w.err = w.Driver.Call("windows.ToggleFullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// IsFullScreen satisfies the app.Window interface.
func (w *Window) IsFullScreen() bool {
	return w.isFullscreen
}

// SetIsFullScreen set the full screen status with the given value.
func (w *Window) SetIsFullScreen(v bool) {
	w.isFullscreen = v
}

// Minimize satisfies the app.Window interface.
func (w *Window) Minimize() {
	if w.isMinimized {
		w.err = nil
		return
	}

	w.err = w.Driver.Call("windows.ToggleMinimize", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// Deminimize satisfies the app.Window interface.
func (w *Window) Deminimize() {
	if !w.isMinimized {
		w.err = nil
		return
	}

	w.err = w.Driver.Call("windows.ToggleMinimize", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}

// IsMinimized satisfies the app.Window interface.
func (w *Window) IsMinimized() bool {
	return w.isMinimized
}

// SetIsMinimized set the full minimized status with the given value.
func (w *Window) SetIsMinimized(v bool) {
	w.isMinimized = v
}

// Close satisfies the app.Window interface.
func (w *Window) Close() {
	w.err = w.Driver.Call("windows.Close", nil, struct {
		ID string
	}{
		ID: w.id,
	})
}
