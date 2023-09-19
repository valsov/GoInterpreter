package object

func NewEnvironment() *Environment {
	store := map[string]Object{}
	return &Environment{Store: store}
}

type Environment struct {
	Store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, found := e.Store[name]
	return obj, found
}

func (e *Environment) Set(name string, obj Object) {
	e.Store[name] = obj
}
