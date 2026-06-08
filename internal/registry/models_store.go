package registry

// ModelStore persists model definitions.
// See docs/contracts/registry.md and tasks/006-registry-persistence.md.
type ModelStore interface {
	RegisterModel(def ModelDefinition) error
	GetModel(id string) (ModelDefinition, error)
	ListModels() []ModelDefinition
	Close() error
}
