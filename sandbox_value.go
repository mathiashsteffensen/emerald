package emerald

import (
	"emerald/core"
	"emerald/internal/sandboxwire"
	"emerald/object"
	"fmt"
	"math"
)

const maxSandboxValueDepth = 64

func encodeSandboxValue(value object.EmeraldValue) (*sandboxwire.Value, error) {
	return encodeSandboxValueAt(value, map[object.HeapObject]struct{}{}, 0)
}

func encodeSandboxValueAt(
	value object.EmeraldValue,
	seenReferences map[object.HeapObject]struct{},
	depth int,
) (*sandboxwire.Value, error) {
	if depth > maxSandboxValueDepth {
		return nil, unsupportedSandboxValue(value, "maximum nesting depth exceeded")
	}

	switch value.TypeID {
	case object.NIL_VALUE:
		return &sandboxwire.Value{Type: sandboxwire.ValueNil}, nil
	case object.TRUE_VALUE:
		return &sandboxwire.Value{Type: sandboxwire.ValueBool, Bool: true}, nil
	case object.FALSE_VALUE:
		return &sandboxwire.Value{Type: sandboxwire.ValueBool, Bool: false}, nil
	case object.INTEGER_VALUE:
		return &sandboxwire.Value{Type: sandboxwire.ValueInt, Int: int64(value.Num)}, nil
	case object.FLOAT_VALUE:
		return &sandboxwire.Value{Type: sandboxwire.ValueFloat, FloatBits: value.Num}, nil
	}

	switch heapValue := value.Heap.(type) {
	case *core.StringInstance:
		if err := markSandboxReference(value, seenReferences); err != nil {
			return nil, err
		}
		return &sandboxwire.Value{Type: sandboxwire.ValueString, String: heapValue.Value}, nil
	case *core.SymbolInstance:
		if err := markSandboxReference(value, seenReferences); err != nil {
			return nil, err
		}
		return &sandboxwire.Value{Type: sandboxwire.ValueSymbol, Symbol: heapValue.Value}, nil
	case *core.ArrayInstance:
		if err := markSandboxReference(value, seenReferences); err != nil {
			return nil, err
		}
		encoded := &sandboxwire.Value{
			Type:  sandboxwire.ValueArray,
			Array: make([]sandboxwire.Value, len(heapValue.Value)),
		}
		for i, item := range heapValue.Value {
			itemValue, err := encodeSandboxValueAt(item, seenReferences, depth+1)
			if err != nil {
				return nil, err
			}
			encoded.Array[i] = *itemValue
		}
		return encoded, nil
	case *core.HashInstance:
		if err := markSandboxReference(value, seenReferences); err != nil {
			return nil, err
		}
		encoded := &sandboxwire.Value{Type: sandboxwire.ValueHash}
		var encodeErr error
		heapValue.Each(func(key object.EmeraldValue, item object.EmeraldValue) {
			if encodeErr != nil {
				return
			}
			encodedKey, err := encodeSandboxValueAt(key, seenReferences, depth+1)
			if err != nil {
				encodeErr = err
				return
			}
			encodedValue, err := encodeSandboxValueAt(item, seenReferences, depth+1)
			if err != nil {
				encodeErr = err
				return
			}
			encoded.Hash = append(encoded.Hash, sandboxwire.HashEntry{
				Key:   *encodedKey,
				Value: *encodedValue,
			})
		})
		if encodeErr != nil {
			return nil, encodeErr
		}
		return encoded, nil
	default:
		return nil, unsupportedSandboxValue(value, "type is not transferable")
	}
}

func markSandboxReference(value object.EmeraldValue, seen map[object.HeapObject]struct{}) error {
	if _, ok := seen[value.Heap]; ok {
		return unsupportedSandboxValue(value, "aliases and cycles are not transferable")
	}
	seen[value.Heap] = struct{}{}
	return nil
}

