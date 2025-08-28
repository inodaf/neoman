package ports

type RegistryScope int

const (
	LocalScope RegistryScope = iota
	RemoteScope
)

type RegistryEntry struct {
	Scope       RegistryScope
	Owner       string
	Project     string
	ProjectPath string
}

type ProjectRegistry interface {
	HasEntry(entry RegistryEntry) bool
	AddEntry(entry RegistryEntry) error
}
