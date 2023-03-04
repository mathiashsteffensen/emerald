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
	InitComparable()
	InitNumeric()
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
	InitRuntimeError()
	InitArgumentError()
	InitTypeError()
	InitLoadError()
	InitNoMethodError()

	// Initialize remaining core classes & modules
	InitRegexp()
	InitMatchData()
	InitRange()
	InitIO()
	InitDir()
}
