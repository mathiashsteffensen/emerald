// Package sandboxwire defines the private protocol between Sandbox and its
// worker process.
package sandboxwire

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
)

// Version is the current wire protocol version.
const Version uint32 = 1

// HeaderSize is the number of bytes in a frame length prefix.
const HeaderSize = 4

var (
	// ErrInvalidFrame is returned for a zero-length frame or invalid frame limit.
	ErrInvalidFrame = errors.New("sandboxwire: invalid frame")
	// ErrFrameTooLarge is returned when a frame exceeds the configured limit.
	ErrFrameTooLarge = errors.New("sandboxwire: frame too large")
	// ErrTruncatedFrame is returned when a frame header or payload ends early.
	ErrTruncatedFrame = errors.New("sandboxwire: truncated frame")
	// ErrInvalidValue is returned when a tagged value is malformed.
	ErrInvalidValue = errors.New("sandboxwire: invalid value")
)

// Request asks a worker to evaluate Content as FileName. DeadlineUnixNano and
// MemoryLimitBytes are supplied by the parent and apply to that evaluation.
type Request struct {
	Version          uint32 `json:"version"`
	FileName         string `json:"file_name"`
	Content          string `json:"content"`
	DeadlineUnixNano int64  `json:"deadline_unix_nano"`
	MemoryLimitBytes int64  `json:"memory_limit_bytes"`
}

// Phase identifies the stage of evaluation reported by the worker.
type Phase string

const (
	PhaseCompile Phase = "compile"
	PhaseExecute Phase = "execute"
)

