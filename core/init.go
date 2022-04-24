package core

func init() {
	// Initialize object hierarchy base
	InitClass()
	InitBasicObject()
	InitKernel()
	InitObject()
	InitModule()

	// Initialize primitives
	InitTrueClass()
	InitFalseClass()
	InitNilClass()
	InitInteger()
	InitString()
	InitSymbol()

	// Initialize composite data types
	InitArray()
	InitHash()

	// Initialize exception hierarchy
	InitException()
	InitStandardError()
	InitArgumentError()
}
