// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package proto provides functionality for handling protocol buffer messages.
// In particular, it provides marshaling and unmarshaling between a protobuf
// message and the binary wire format.
//
// See https://developers.google.com/protocol-buffers/docs/gotutorial for
// more information.
//
// Deprecated: Use the "google.golang.org/protobuf/proto" package instead.
package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	ProtoPackageIsVersion1 = true
	ProtoPackageIsVersion2 = true
	ProtoPackageIsVersion3 = true
	ProtoPackageIsVersion4 = true
)

// Message is a protocol buffer message.
type Message = protoiface.MessageV1

// Marshaler is implemented by messages that can marshal themselves.
// This interface is used by the following functions: Size, Marshal,
// Buffer.Marshal, and Buffer.EncodeMessage.
//
// Deprecated: Do not implement.
type Marshaler interface {
	// Marshal formats the encoded bytes of the message.
	// It should be deterministic and emit valid protobuf wire data.
	// The caller takes ownership of the returned buffer.
	Marshal() ([]byte, error)
}

// Unmarshaler is implemented by messages that can unmarshal themselves.
// This interface is used by the following functions: Unmarshal, UnmarshalMerge,
// Buffer.Unmarshal, Buffer.DecodeMessage, and Buffer.DecodeGroup.
//
// Deprecated: Do not implement.
type Unmarshaler interface {
	// Unmarshal parses the encoded bytes of the protobuf wire input.
	// The provided buffer is only valid for during method call.
	// It should not reset the receiver message.
	Unmarshal([]byte) error
}

// Merger is implemented by messages that can merge themselves.
// This interface is used by the following functions: Clone and Merge.
//
// Deprecated: Do not implement.
type Merger interface {
	// Merge merges the contents of src into the receiver message.
	// It clones all data structures in src such that it aliases no mutable
	// memory referenced by src.
	Merge(src Message)
}

// RequiredNotSetError is an error type returned when
// marshaling or unmarshaling a message with missing required fields.
type RequiredNotSetError struct {
	err error
}

func (e *RequiredNotSetError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "proto: required field not set"
}
func (e *RequiredNotSetError) RequiredNotSet() bool {
	return true
}

func checkRequiredNotSet(m proto.Message) error {
	if err := proto.IsInitialized(m); err != nil {
		return &RequiredNotSetError{err: err}
	}
	return nil
}

// Clone returns a deep copy of src.
func Clone(src Message) Message {
	srcMsg := protoimpl.X.MessageOf(src)
	if srcMsg == nil || !srcMsg.IsValid() {
		return src
	}

	dst := protoimpl.X.ProtoMessageV1Of(srcMsg.New().Interface())
	Merge(dst, src)
	return dst
}

// Merge merges src into dst, which must be messages of the same type.
//
// Populated scalar fields in src are copied to dst, while populated
// singular messages in src are merged into dst by recursively calling Merge.
// The elements of every list field in src is appended to the corresponded
// list fields in dst. The entries of every map field in src is copied into
// the corresponding map field in dst, possibly replacing existing entries.
// The unknown fields of src are appended to the unknown fields of dst.
func Merge(dst, src Message) {
	// TODO: Drop this type assertion if the aberrant wrapper in v2 calls this.
	if m, ok := dst.(Merger); ok {
		m.Merge(src)
		return
	}
	proto.Merge(
		protoimpl.X.ProtoMessageV2Of(dst),
		protoimpl.X.ProtoMessageV2Of(src),
	)
}

// Equal reports whether two messages are equal.
// If two messages marshal to the same bytes under deterministic serialization,
// then Equal is guaranteed to report true.
//
// Two messages are equal if they are the same protobuf message type,
// have the same set of populated known and extension field values,
// and the same set of unknown fields values.
//
// Scalar values are compared with the equivalent of the == operator in Go,
// except bytes values which are compared using bytes.Equal and
// floating point values which specially treat NaNs as equal.
// Message values are compared by recursively calling Equal.
// Lists are equal if each element value is also equal.
// Maps are equal if they have the same set of keys, where the pair of values
// for each key is also equal.
func Equal(x, y Message) bool {
	return proto.Equal(
		protoimpl.X.ProtoMessageV2Of(x),
		protoimpl.X.ProtoMessageV2Of(y),
	)
}

func isMessageSet(md protoreflect.MessageDescriptor) bool {
	ms, ok := md.(interface{ IsMessageSet() bool })
	return ok && ms.IsMessageSet()
}
