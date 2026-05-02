package core

func (rt *Runtime) InitNumeric() {
	rt.Numeric = rt.DefineClass("Numeric", rt.Object)
}
