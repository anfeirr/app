// Package app is a package to build GUI apps with Go, HTML and CSS.
package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrNotSupported describes an error that occurs when an unsupported
	// feature is used.
	ErrNotSupported = errors.New("not supported")

	// ErrElemNotSet describes an error that reports if an element is set.
	ErrElemNotSet = errors.New("element not set")

	// ErrCompoNotMounted describes an error that reports whether a component
	// is mounted.
	ErrCompoNotMounted = errors.New("component not mounted")

	// Logger is a function that formats using the default formats for its
	// operands and logs the resulting string.
	// It is used by Log, Logf, Panic and Panicf to generate logs.
	Logger func(format string, a ...interface{})

	driver    Driver
	target    = "web"
	ui        = make(chan func(), 4096)
	factory   = NewFactory()
	events    = NewEventRegistry(ui)
	messages  = newMsgRegistry()
	whenDebug func(func())
)

const (
	// Running is the event emitted when the app starts to run.
	Running Event = "app.running"

	// Reopened is the event emitted when the app is reopened.
	Reopened Event = "app.reopened"

	// Focused is the event emitted when the app gets focus.
	Focused Event = "app.focused"

	// Blurred is the event emitted when the app loses focus.
	Blurred Event = "app.blurred"

	// OpenFilesRequested is the event emitted when the app is requested to
	// open files. The arg passed to subscribed funcs is a []string containing
	// the path of the requested files.
	OpenFilesRequested Event = "app.openFilesRequested"

	// OpenURLRequested is the event emitted when the app is requested to open
	// an URL. The arg passed to subscribed funcs is a *url.URL.
	OpenURLRequested Event = "app.openURLrequested"

	// Closed is the event emitted when the app is closed. Final cleanups
	// should be done by subscribing to this event.
	Closed Event = "app.closed"
)

func init() {
	EnableDebug(false)
}

// Import imports the given components into the app.
// Components must be imported in order the be used by the app package.
// This allows components to be created dynamically when they are found into
// markup.
func Import(c ...Compo) {
	for _, compo := range c {
		if _, err := factory.RegisterCompo(compo); err != nil {
			Panicf("import component failed: %s", err)
		}
	}
}

// Run runs the app with the given driver as backend.
func Run(drivers ...Driver) {
	for _, d := range drivers {
		if d.Target() == target {
			driver = d
			break
		}
	}

	if driver == nil {
		panic(errors.Errorf("no driver set for %s", target))
	}

	driver = DriverWithLogs(driver)

	driver.Run(DriverConfig{
		UI:      ui,
		Factory: factory,
		Events:  events,
	})
}

// CurrentDriver returns the current driver.
func CurrentDriver() Driver {
	return driver
}

// Name returns the application name.
//
// It panics if called before Run.
func Name() string {
	return driver.AppName()
}

// Resources returns the given path prefixed by the resources directory
// location.
// Resources should be used only for read only operations.
//
// It panics if called before Run.
func Resources(path ...string) string {
	return driver.Resources(path...)
}

// Storage returns the given path prefixed by the storage directory
// location.
//
// It panics if called before Run.
func Storage(path ...string) string {
	return driver.Storage(path...)
}

// Render renders the given component.
// It should be called when the display of component c have to be updated.
//
// It panics if called before Run.
func Render(c Compo) {
	driver.UI(func() {
		driver.Render(c)
	})
}

// New creates and displays the element described by the given configuration.
//
// It panics if called before Run.
func New(c ElemConfig) Elem {
	return driver.New(c)
}

// ElemByCompo returns the element where the given component is mounted.
//
// It panics if called before Run.
func ElemByCompo(c Compo) Elem {
	return driver.ElemByCompo(c)
}

// MenuBar returns the menu bar.
//
// It panics if called before Run.
func MenuBar() Menu {
	return driver.MenuBar()
}

// Dock returns the dock tile.
//
// It panics if called before Run.
func Dock() DockTile {
	return driver.DockTile()
}

// Close stops the app. It makes Run() to return an error.
//
// It panics if called before Run.
func Close() {
	driver.Close()
}

// UI calls a function on the UI goroutine.
func UI(f func()) {
	driver.UI(f)
}

// Handle handles the message for the given key.
func Handle(key string, h Handler) {
	messages.handle(key, h)
}

// Post posts the given messages.
// Messages are handled in another goroutine.
func Post(msgs ...Msg) {
	messages.post(msgs...)
}

// NewMsg creates a message.
func NewMsg(key string) Msg {
	return &msg{key: key}
}

// Emit emits the event with the given arguments.
func Emit(e Event, args ...interface{}) {
	events.Emit(e, args...)
}

// NewSubscriber creates an event subscriber to return when implementing the
// app.EventSubscriber interface.
func NewSubscriber() *Subscriber {
	return &Subscriber{
		Events: events,
	}
}

// Log formats using the default formats for its operands and logs the resulting
// string.
// Spaces are always added between operands and a newline is appended.
func Log(a ...interface{}) {
	format := ""

	for range a {
		format += "%v "
	}

	format = format[:len(format)-1]
	Logger(format, a...)
}

// Logf formats according to a format specifier and logs the resulting string.
func Logf(format string, a ...interface{}) {
	Logger(format, a...)
}

// Panic is equivalent to Log() followed by a call to panic().
func Panic(a ...interface{}) {
	Log(a...)
	panic(strings.TrimSpace(fmt.Sprintln(a...)))
}

// Panicf is equivalent to Logf() followed by a call to panic().
func Panicf(format string, a ...interface{}) {
	Logf(format, a...)
	panic(fmt.Sprintf(format, a...))
}

// EnableDebug is a function that set whether debug mode is enabled.
func EnableDebug(v bool) {
	whenDebug = func(f func()) {}

	if v {
		whenDebug = func(f func()) {
			f()
		}
	}
}

// WhenDebug execute the given function when debug mode is enabled.
func WhenDebug(f func()) {
	whenDebug(f)
}

// CompoName returns the name of the given component.
// The returned name is the one to use in html tags.
func CompoName(c Compo) string {
	v := reflect.ValueOf(c)
	v = reflect.Indirect(v)

	name := strings.ToLower(v.Type().String())
	return strings.TrimPrefix(name, "main.")
}

// Pretty is an helper function that returns a prettified string representation
// of the given value.
// Returns an empty string if the value can't be prettified.
func Pretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}
