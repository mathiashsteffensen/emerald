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
	ArgumentError *object.Class
	Array         *object.Class
	BasicObject   *object.Class
	Class         *object.Class
	Comparable    *object.Module
	Dir           *object.Class
	Emerald       *object.Module
	Enumerable    *object.Module
	Exception     *object.Class
	FalseClass    *object.Class
	File          *object.Class
	Float         *object.Class
	Hash          *object.Class
	IO            *object.Class
	Integer       *object.Class
	Kernel        *object.Module
	LoadError     *object.Class
	MatchData     *object.Class
	Module        *object.Class
	NameError     *object.Class
	NilClass      *object.Class
	NoMethodError *object.Class
	Numeric       *object.Class
	Object        *object.Class
	Range         *object.Class
	Regexp        *object.Class
	RuntimeError  *object.Class
	StandardError *object.Class
	String        *object.Class
	Symbol        *object.Class
	TCPServer     *object.Class
	TCPSocket     *object.Class
	Time          *object.Class
	TrueClass     *object.Class
	TypeError     *object.Class

	// Singletons
	TRUE       object.EmeraldValue
	FALSE      object.EmeraldValue
	NULL       *object.Instance
	MainObject *object.Instance
}

func NewRuntime() *Runtime {
	rt := &Runtime{
		Heap:                   heap.NewHeap(),
		GlobalSymbolInternPool: SymbolInternStore{},
		RequiredFilesHash:      nil,
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
