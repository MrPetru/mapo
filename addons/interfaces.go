package addons

const (
	String = 0
)

type Addons interface {
	NewAddon(string) Addon
}

type Addon interface {
	//SetConstructor(constructor)
	SetConstructor(func(EntityContainer))
}

//type constructor func(EntityContainer)

type EntityContainer interface {
	NewEntity(string) Entity
	GetEntity(string) Entity
	GetEntityList(string) EntityList
}

type Entity interface {
	AddAttribute(string, int)
	SetAttribute(string, string)
	AddMethod(string, string, Method)
	Restore(string) error
	Store() (string, error)
}

type EntityList interface {
	Restore() error
}

type Method func(EntityContainer, RequestData) interface{}

type RequestData interface {
	GetValue(string) string
}
