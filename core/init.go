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
	InitFloat()
	InitString()
	InitSymbol()

	// Initialize composite data types
	InitEnumerable()
	InitArray()
	InitHash()

	// Initialize exception hierarchy
	InitException()
	InitStandardError()
	InitArgumentError()
	InitTypeError()
}
