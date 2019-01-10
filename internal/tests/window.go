package tests

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
)

func newWindow(d *Driver, c app.WindowConfig) app.Window {
	w := &core.Window{
		Driver: d,
		Dom: dom.Engine{
			Factory:   d.factory,
			Resources: d.Resources,
			AttrTransforms: []dom.Transform{
				dom.JsToGoHandler,
				dom.HrefCompoFmt,
			},
			UI: d.UI,
		},
		OpenDefaultBrowser: func(string) error {
			return app.ErrNotSupported
		},
		JSToGo: "alert",
	}

	w.Create(c)
	if w.Err() != nil {
		return w
	}

	d.elems.Put(w)

	if len(c.URL) != 0 {
		w.Load(c.URL)
	}

	return w
}

// func testWindow(t *testing.T, w app.Window) {
// 	// app.Elem
// 	called := false
// 	w.WhenWindow(func(w app.Window) {
// 		called = true
// 	})
// 	assert.True(t, called)

// 	called = false
// 	w.WhenPage(func(p app.Page) {
// 		called = true
// 	})
// 	assert.False(t, called)

// 	called = false
// 	w.WhenWebView(func(w app.WebView) {
// 		called = true
// 	})
// 	assert.True(t, called)

// 	called = false
// 	w.WhenMenu(func(m app.Menu) {
// 		called = true
// 	})
// 	assert.False(t, called)

// 	called = false
// 	w.WhenDockTile(func(d app.DockTile) {
// 		called = true
// 	})
// 	assert.False(t, called)

// 	called = false
// 	w.WhenStatusMenu(func(s app.StatusMenu) {
// 		called = true
// 	})
// 	assert.False(t, called)

// 	w.WhenErr(func(err error) {
// 		t.Log(err)
// 	})

// 	t.Run("navigator", func(t *testing.T) {
// 		testWebView(t, w, false)
// 	})

// 	t.Run("compo", func(t *testing.T) {
// 		testElemWithCompo(t, w)
// 	})

// 	w.Position()
// 	assertElem(t, w)

// 	w.Move(42, 42)
// 	assertElem(t, w)

// 	w.Center()
// 	assertElem(t, w)

// 	w.Size()
// 	assertElem(t, w)

// 	w.Resize(42, 42)
// 	assertElem(t, w)

// 	w.Focus()
// 	assertElem(t, w)

// 	w.FullScreen()
// 	assertElem(t, w)

// 	w.ExitFullScreen()
// 	assertElem(t, w)

// 	w.Minimize()
// 	assertElem(t, w)

// 	w.Deminimize()
// 	assertElem(t, w)

// 	w.Close()
// 	assertElem(t, w)
// }
