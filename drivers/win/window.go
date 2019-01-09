// +build windows

package win

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

// Window implements the app.Window interface.
type Window struct {
	core.Window

	id           string
	dom          dom.Engine
	history      core.History
	compo        app.Compo
	isFullScreen bool
	isMinimized  bool
	isFocus      bool
}

func newWindow(c app.WindowConfig) *Window {
	id := uuid.New().String()

	w := &Window{
		id: id,
		dom: dom.Engine{
			Factory:   driver.factory,
			Resources: resourcesDir,
			AttrTransforms: []dom.Transform{
				dom.JsToGoHandler,
				dom.HrefCompoFmt,
			},
			UI: driver.UI,
		},
		isFocus: true,
	}

	w.dom.Sync = w.render

	in := struct {
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
	}

	in.MinWidth, in.MaxWidth = normalizeWidowSize(in.MinWidth, in.MaxWidth)
	in.MinHeight, in.MaxHeight = normalizeWidowSize(in.MinHeight, in.MaxHeight)

	if err := driver.winRPC.Call("windows.New", nil, in); err != nil {
		w.SetErr(err)
		return w
	}

	driver.elems.Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}

	return w
}

func normalizeWidowSize(min, max float64) (float64, float64) {
	min = math.Max(0, min)
	min = math.Min(min, 10000)

	if max == 0 {
		max = 10000
	}
	max = math.Max(0, max)
	max = math.Min(max, 10000)

	min = math.Min(min, max)
	return min, max
}

// ID satisfies the app.Window interface.
func (w *Window) ID() string {
	return w.id
}

// Load satisfies the app.Window interface.
func (w *Window) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		w.SetErr(err)
	}()

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	// Redirect web page to default web browser.
	if !driver.factory.IsCompoRegistered(n) {
		err = exec.Command("start", u).Run()
		return
	}

	var c app.Compo
	if c, err = driver.factory.NewCompo(n); err != nil {
		return
	}

	w.compo = c

	if u != w.history.Current() {
		w.history.NewEntry(u)
	}

	htmlConf := app.HTMLConfig{}
	if configurator, ok := c.(app.Configurator); ok {
		htmlConf = configurator.Config()
	}

	if len(htmlConf.CSS) == 0 {
		htmlConf.CSS = file.Filenames(driver.Resources("css"), ".css")
	}

	if len(htmlConf.Javascripts) == 0 {
		htmlConf.Javascripts = file.Filenames(driver.Resources("js"), ".js")
	}

	page := dom.Page{
		Title:         htmlConf.Title,
		Metas:         htmlConf.Metas,
		CSS:           pageFiles(htmlConf.CSS),
		Javascripts:   pageFiles(htmlConf.Javascripts),
		GoRequest:     "window.external.notify",
		RootCompoName: n,
	}

	if err = driver.winRPC.Call("windows.Load", nil, struct {
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
		BaseURL: driver.Resources(),
	}); err != nil {
		return
	}

	err = w.dom.New(c)

	if nav, ok := c.(app.Navigable); ok {
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
	return w.dom.Contains(c)
}

// Render satisfies the app.Window interface.
func (w *Window) Render(c app.Compo) {
	w.SetErr(w.dom.Render(c))
}

func (w *Window) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encode changes failed")
	}

	return driver.winRPC.Call("windows.Render", nil, struct {
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
		w.SetErr(errors.New("no component loaded"))
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
		w.SetErr(errors.New("no previous component"))
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
		w.SetErr(errors.New("no next component"))
		return
	}

	w.Load(u)
}

// EvalJS satisfies the app.Window interface.
func (w *Window) EvalJS(result interface{}, eval string, args ...interface{}) error {
	ev, err := core.FormatJS(eval, args...)
	if err != nil {
		return err
	}

	out := struct {
		Result interface{}
	}{
		Result: result,
	}

	err = driver.platform.Call("windows.EvalJS", &out, struct {
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

	err := driver.winRPC.Call("windows.Position", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
	return out.X, out.Y
}

// Size satisfies the app.Window interface.
func (w *Window) Size() (width, height float64) {
	out := struct {
		Width  float64
		Heigth float64
	}{}

	err := driver.winRPC.Call("windows.Size", &out, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
	return out.Width, out.Heigth
}

// Resize satisfies the app.Window interface.
func (w *Window) Resize(width, height float64) {
	err := driver.winRPC.Call("windows.Resize", nil, struct {
		ID     string
		Width  float64
		Height float64
	}{
		ID:     w.id,
		Width:  width,
		Height: height,
	})

	w.SetErr(err)
}

// Focus satisfies the app.Window interface.
func (w *Window) Focus() {
	err := driver.winRPC.Call("windows.Focus", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// IsFocus satisfies the app.Window interface.
func (w *Window) IsFocus() bool {
	return w.isFocus
}

// FullScreen satisfies the app.Window interface.
func (w *Window) FullScreen() {
	err := driver.winRPC.Call("windows.FullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// ExitFullScreen satisfies the app.Window interface.
func (w *Window) ExitFullScreen() {
	err := driver.winRPC.Call("windows.ExitFullScreen", nil, struct {
		ID string
	}{
		ID: w.id,
	})

	w.SetErr(err)
}

// IsFullScreen satisfies the app.Window interface.
func (w *Window) IsFullScreen() bool {
	return w.isFullScreen
}

func onWindowCallback(w *Window, in map[string]interface{}) {
	mappingStr := in["Mapping"].(string)

	var m dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &m); err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	c, err := w.dom.CompoByID(m.CompoID)
	if err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	app.Render(c)
}

func onWindowNavigate(w *Window, in map[string]interface{}) {
	e := app.ElemByCompo(w.Compo())

	e.WhenWindow(func(w app.Window) {
		w.Load(in["URL"].(string))
	})
}

func onWindowResize(w *Window, in map[string]interface{}) {
	driver.events.Emit(app.WindowResized, w)
}

func onWindowFocus(w *Window, in map[string]interface{}) {
	w.isFocus = true
	driver.events.Emit(app.WindowFocused, w)
}

func onWindowBlur(w *Window, in map[string]interface{}) {
	w.isFocus = false
	driver.events.Emit(app.WindowBlurred, w)
}

func onWindowFullScreen(w *Window, in map[string]interface{}) {
	w.isFullScreen = true
	driver.events.Emit(app.WindowEnteredFullScreen, w)
}

func onWindowExitFullScreen(w *Window, in map[string]interface{}) {
	w.isFullScreen = false
	driver.events.Emit(app.WindowExitedFullScreen, w)
}

func onWindowClose(w *Window, in map[string]interface{}) {
	fmt.Println("gonna close window")
	driver.events.Emit(app.WindowClosed, w)
	w.dom.Close()
	driver.elems.Delete(w)
	fmt.Println("bye")
}

func handleWindow(h func(w *Window, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*Window), in)
	}
}

func resourcesDir(p ...string) string {
	r := path.Join(p...)
	r = strings.TrimLeft(r, "/")
	return "ms-appx-web:///Resources/" + r
}

func pageFiles(files []string) []string {
	pfiles := make([]string, len(files))
	resources := driver.Resources()

	for i, f := range files {
		f = strings.TrimPrefix(f, resources)
		f = strings.Replace(f, "\\", "/", -1)
		pfiles[i] = resourcesDir(f)
	}

	return pfiles
}
