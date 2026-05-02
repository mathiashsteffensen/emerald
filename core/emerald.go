package core

import (
	"emerald/debug"
	"emerald/object"
)

func (rt *Runtime) InitEmerald() {
	rt.Emerald = rt.DefineModule("Emerald")

	// This should be replaced by `extend self` when the feature is available
	rt.DefineMethod(rt.Emerald, "version", rt.emeraldVersion())
	rt.DefineSingletonMethod(rt.Emerald, "version", rt.emeraldVersion())
}

func (rt *Runtime) emeraldVersion() object.BuiltInMethod {
	version := rt.NewString(debug.EMERALD_VERSION)

	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return version
	}
}
