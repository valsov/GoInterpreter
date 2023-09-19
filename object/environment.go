package object

func NewEnvironment() *Environment {
	store := map[string]Object{}
	return &Environment{store: store}
}

func NewEnclosedEnvironment(outerEnv *Environment) *Environment {
	newEnv := NewEnvironment()
	newEnv.outer = outerEnv
	return newEnv
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	// Find in current env first, then try outer envs
	obj, found := e.store[name]
	if !found && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, found
}

func (e *Environment) Set(name string, obj Object) {
	e.store[name] = obj
}
