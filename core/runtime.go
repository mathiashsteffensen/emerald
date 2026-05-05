package core

import (
	"emerald/bytecode"
	"emerald/heap"
	"emerald/object"
)

type Runtime struct {
	Heap *heap.Heap

	CompileBlock func(fileName string, content string) *bytecode.Bytecode
	EvalBlock    func(block *object.ClosedBlock, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue
	Send         func(self object.EmeraldValue, name string, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue
	OnRaise      func(err object.EmeraldError)

	GlobalSymbolInternPool SymbolInternStore
	RequiredFilesHash      object.EmeraldValue

	// Core classes & modules
	ArgumentError object.EmeraldValue
	Array         object.EmeraldValue
	BasicObject   object.EmeraldValue
	Class         object.EmeraldValue
	Comparable    object.EmeraldValue
	Dir           object.EmeraldValue
	Emerald       object.EmeraldValue
	Enumerable    object.EmeraldValue
	Exception     object.EmeraldValue
	FalseClass    object.EmeraldValue
	File          object.EmeraldValue
	Float         object.EmeraldValue
	Hash          object.EmeraldValue
	IO            object.EmeraldValue
	Integer       object.EmeraldValue
	Kernel        object.EmeraldValue
	LoadError     object.EmeraldValue
	MatchData     object.EmeraldValue
	Module        object.EmeraldValue
	NameError     object.EmeraldValue
	NilClass      object.EmeraldValue
	NoMethodError object.EmeraldValue
	Numeric       object.EmeraldValue
	Object        object.EmeraldValue
	Range         object.EmeraldValue
	Regexp        object.EmeraldValue
	RuntimeError  object.EmeraldValue
	StandardError object.EmeraldValue
	String        object.EmeraldValue
	Symbol        object.EmeraldValue
	TCPServer     object.EmeraldValue
	TCPSocket     object.EmeraldValue
	Time          object.EmeraldValue
	TrueClass     object.EmeraldValue
	TypeError     object.EmeraldValue

	// Singletons
	TRUE       object.EmeraldValue
	FALSE      object.EmeraldValue
	NULL       object.EmeraldValue
	MainObject object.EmeraldValue
}

func NewRuntime() *Runtime {
	rt := &Runtime{
		Heap:                   heap.NewHeap(),
		GlobalSymbolInternPool: SymbolInternStore{},
		RequiredFilesHash:      object.EmeraldValue{},
	}

	return rt
}

func (rt *Runtime) Init() {
	// Initialize object hierarchy base
	rt.InitClass()
	rt.InitBasicObject()
	rt.InitKernel()
	rt.InitObject()
	rt.InitModule()

	// Initialize primitives
	rt.InitTrueClass()
	rt.InitFalseClass()
	rt.InitNilClass()
	rt.InitComparable()
	rt.InitNumeric()
	rt.InitInteger()
	rt.InitFloat()
	rt.InitString()
	rt.InitSymbol()

	// Initialize composite data types
	rt.InitEnumerable()
	rt.InitArray()
	rt.InitHash()

	// Initialize exception hierarchy
	rt.InitException()
	rt.InitStandardError()
	rt.InitRuntimeError()
	rt.InitNameError()
	rt.InitArgumentError()
	rt.InitTypeError()
	rt.InitLoadError()
	rt.InitNoMethodError()

	// Initialize our Emerald module
	rt.InitEmerald()

	// Initialize remaining core classes & modules
	rt.InitRegexp()
	rt.InitMatchData()
	rt.InitRange()
	rt.InitIO()
	rt.InitDir()
	rt.InitFile()

	// Networking core classes
	rt.InitTCPServer()
	rt.InitTCPSocket()

	rt.RequiredFilesHash = rt.NewHash()
}
