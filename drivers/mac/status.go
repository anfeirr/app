// +build darwin,amd64

package mac

// // StatusMenu represents a menu that lives in the status bar.
// type StatusMenu struct {
// 	Menu
// }

// func newStatusMenu(c app.StatusMenuConfig) *StatusMenu {
// 	s := &StatusMenu{
// 		Menu: Menu{
// 			id: uuid.New().String(),
// 			dom: dom.Engine{
// 				Factory:   driver.factory,
// 				Resources: driver.Resources,
// 				AllowedNodes: []string{
// 					"menu",
// 					"menuitem",
// 				},
// 				UI: driver.UI,
// 			},
// 			typ:            "status menu",
// 			keepWhenClosed: true,
// 		},
// 	}

// 	s.dom.Sync = s.render

// 	if err := driver.platform.Call("statusMenus.New", nil, struct {
// 		ID   string
// 		Text string
// 		Icon string
// 	}{
// 		ID:   s.id,
// 		Text: c.Text,
// 		Icon: c.Icon,
// 	}); err != nil {
// 		s.SetErr(err)
// 		return s
// 	}

// 	driver.elems.Put(s)

// 	if len(c.URL) != 0 {
// 		s.Load(c.URL)
// 	}

// 	return s
// }

// // WhenStatusMenu satisfies the app.StatusMenu interface.
// func (s *StatusMenu) WhenStatusMenu(f func(app.StatusMenu)) {
// 	f(s)
// }

// // Load the app.StatusMenu interface.
// func (s *StatusMenu) Load(urlFmt string, v ...interface{}) {
// 	s.Menu.Load(urlFmt, v...)
// 	if s.Err() != nil {
// 		return
// 	}

// 	err := driver.platform.Call("statusMenus.SetMenu", nil, struct {
// 		ID string
// 	}{
// 		ID: s.id,
// 	})

// 	s.SetErr(err)
// }

// // SetIcon satisfies the app.StatusMenu interface.
// func (s *StatusMenu) SetIcon(path string) {
// 	if _, err := os.Stat(path); err != nil && len(path) != 0 {
// 		s.SetErr(err)
// 		return
// 	}

// 	err := driver.platform.Call("statusMenus.SetIcon", nil, struct {
// 		ID   string
// 		Icon string
// 	}{
// 		ID:   s.id,
// 		Icon: path,
// 	})

// 	s.SetErr(err)
// }

// // SetText satisfies the app.StatusMenu interface.
// func (s *StatusMenu) SetText(text string) {
// 	err := driver.platform.Call("statusMenus.SetText", nil, struct {
// 		ID   string
// 		Text string
// 	}{
// 		ID:   s.id,
// 		Text: text,
// 	})

// 	s.SetErr(err)
// }

// // Close satisfies the app.StatusMenu interface.
// func (s *StatusMenu) Close() {
// 	err := driver.platform.Call("statusMenus.Close", nil, struct {
// 		ID string
// 	}{
// 		ID: s.id,
// 	})

// 	s.SetErr(err)
// 	driver.elems.Delete(s)
// }
