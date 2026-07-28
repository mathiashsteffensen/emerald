package sandboxwire

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"
)

const testMaxFrameBytes = 16 * 1024

func TestWriteReadRoundTrip(t *testing.T) {
	want := Response{
		Version: Version,
		Value: &Value{
			Type: ValueHash,
			Hash: []HashEntry{
				{Key: Value{Type: ValueSymbol, Symbol: "answer"}, Value: Value{Type: ValueInt, Int: 42}},
				{Key: Value{Type: ValueString, String: "items"}, Value: Value{
					Type: ValueArray,
					Array: []Value{
						{Type: ValueNil},
						{Type: ValueBool, Bool: true},
						{Type: ValueFloat, FloatBits: math.Float64bits(1.5)},
					},
				}},
			},
		},
	}

	var frame bytes.Buffer
	if err := Write(&frame, testMaxFrameBytes, want); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	var got Response
	if err := Read(&frame, testMaxFrameBytes, &got); err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestRequestRoundTrip(t *testing.T) {
	want := Request{
		Version:          Version,
		FileName:         "example.em",
		Content:          "21 * 2",
		DeadlineUnixNano: 123456789,
		MemoryLimitBytes: 64 * 1024 * 1024,
	}

	var frame bytes.Buffer
	if err := Write(&frame, testMaxFrameBytes, want); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	var got Request
	if err := Read(&frame, testMaxFrameBytes, &got); err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if got != want {
		t.Fatalf("round trip mismatch: got %#v, want %#v", got, want)
	}
}

func TestReadRejectsZeroLengthFrame(t *testing.T) {
	err := Read(bytes.NewReader([]byte{0, 0, 0, 0}), testMaxFrameBytes, &Request{})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("Read() error = %v, want ErrInvalidFrame", err)
	}
}

func TestReadRejectsOversizedFrame(t *testing.T) {
	frame := frameWithLength(uint32(testMaxFrameBytes + 1))
	err := Read(bytes.NewReader(frame), testMaxFrameBytes, &Request{})
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("Read() error = %v, want ErrFrameTooLarge", err)
	}
}

func TestReadRejectsTruncatedHeader(t *testing.T) {
	err := Read(bytes.NewReader([]byte{0, 0, 0}), testMaxFrameBytes, &Request{})
	if !errors.Is(err, ErrTruncatedFrame) {
		t.Fatalf("Read() error = %v, want ErrTruncatedFrame", err)
	}
}

func TestReadRejectsTruncatedPayload(t *testing.T) {
	frame := append(frameWithLength(5), '{', '}')
	err := Read(bytes.NewReader(frame), testMaxFrameBytes, &Request{})
	if !errors.Is(err, ErrTruncatedFrame) {
		t.Fatalf("Read() error = %v, want ErrTruncatedFrame", err)
	}
}

func TestWriteRejectsOversizedFrame(t *testing.T) {
	err := Write(&bytes.Buffer{}, 8, Request{Content: "this is too large"})
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("Write() error = %v, want ErrFrameTooLarge", err)
	}
}

func TestReadRejectsUnknownFields(t *testing.T) {
	payload := []byte(`{"version":1,"unexpected":true}`)
	frame := append(frameWithLength(uint32(len(payload))), payload...)
	err := Read(bytes.NewReader(frame), testMaxFrameBytes, &Request{})
	if err == nil {
		t.Fatal("Read() unexpectedly accepted an unknown field")
	}
}

func TestValueUnmarshalAcceptsOmittedActiveZeroValues(t *testing.T) {
	tests := []struct {
		name string
		json string
		want Value
	}{
		{name: "nil", json: `{"type":"nil"}`, want: Value{Type: ValueNil}},
		{name: "bool", json: `{"type":"bool"}`, want: Value{Type: ValueBool}},
		{name: "int", json: `{"type":"int"}`, want: Value{Type: ValueInt}},
		{name: "float", json: `{"type":"float"}`, want: Value{Type: ValueFloat}},
		{name: "string", json: `{"type":"string"}`, want: Value{Type: ValueString}},
		{name: "symbol", json: `{"type":"symbol"}`, want: Value{Type: ValueSymbol}},
		{name: "array", json: `{"type":"array"}`, want: Value{Type: ValueArray}},
		{name: "hash", json: `{"type":"hash"}`, want: Value{Type: ValueHash}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got Value
			if err := json.Unmarshal([]byte(test.json), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Unmarshal() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestValueUnmarshalAcceptsExplicitActiveZeroValues(t *testing.T) {
	tests := []struct {
		name string
		json string
		want Value
	}{
		{name: "bool", json: `{"type":"bool","bool":false}`, want: Value{Type: ValueBool}},
		{name: "int", json: `{"type":"int","int":0}`, want: Value{Type: ValueInt}},
		{name: "float", json: `{"type":"float","float_bits":0}`, want: Value{Type: ValueFloat}},
		{name: "string", json: `{"type":"string","string":""}`, want: Value{Type: ValueString}},
		{name: "symbol", json: `{"type":"symbol","symbol":""}`, want: Value{Type: ValueSymbol}},
		{name: "array", json: `{"type":"array","array":[]}`, want: Value{Type: ValueArray, Array: []Value{}}},
		{name: "hash", json: `{"type":"hash","hash":[]}`, want: Value{Type: ValueHash, Hash: []HashEntry{}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got Value
			if err := json.Unmarshal([]byte(test.json), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Unmarshal() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestValueUnmarshalRejectsMalformedTaggedValues(t *testing.T) {
	tests := map[string]string{
		"unknown field":            `{"type":"nil","unexpected":true}`,
		"incompatible field":       `{"type":"nil","array":[]}`,
		"unknown type":             `{"type":"object"}`,
		"missing type":             `{"string":"ok"}`,
		"null scalar":              `{"type":"bool","bool":null}`,
		"null collection":          `{"type":"array","array":null}`,
		"duplicate field":          `{"type":"int","type":"nil"}`,
		"nested incompatible":      `{"type":"array","array":[{"type":"nil","array":[]}]}`,
		"hash entry unknown field": `{"type":"hash","hash":[{"key":{"type":"int"},"value":{"type":"nil"},"extra":true}]}`,
	}

	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			var value Value
			err := json.Unmarshal([]byte(payload), &value)
			if !errors.Is(err, ErrInvalidValue) {
				t.Fatalf("Unmarshal() error = %v, want ErrInvalidValue", err)
			}
		})
	}
}

func TestReadRejectsMalformedTaggedValue(t *testing.T) {
	payload := []byte(`{"version":1,"value":{"type":"nil","array":[]}}`)
	frame := append(frameWithLength(uint32(len(payload))), payload...)

	err := Read(bytes.NewReader(frame), testMaxFrameBytes, &Response{})
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("Read() error = %v, want ErrInvalidValue", err)
	}
}

func frameWithLength(length uint32) []byte {
	frame := make([]byte, HeaderSize)
	binary.BigEndian.PutUint32(frame, length)
	return frame
}
