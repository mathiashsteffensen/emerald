package core

func init() {
	// Initialize object hierarchy base
	InitBasicObject()
	InitKernel()
	InitObject()

	// Initialize primitives
	InitTrueClass()
	InitFalseClass()
	InitNilClass()
	InitInteger()
	InitString()
	InitSymbol()

	// Initialize composite data types
	InitArray()

	// Initialize exception hierarchy
	InitException()
	InitStandardError()
	InitArgumentError()
}
