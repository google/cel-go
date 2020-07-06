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
// Decodes a base64-encoded string to an unencoded string value.
//
// This function will return an error if the input string is not base64-encoded.
//
//     base64.decode(<string>) -> <string>
//
// Examples:
//
//     base64.decode('aGVsbG8=')  // return 'hello'
//     base64.decode('aGVsbG8')   // error
//
// Base64.Encode
//
// Encodes a bytes or string input to a base64-encoded string.
//
//     base64.encode(<bytes>) -> <string>
//     base64.encode(<string>) -> <string>
//
// Examples:
//
//     base64.encode('hello')  // return 'aGVsbG8='
//     base64.encode(b'hello') // return 'aGVsbG8='
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
					decls.String)),
			decls.NewFunction("base64.encode",
				decls.NewOverload("base64_encode_string",
					[]*exprpb.Type{decls.String},
					decls.String),
				decls.NewOverload("base64_encode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("json.decode",
				decls.NewOverload("json_decode_string",
					[]*exprpb.Type{decls.String},
					decls.Dyn)),
			decls.NewFunction("json.encode",
				decls.NewOverload("json_encode_dyn",
					[]*exprpb.Type{decls.Dyn},
					decls.String)),
		),
	}
}

func (encoderLib) ProgramOptions() []cel.ProgramOption {
	wrappedBase64EncodeBytes := callInBytesOutStr(base64EncodeBytes)
	wrappedBase64EncodeString := callInStrOutStr(base64EncodeString)
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "base64.decode",
				Unary:    callInStrOutStr(base64DecodeString),
			},
			&functions.Overload{
				Operator: "base64_decode_string",
				Unary:    callInStrOutStr(base64DecodeString),
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
					return types.NoSuchOverloadErr()
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

func base64DecodeString(str string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func base64EncodeBytes(bytes []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func base64EncodeString(str string) (string, error) {
	return base64EncodeBytes([]byte(str))
}
