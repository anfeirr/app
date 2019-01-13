// +build darwin,amd64

package mac

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
)

func newWindow(c app.WindowConfig) app.Window {
	w := &core.Window{
		Driver: driver,
		Dom: dom.Engine{
			Factory:   driver.factory,
			Resources: driver.Resources,
			AttrTransforms: []dom.Transform{
				dom.JsToGoHandler,
				dom.HrefCompoFmt,
			},
			UI: driver.UI,
		},
		OpenDefaultBrowser: openDefaultBrowser,
		JSToGo:             "window.webkit.messageHandlers.golangRequest.postMessage",
	}

	if c.Width == 0 {
		c.Width = 1280
	}

	if c.Height == 0 {
		c.Height = 720
	}

	c.MinWidth, c.MaxWidth = normalizeWidowSize(c.MinWidth, c.MaxWidth)
	c.MinHeight, c.MaxHeight = normalizeWidowSize(c.MinHeight, c.MaxHeight)

	w.Create(c)
	if w.Err() != nil {
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

func openDefaultBrowser(url string) error {
	return exec.Command("open", url).Run()
}

func onWindowCallback(w *core.Window, in map[string]interface{}) {
	mappingStr := in["Mapping"].(string)

	var m dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &m); err != nil {
		app.Logf("window callback failed: %s", err)
		return
	}

	if m.Override == "Files" {
		data, _ := json.Marshal(driver.droppedFiles)
		driver.droppedFiles = nil

		m.JSONValue = strings.Replace(
			m.JSONValue,
			`"FileOverride":"xxx"`,
			fmt.Sprintf(`"Files":%s`, data),
			1,
		)
	}

	c, err := w.Dom.CompoByID(m.CompoID)
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

func onWindowNavigate(w *core.Window, in map[string]interface{}) {
	e := app.ElemByCompo(w.Compo())

	e.WhenWindow(func(w app.Window) {
		w.Load(in["URL"].(string))
	})
}

func onWindowAlert(w *core.Window, in map[string]interface{}) {
	app.Logf("%s", in["Alert"])
}

func onWindowMove(w *core.Window, in map[string]interface{}) {
	driver.events.Emit(app.WindowMoved, w)
}

func onWindowResize(w *core.Window, in map[string]interface{}) {
	driver.events.Emit(app.WindowResized, w)
}

func onWindowFocus(w *core.Window, in map[string]interface{}) {
	w.SetIsFocus(true)
	driver.events.Emit(app.WindowFocused, w)
}

func onWindowBlur(w *core.Window, in map[string]interface{}) {
	w.SetIsFocus(false)
	driver.events.Emit(app.WindowBlurred, w)
}

func onWindowFullScreen(w *core.Window, in map[string]interface{}) {
	w.SetIsFullScreen(true)
	driver.events.Emit(app.WindowEnteredFullScreen, w)
}

func onWindowExitFullScreen(w *core.Window, in map[string]interface{}) {
	w.SetIsFullScreen(false)
	driver.events.Emit(app.WindowExitedFullScreen, w)
}

func onWindowMinimize(w *core.Window, in map[string]interface{}) {
	w.SetIsMinimized(true)
	driver.events.Emit(app.WindowMinimized, w)
}

func onWindowDeminimize(w *core.Window, in map[string]interface{}) {
	w.SetIsMinimized(false)
	driver.events.Emit(app.WindowDeminimized, w)
}

func onWindowClose(w *core.Window, in map[string]interface{}) {
	driver.events.Emit(app.WindowClosed, w)
	w.Dom.Close()
	driver.elems.Delete(w)
}

func handleWindow(h func(w *core.Window, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*core.Window), in)
	}
}
