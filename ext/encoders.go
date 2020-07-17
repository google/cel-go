// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ext

import (
	"encoding/base64"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Encoders returns a cel.EnvOption to configure extended functions for string, byte, and object
// encodings.
//
// Base64.Decode
//
// Decodes bytes or a base64-encoded string to a sequence of bytes.
//
// This function will return an error if the UTF-8 string input is not base64-encoded.
//
//     base64.decode(<bytes>)  -> <bytes>
//     base64.decode(<string>) -> <bytes>
//
// Examples:
//
//     base64.decode(b'aGVsbG8=') // return b'hello'
//     base64.decode('aGVsbG8=')  // return b'hello'
//     base64.decode('aGVsbG8')   // error
//
// Base64.Encode
//
// Encodes a bytes or string input to base64-encoded bytes.
//
//     base64.encode(<bytes>)  -> <bytes>
//     base64.encode(<string>) -> <bytes>
//
// Examples:
//
//     base64.encode(b'hello') // return b'aGVsbG8='
//     base64.encode('hello')  // return b'aGVsbG8='
func Encoders() cel.EnvOption {
	return cel.Lib(encoderLib{})
}

type encoderLib struct{}

func (encoderLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewFunction("base64.decode",
				decls.NewOverload("base64_decode_string",
					[]*exprpb.Type{decls.String},
					decls.Bytes),
				decls.NewOverload("base64_decode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.Bytes)),
			decls.NewFunction("base64.encode",
				decls.NewOverload("base64_encode_string",
					[]*exprpb.Type{decls.String},
					decls.Bytes),
				decls.NewOverload("base64_encode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.Bytes)),
		),
	}
}

func (encoderLib) ProgramOptions() []cel.ProgramOption {
	wrappedBase64EncodeBytes := callInBytesOutBytes(base64EncodeBytes)
	wrappedBase64EncodeString := callInStrOutBytes(base64EncodeString)
	wrappedBase64DecodeBytes := callInBytesOutBytes(base64DecodeBytes)
	wrappedBase64DecodeString := callInStrOutBytes(base64DecodeString)
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "base64.decode",
				Unary: func(val ref.Val) ref.Val {
					switch val.(type) {
					case types.Bytes:
						return wrappedBase64DecodeBytes(val)
					case types.String:
						return wrappedBase64DecodeString(val)
					}
					return types.MaybeNoSuchOverloadErr(val)
				},
			},
			&functions.Overload{
				Operator: "base64_decode_bytes",
				Unary:    wrappedBase64DecodeBytes,
			},
			&functions.Overload{
				Operator: "base64_decode_string",
				Unary:    wrappedBase64DecodeString,
			},
			&functions.Overload{
				Operator: "base64.encode",
				Unary: func(val ref.Val) ref.Val {
					switch val.(type) {
					case types.Bytes:
						return wrappedBase64EncodeBytes(val)
					case types.String:
						return wrappedBase64EncodeString(val)
					}
					return types.MaybeNoSuchOverloadErr(val)
				},
			},
			&functions.Overload{
				Operator: "base64_encode_bytes",
				Unary:    wrappedBase64EncodeBytes,
			},
			&functions.Overload{
				Operator: "base64_encode_string",
				Unary:    wrappedBase64EncodeString,
			},
		),
	}
}

func base64DecodeBytes(bytes []byte) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
	n, err := base64.StdEncoding.Decode(buf, bytes)
	return buf[:n], err
}

func base64DecodeString(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func base64EncodeBytes(bytes []byte) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(buf, bytes)
	return buf, nil
}

func base64EncodeString(str string) ([]byte, error) {
	return base64EncodeBytes([]byte(str))
}
