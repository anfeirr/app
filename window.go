package app

// Window is the interface that describes a window.
type Window interface {
	View
	Closer

	// Evaluates a javascript expression formatted with the given values and
	// stores the result in the given value.
	//
	// It returns an error when result is not a pointer or when a given argument
	// can't be converted to json.
	EvalJS(result interface{}, eval string, args ...interface{}) error

	// Position returns the window position.
	Position() (x, y float64)

	// Move moves the window to the position (x, y).
	Move(x, y float64)

	// Center moves the window to the center of the screen.
	Center()

	// Size returns the window size.
	Size() (width, height float64)

	// Resize resizes the window to width * height.
	Resize(width, height float64)

	// Focus gives the focus to the window.
	// The window will be put in front, above the other elements.
	Focus()

	// Reports whether the window is the focus.
	IsFocus() bool

	// FullScreen takes the window into full screen mode.
	FullScreen()

	// ExitFullScreen takes the window out of full screen mode.
	ExitFullScreen()

	// Reports whether the window is in full screen mode.
	IsFullScreen() bool

	// Minimize takes the window into minimized mode.
	Minimize()

	// Deminimize takes the window out of minimized mode.
	Deminimize()

	// Reports whether the window is minimized.
	IsMinimized() bool
}

// WindowConfig is a struct that describes a window. It implements the
// app.ElemConfig interface.
type WindowConfig struct {
	// The URL of the component to load when the window is created.
	URL string `json:",omitempty"`

	// The title.
	Title string `json:",omitempty"`

	// The default position on x axis.
	X float64 `json:",omitempty"`

	// The default position on y axis.
	Y float64 `json:",omitempty"`

	// The default width.
	Width float64 `json:",omitempty"`

	// The minimum width.
	MinWidth float64 `json:",omitempty"`

	// The maximum width.
	MaxWidth float64 `json:",omitempty"`

	// The default height.
	Height float64 `json:",omitempty"`

	// The minimum height.
	MinHeight float64 `json:",omitempty"`

	// The maximum height.
	MaxHeight float64 `json:",omitempty"`

	// The background color (#rrggbb).
	BackgroundColor string `json:",omitempty"`

	// Enables frosted effect.
	FrostedBackground bool `json:",omitempty"`

	// Reports whether the window is resizable.
	FixedSize bool `json:",omitempty"`

	// Reports whether the close button is hidden.
	CloseHidden bool `json:",omitempty"`

	// Reports whether the minimize button is hidden.
	MinimizeHidden bool `json:",omitempty"`
}

// Dump satisfies the app.ElemConfig interface.
func (c WindowConfig) Dump() string {
	return Pretty(c)
}

// Constants that enumerates window events.
const (
	// WindowMoved is the event emitted when a window is moved. The arg passed
	// to subscribed funcs is a app.Window.
	WindowMoved Event = "app.window.moved"

	// WindowResized is the event emitted when a window is resized. The arg
	// passed to subscribed funcs is a app.Window.
	WindowResized Event = "app.window.resized"

	// WindowFocused is the event emitted when a window gets focus. The arg
	// passed to subscribed funcs is a app.Window.
	WindowFocused Event = "app.window.focused"

	// WindowBlurred is the event emitted when a window loses focus. The arg
	// passed to subscribed funcs is a app.Window.
	WindowBlurred Event = "app.window.blurred"

	// WindowEnteredFullScreen is the event emitted when a window goes full
	// screen. The arg passed to subscribed funcs is a app.Window.
	WindowEnteredFullScreen Event = "app.window.enteredFullscreen"

	// WindowExitedFullScreen is the event emitted when a window exits full
	// screen. The arg passed to subscribed funcs is a app.Window.
	WindowExitedFullScreen Event = "app.window.exitedFullscreen"

	// WindowMinimized is the event emitted when a window is minimized. The arg
	// passed to subscribed funcs is a app.Window.
	WindowMinimized Event = "app.window.minimized"

	// WindowDeminimized is the event emitted when a window is deminimized. The
	// arg passed to subscribed funcs is a app.Window.
	WindowDeminimized Event = "app.window.deminimized"

	// WindowClosed is the event emitted when a window is closed. The arg passed
	// to subscribed funcs is a app.Window.
	WindowClosed Event = "app.window.closed"
)