// Response is emitted by a worker. A phase-only response reports progress; a
// terminal response has either Value or Error set.
type Response struct {
	Version uint32 `json:"version"`
	Phase   Phase  `json:"phase,omitempty"`
	Value   *Value `json:"value,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Error carries an error reported by the worker. Its interpretation is left to
// the parent-side sandbox API.
type Error struct {
	Kind       string `json:"kind,omitempty"`
	ClassName  string `json:"class_name,omitempty"`
	Message    string `json:"message"`
	Phase      Phase  `json:"phase,omitempty"`
	Stacktrace []byte `json:"stacktrace,omitempty"`
}

// ValueType tags the active field in a Value.
type ValueType string

const (
	ValueNil    ValueType = "nil"
	ValueBool   ValueType = "bool"
	ValueInt    ValueType = "int"
	ValueFloat  ValueType = "float"
	ValueString ValueType = "string"
	ValueSymbol ValueType = "symbol"
	ValueArray  ValueType = "array"
	ValueHash   ValueType = "hash"
)

// Value is a transferable Emerald result. Type identifies which of the value
// fields is meaningful.
type Value struct {
	Type      ValueType   `json:"type"`
	Bool      bool        `json:"bool,omitempty"`
	Int       int64       `json:"int,omitempty"`
	FloatBits uint64      `json:"float_bits,omitempty"`
	String    string      `json:"string,omitempty"`
	Symbol    string      `json:"symbol,omitempty"`
	Array     []Value     `json:"array,omitempty"`
	Hash      []HashEntry `json:"hash,omitempty"`
}

// HashEntry represents one hash key/value pair. Hashes are represented as
// entries so non-string keys retain their type.
type HashEntry struct {
	Key   Value `json:"key"`
	Value Value `json:"value"`
}

// UnmarshalJSON accepts only the fields compatible with Value.Type. Active
// fields may be omitted because zero values are omitted by Value's JSON tags.
func (v *Value) UnmarshalJSON(data []byte) error {
	fields, err := decodeObject(data)
	if err != nil {
		return err
	}

	rawType, ok := fields["type"]
	if !ok {
		return fmt.Errorf("%w: missing type", ErrInvalidValue)
	}

	var valueType ValueType
	if err := decodeNonNull(rawType, &valueType); err != nil {
		return fmt.Errorf("%w: type: %v", ErrInvalidValue, err)
	}

	decoded := Value{Type: valueType}
	switch valueType {
	case ValueNil:
		if err := allowOnly(fields, "type"); err != nil {
			return err
		}
	case ValueBool:
		if err := allowOnly(fields, "type", "bool"); err != nil {
			return err
		}
		if raw, ok := fields["bool"]; ok {
			if err := decodeNonNull(raw, &decoded.Bool); err != nil {
				return fmt.Errorf("%w: bool: %v", ErrInvalidValue, err)
			}
		}
	case ValueInt:
		if err := allowOnly(fields, "type", "int"); err != nil {
			return err
		}
		if raw, ok := fields["int"]; ok {
			if err := decodeNonNull(raw, &decoded.Int); err != nil {
				return fmt.Errorf("%w: int: %v", ErrInvalidValue, err)
			}
		}
	case ValueFloat:
		if err := allowOnly(fields, "type", "float_bits"); err != nil {
			return err
		}
		if raw, ok := fields["float_bits"]; ok {
			if err := decodeNonNull(raw, &decoded.FloatBits); err != nil {
				return fmt.Errorf("%w: float_bits: %v", ErrInvalidValue, err)
			}
		}
	case ValueString:
		if err := allowOnly(fields, "type", "string"); err != nil {
			return err
		}
		if raw, ok := fields["string"]; ok {
			if err := decodeNonNull(raw, &decoded.String); err != nil {
				return fmt.Errorf("%w: string: %v", ErrInvalidValue, err)
			}
		}
	case ValueSymbol:
		if err := allowOnly(fields, "type", "symbol"); err != nil {
			return err
		}
		if raw, ok := fields["symbol"]; ok {
			if err := decodeNonNull(raw, &decoded.Symbol); err != nil {
				return fmt.Errorf("%w: symbol: %v", ErrInvalidValue, err)
			}
		}
	case ValueArray:
		if err := allowOnly(fields, "type", "array"); err != nil {
			return err
		}
		if raw, ok := fields["array"]; ok {
			if err := decodeNonNull(raw, &decoded.Array); err != nil {
				return fmt.Errorf("%w: array: %v", ErrInvalidValue, err)
			}
		}
	case ValueHash:
		if err := allowOnly(fields, "type", "hash"); err != nil {
			return err
		}
		if raw, ok := fields["hash"]; ok {
			if err := decodeNonNull(raw, &decoded.Hash); err != nil {
				return fmt.Errorf("%w: hash: %v", ErrInvalidValue, err)
			}
		}
	default:
		return fmt.Errorf("%w: unknown type %q", ErrInvalidValue, valueType)
	}

	*v = decoded
	return nil
}

// UnmarshalJSON accepts a hash entry with exactly one key and one value.
func (e *HashEntry) UnmarshalJSON(data []byte) error {
	fields, err := decodeObject(data)
	if err != nil {
		return err
	}
	if err := allowOnly(fields, "key", "value"); err != nil {
		return err
	}

	rawKey, ok := fields["key"]
	if !ok {
		return fmt.Errorf("%w: hash entry missing key", ErrInvalidValue)
	}
	rawValue, ok := fields["value"]
	if !ok {
		return fmt.Errorf("%w: hash entry missing value", ErrInvalidValue)
	}

	var decoded HashEntry
	if err := json.Unmarshal(rawKey, &decoded.Key); err != nil {
		return fmt.Errorf("%w: hash key: %v", ErrInvalidValue, err)
	}
	if err := json.Unmarshal(rawValue, &decoded.Value); err != nil {
		return fmt.Errorf("%w: hash value: %v", ErrInvalidValue, err)
	}

	*e = decoded
	return nil
}

// Write serializes message as JSON and writes it as a big-endian uint32 length
// prefix followed by its payload. The serialized payload must not exceed
// maxFrameBytes.
func Write(w io.Writer, maxFrameBytes int, message any) error {
	if err := validateMaxFrameBytes(maxFrameBytes); err != nil {
		return err
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("sandboxwire: marshal frame: %w", err)
	}
	if len(payload) == 0 {
		return fmt.Errorf("%w: zero-length payload", ErrInvalidFrame)
	}
	if len(payload) > maxFrameBytes || uint64(len(payload)) > math.MaxUint32 {
		return fmt.Errorf("%w: %d bytes exceeds %d", ErrFrameTooLarge, len(payload), maxFrameBytes)
	}

	var header [HeaderSize]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if err := writeAll(w, header[:]); err != nil {
		return fmt.Errorf("sandboxwire: write frame header: %w", err)
	}
	if err := writeAll(w, payload); err != nil {
		return fmt.Errorf("sandboxwire: write frame payload: %w", err)
	}

	return nil
}

// Read reads one big-endian uint32 length-prefixed JSON frame into message.
// It rejects zero-length, oversized, and truncated frames before decoding.
func Read(r io.Reader, maxFrameBytes int, message any) error {
	if err := validateMaxFrameBytes(maxFrameBytes); err != nil {
		return err
	}

	var header [HeaderSize]byte
	n, err := io.ReadFull(r, header[:])
	if err != nil {
		if err == io.EOF && n == 0 {
			return io.EOF
		}
		return fmt.Errorf("%w: frame header: %v", ErrTruncatedFrame, err)
	}

	frameBytes := binary.BigEndian.Uint32(header[:])
	if frameBytes == 0 {
		return fmt.Errorf("%w: zero-length payload", ErrInvalidFrame)
	}
	if uint64(frameBytes) > uint64(maxFrameBytes) {
		return fmt.Errorf("%w: %d bytes exceeds %d", ErrFrameTooLarge, frameBytes, maxFrameBytes)
	}

	payload := make([]byte, int(frameBytes))
	if _, err := io.ReadFull(r, payload); err != nil {
		return fmt.Errorf("%w: frame payload: %v", ErrTruncatedFrame, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(message); err != nil {
		return fmt.Errorf("sandboxwire: decode frame: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("sandboxwire: decode frame: multiple JSON values")
		}
		return fmt.Errorf("sandboxwire: decode frame: %w", err)
	}

	return nil
}

func validateMaxFrameBytes(maxFrameBytes int) error {
	if maxFrameBytes <= 0 || uint64(maxFrameBytes) > math.MaxUint32 {
		return fmt.Errorf("%w: max frame bytes must be between 1 and %d", ErrInvalidFrame, uint64(math.MaxUint32))
	}

	return nil
}

func decodeObject(data []byte) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: decode object: %v", ErrInvalidValue, err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '{' {
		return nil, fmt.Errorf("%w: expected object", ErrInvalidValue)
	}

	fields := make(map[string]json.RawMessage)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("%w: read field name: %v", ErrInvalidValue, err)
		}
		name, ok := token.(string)
		if !ok {
			return nil, fmt.Errorf("%w: field name is not a string", ErrInvalidValue)
		}
		if _, exists := fields[name]; exists {
			return nil, fmt.Errorf("%w: duplicate field %q", ErrInvalidValue, name)
		}

		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return nil, fmt.Errorf("%w: decode field %q: %v", ErrInvalidValue, name, err)
		}
		fields[name] = raw
	}

	token, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: end object: %v", ErrInvalidValue, err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '}' {
		return nil, fmt.Errorf("%w: expected end object", ErrInvalidValue)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("%w: multiple JSON values", ErrInvalidValue)
		}
		return nil, fmt.Errorf("%w: trailing JSON: %v", ErrInvalidValue, err)
	}

	return fields, nil
}

func allowOnly(fields map[string]json.RawMessage, names ...string) error {
	allowed := make(map[string]struct{}, len(names))
	for _, name := range names {
		allowed[name] = struct{}{}
	}
	for name := range fields {
		if _, ok := allowed[name]; !ok {
			return fmt.Errorf("%w: field %q is incompatible with value type", ErrInvalidValue, name)
		}
	}

	return nil
}

func decodeNonNull(raw json.RawMessage, value any) error {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return errors.New("must not be null")
	}
	if err := json.Unmarshal(raw, value); err != nil {
		return err
	}

	return nil
}

func writeAll(w io.Writer, payload []byte) error {
	for len(payload) > 0 {
		n, err := w.Write(payload)
		if n > 0 {
			payload = payload[n:]
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}

	return nil
}
