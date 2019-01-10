package app

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Returns the targetted operating system name.
	Target() string

	// Calls the platform method with the given input and stores the result into
	// the given ouput. Result is ignored if the output is nil.
	//
	// It panics if the given output is not a pointer.
	Call(method string, out interface{}, in interface{}) error

	// Runs the application with the components registered in the given factory.
	Run(DriverConfig)

	// Returns the component factory used to create components.
	Factory() *Factory

	// Returns the appliction name.
	AppName() string

	// Returns the given path prefixed by the resources directory location.
	Resources(path ...string) string

	// Returns the given path prefixed by the storage directory location.
	Storage(path ...string) string

	// Renders the given component.
	Render(Compo)

	// Creates and displays the element described in the given configuration.
	New(ElemConfig) Elem

	// Returns the element where the given component is mounted.
	ElemByCompo(Compo) Elem

	// Returns the current menu bar element.
	MenuBar() Menu

	// Returns the dock tile element.
	DockTile() DockTile

	// UI calls a function on the UI goroutine.
	UI(func())

	// Close stops the driver. It makes Run() to return an error.
	Close()
}

// DriverConfig contains driver configuration.
type DriverConfig struct {
	// The channel to send function to execute on UI goroutine.
	UI chan func()

	// The factory used to create components.
	Factory *Factory

	// The event registery to emit events.
	Events *EventRegistry
}
