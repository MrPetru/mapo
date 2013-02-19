package addons

const (
	String = 0
)

type Addons interface {
	NewAddon(string) Addon
}

type Addon interface {
	SetConstructor(constructor)
}

type constructor func(EntityContainer)

type EntityContainer interface {
	NewEntity(string) Entity
	GetEntity(string) Entity
	GetEntityList(string) EntityList
}

type Entity interface {
	AddAttribute(string, int)
	SetAttribute(string, string)
	AddMethod(string, string, method)
	Restore(string, string) error
	Store(string) (string, error)
}

type EntityList interface {
	Restore(string) error
}

type method func(EntityContainer, RequestData) interface{}

type RequestData interface {
	GetValue(string) string
}
