package object

type Environment interface {
	Set(name string, val EmeraldValue)
	Get(name string) (EmeraldValue, bool)
}

type environment struct {
	outerEnv Environment
	env      map[string]EmeraldValue
}

var defaultEnvironment = &environment{nil, map[string]EmeraldValue{}}

func NewEnvironment() Environment {
	return &environment{env: defaultEnvironment.env}
}

func NewEnclosedEnvironment(outer Environment) Environment {
	env := NewEnvironment().(*environment)

	env.outerEnv = outer

	return env
}

func (env *environment) Set(name string, val EmeraldValue) {
	env.env[name] = val
}

func (env *environment) Get(name string) (EmeraldValue, bool) {
	val, ok := env.env[name]
	return val, ok
}