func unsupportedSandboxValue(value object.EmeraldValue, reason string) error {
	typ := value.TypeID.String()
	if value.Heap != nil {
		switch value.Heap.(type) {
		case *core.StringInstance:
			typ = "String"
		case *core.SymbolInstance:
			typ = "Symbol"
		case *core.ArrayInstance:
			typ = "Array"
		case *core.HashInstance:
			typ = "Hash"
		default:
			typ = fmt.Sprintf("%T", value.Heap)
		}
	}
	return &SandboxUnsupportedResultError{Type: typ, Reason: reason}
}

func decodeSandboxValue(encoded *sandboxwire.Value, frameLimit int) (object.EmeraldValue, error) {
	if encoded == nil {
		return object.EmeraldValue{}, sandboxInvalidWorkerValue("missing result")
	}
	rt := core.NewRuntime()
	rt.InitSandbox()
	budget := sandboxDecodeBudget{remaining: frameLimit}
	return decodeSandboxValueAt(rt, *encoded, &budget, 0)
}

type sandboxDecodeBudget struct {
	remaining int
}

func (b *sandboxDecodeBudget) consume(amount int) error {
	if amount < 0 || amount > b.remaining {
		return sandboxInvalidWorkerValue("result exceeds transfer budget")
	}
	b.remaining -= amount
	return nil
}

func decodeSandboxValueAt(
	rt *core.Runtime,
	encoded sandboxwire.Value,
	budget *sandboxDecodeBudget,
	depth int,
) (object.EmeraldValue, error) {
	if depth > maxSandboxValueDepth {
		return object.EmeraldValue{}, sandboxInvalidWorkerValue("result nesting is too deep")
	}
	if err := budget.consume(1); err != nil {
		return object.EmeraldValue{}, err
	}

	switch encoded.Type {
	case sandboxwire.ValueNil:
		return rt.NULL, nil
	case sandboxwire.ValueBool:
		return rt.NativeBoolToBooleanObject(encoded.Bool), nil
	case sandboxwire.ValueInt:
		return rt.NewInteger(encoded.Int), nil
	case sandboxwire.ValueFloat:
		return rt.NewFloat(math.Float64frombits(encoded.FloatBits)), nil
	case sandboxwire.ValueString:
		if err := budget.consume(len(encoded.String)); err != nil {
			return object.EmeraldValue{}, err
		}
		return rt.NewString(encoded.String), nil
	case sandboxwire.ValueSymbol:
		if err := budget.consume(len(encoded.Symbol)); err != nil {
			return object.EmeraldValue{}, err
		}
		return rt.NewSymbol(encoded.Symbol), nil
	case sandboxwire.ValueArray:
		if err := budget.consume(len(encoded.Array)); err != nil {
			return object.EmeraldValue{}, err
		}
		values := make([]object.EmeraldValue, len(encoded.Array))
		for i, item := range encoded.Array {
			value, err := decodeSandboxValueAt(rt, item, budget, depth+1)
			if err != nil {
				return object.EmeraldValue{}, err
			}
			values[i] = value
		}
		return rt.NewArray(values), nil
	case sandboxwire.ValueHash:
		if err := budget.consume(len(encoded.Hash)); err != nil {
			return object.EmeraldValue{}, err
		}
		hashValue := rt.NewHash()
		hash := hashValue.Heap.(*core.HashInstance)
		seenKeys := map[string]struct{}{}
		for _, entry := range encoded.Hash {
			key, err := decodeSandboxValueAt(rt, entry.Key, budget, depth+1)
			if err != nil {
				return object.EmeraldValue{}, err
			}
			keyHash := key.HashKey()
			if _, exists := seenKeys[keyHash]; exists {
				return object.EmeraldValue{}, sandboxInvalidWorkerValue("result hash has duplicate keys")
			}
			seenKeys[keyHash] = struct{}{}
			value, err := decodeSandboxValueAt(rt, entry.Value, budget, depth+1)
			if err != nil {
				return object.EmeraldValue{}, err
			}
			hash.Set(key, value)
		}
		return hashValue, nil
	default:
		return object.EmeraldValue{}, sandboxInvalidWorkerValue("result has an unknown type")
	}
}

func sandboxInvalidWorkerValue(reason string) error {
	return &SandboxWorkerError{Detail: "invalid worker result: " + reason}
}
