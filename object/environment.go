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
	env := &environment{env: map[string]EmeraldValue{}}

	for k, v := range defaultEnvironment.env {
		env.env[k] = v
	}

	return env
}

func NewEnclosedEnvironment(outer Environment) Environment {
	env := &environment{env: map[string]EmeraldValue{}}

	env.outerEnv = outer

	return env
}

func (env *environment) Set(name string, val EmeraldValue) {
	env.env[name] = val
}

func (env *environment) Get(name string) (EmeraldValue, bool) {
	val, ok := env.env[name]
	if ok {
		return val, ok
	}

	if env.outerEnv != nil {
		return env.outerEnv.Get(name)
	}

	return val, ok
}
