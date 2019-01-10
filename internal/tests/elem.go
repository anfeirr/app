package tests

// UnsupportedConfig implements the app.ElemConfig interface.
type UnsupportedConfig struct{}

// Dump satisfies the app.ElemConfig interface.
func (c UnsupportedConfig) Dump() string {
	return ""
}
