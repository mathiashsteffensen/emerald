package core

import (
	"emerald/debug"
	"emerald/object"
)

var Emerald *object.Module

func InitEmerald() {
	Emerald = DefineModule("Emerald")

	// This should be replaced by `extend self` when the feature is available
	DefineMethod(Emerald, "version", emeraldVersion())
	DefineSingletonMethod(Emerald, "version", emeraldVersion())
}

func emeraldVersion() object.BuiltInMethod {
	version := NewString(debug.EMERALD_VERSION)

	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return version
	}
}
