package app

// Elem is the interface that describes an app element.
type Elem interface {
	// ID returns the element identifier.
	ID() string

	// Contains reports whether the component is mounted in the element.
	Contains(Compo) bool

	// WhenWindow calls the given func when the element is a window.
	WhenWindow(func(Window))

	// WhenPage calls the given func when the element is a page.
	WhenPage(func(Page))

	// WebView calls the given func when the element displays web content.
	WhenWebView(func(WebView))

	// WhenStatusMenu calls the given func when the element is a menu.
	WhenMenu(func(Menu))

	// WhenDockTile calls the given func when the element is a dock tile.
	WhenDockTile(func(DockTile))

	// WhenStatusMenu calls the given func when the element is a status menu.
	WhenStatusMenu(func(StatusMenu))

	// Err returns the error that prevent the element to work.
	Err() error
}

// ElemWithCompo is the interface that describes an element that hosts
// components.
type ElemWithCompo interface {
	Elem

	// Load loads the page specified by the URL.
	// URL can be formated as fmt package functions.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string, v ...interface{})

	// Compo returns the loaded component.
	Compo() Compo

	// Render renders the component.
	Render(Compo)
}

// ElemStore is the interface that describes a store that contains app elements.
type ElemStore interface {
	// Put inserts or update the given element.
	Put(Elem)

	// Delete deletes the given element.
	Delete(Elem)

	// GetByID returns the element registered under the the given id.
	GetByID(string) Elem

	// GetByCompo returns the element where the given component is mounted.
	GetByCompo(Compo) Elem
}

// WebView is the interface that describe an element that displays web content.
type WebView interface {
	ElemWithCompo

	// Reload reloads the current page.
	Reload()

	// CanPrevious reports whether load the previous page is possible.
	CanPrevious() bool

	// Previous loads the previous page.
	Previous()

	// CanNext indicates if loading next page is possible.
	CanNext() bool

	// Next loads the next page.
	Next()

	// Evaluates a javascript expression formatted with the given values and
	// stores the result is the given value.
	//
	// It returns an error when result is not a pointer or when a given argument
	// can't be converted to json.
	EvalJS(result interface{}, eval string, args ...interface{}) error
}

// Closer is the interface that describes an element that can be closed.
type Closer interface {
	// Close closes the element and free its allocated resources.
	Close()
}

// NotificationConfig is a struct that describes a notification.
type NotificationConfig struct {
	Title     string
	Subtitle  string
	Text      string
	ImageName string
	Sound     bool

	OnReply func(reply string) `json:"-"`
}
