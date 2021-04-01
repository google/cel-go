// Copyright 2021 Google LLC
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

package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// asciiCodePointBuffer is an implementation for an array of code points that contain code points
// only from the ASCII character set.
type asciiCodePointBuffer struct {
	arr []byte
	pos int
	src string
}

// Consume implements (antlr.CharStream).Consume.
func (a *asciiCodePointBuffer) Consume() {
	if a.pos >= len(a.arr) {
		panic("cannot consume EOF")
	}
	a.pos++
}

// LA implements (antlr.CharStream).LA.
func (a *asciiCodePointBuffer) LA(offset int) int {
	if offset == 0 {
		return 0
	}
	if offset < 0 {
		offset++
	}
	pos := a.pos + offset - 1
	if pos < 0 || pos >= len(a.arr) {
		return antlr.TokenEOF
	}
	return int(uint(a.arr[pos]))
}

// LT mimics (*antlr.InputStream).LT.
func (a *asciiCodePointBuffer) LT(offset int) int {
	return a.LA(offset)
}

// Mark implements (antlr.CharStream).Mark.
func (a *asciiCodePointBuffer) Mark() int {
	return -1
}

// Release implements (antlr.CharStream).Release.
func (a *asciiCodePointBuffer) Release(marker int) {}

// Index implements (antlr.CharStream).Index.
func (a *asciiCodePointBuffer) Index() int {
	return a.pos
}

// Seek implements (antlr.CharStream).Seek.
func (a *asciiCodePointBuffer) Seek(index int) {
	if index <= a.pos {
		a.pos = index
		return
	}
	if index < len(a.arr) {
		a.pos = index
	} else {
		a.pos = len(a.arr)
	}
}

// Size implements (antlr.CharStream).Size.
func (a *asciiCodePointBuffer) Size() int {
	return len(a.arr)
}

// GetSourceName implements (antlr.CharStream).GetSourceName.
func (a *asciiCodePointBuffer) GetSourceName() string {
	return a.src
}

// GetText implements (antlr.CharStream).GetText.
func (a *asciiCodePointBuffer) GetText(start, stop int) string {
	if stop >= len(a.arr) {
		stop = len(a.arr) - 1
	}
	if start >= len(a.arr) {
		return ""
	}
	return string(a.arr[start : stop+1])
}

// GetTextFromTokens implements (antlr.CharStream).GetTextFromTokens.
func (a *asciiCodePointBuffer) GetTextFromTokens(start, stop antlr.Token) string {
	if start != nil && stop != nil {
		return a.GetText(start.GetTokenIndex(), stop.GetTokenIndex())
	}
	return ""
}

// GetTextFromInterval implements (antlr.CharStream).GetTextFromInterval.
func (a *asciiCodePointBuffer) GetTextFromInterval(i *antlr.Interval) string {
	return a.GetText(i.Start, i.Stop)
}

// String mimics (*antlr.InputStream).String.
func (a *asciiCodePointBuffer) String() string {
	return string(a.arr)
}

var _ antlr.CharStream = &asciiCodePointBuffer{}

// basicCodePointBuffer is an implementation for an array of code points that contain code points
// from both the Latin-1 character set and Basic Multilingual Plane.
type basicCodePointBuffer struct {
	arr []uint16
	pos int
	src string
}

// Consume implements (antlr.CharStream).Consume.
func (b *basicCodePointBuffer) Consume() {
	if b.pos >= len(b.arr) {
		panic("cannot consume EOF")
	}
	b.pos++
}

// LA implements (antlr.CharStream).LA.
func (b *basicCodePointBuffer) LA(offset int) int {
	if offset == 0 {
		return 0
	}
	if offset < 0 {
		offset++
	}
	pos := b.pos + offset - 1
	if pos < 0 || pos >= len(b.arr) {
		return antlr.TokenEOF
	}
	return int(uint(b.arr[pos]))
}

// LT mimics (*antlr.InputStream).LT.
func (b *basicCodePointBuffer) LT(offset int) int {
	return b.LA(offset)
}

// Mark implements (antlr.CharStream).Mark.
func (b *basicCodePointBuffer) Mark() int {
	return -1
}

// Release implements (antlr.CharStream).Release.
func (b *basicCodePointBuffer) Release(marker int) {}

// Index implements (antlr.CharStream).Index.
func (b *basicCodePointBuffer) Index() int {
	return b.pos
}

// Seek implements (antlr.CharStream).Seek.
func (b *basicCodePointBuffer) Seek(index int) {
	if index <= b.pos {
		b.pos = index
		return
	}
	if index < len(b.arr) {
		b.pos = index
	} else {
		b.pos = len(b.arr)
	}
}

// Size implements (antlr.CharStream).Size.
func (b *basicCodePointBuffer) Size() int {
	return len(b.arr)
}

// GetSourceName implements (antlr.CharStream).GetSourceName.
func (b *basicCodePointBuffer) GetSourceName() string {
	return b.src
}

// GetText implements (antlr.CharStream).GetText.
func (b *basicCodePointBuffer) GetText(start, stop int) string {
	if stop >= len(b.arr) {
		stop = len(b.arr) - 1
	}
	if start >= len(b.arr) {
		return ""
	}
	var str strings.Builder
	for i := start; i <= stop; i++ {
		str.WriteRune(rune(uint32(b.arr[i])))
	}
	return str.String()
}

// GetTextFromTokens implements (antlr.CharStream).GetTextFromTokens.
func (b *basicCodePointBuffer) GetTextFromTokens(start, stop antlr.Token) string {
	if start != nil && stop != nil {
		return b.GetText(start.GetTokenIndex(), stop.GetTokenIndex())
	}
	return ""
}

// GetTextFromInterval implements (antlr.CharStream).GetTextFromInterval.
func (b *basicCodePointBuffer) GetTextFromInterval(i *antlr.Interval) string {
	return b.GetText(i.Start, i.Stop)
}

// String mimics (*antlr.InputStream).String.
func (b *basicCodePointBuffer) String() string {
	var str strings.Builder
	for _, v := range b.arr {
		str.WriteRune(rune(uint32(v)))
	}
	return str.String()
}

var _ antlr.CharStream = &basicCodePointBuffer{}

// supplementalCodePointBuffer is an implementation for an array of code points that contain code
// points from the Latin-1 character set, Basic Multilingual Plane, or the Supplemental Multilingual
// Plane.
type supplementalCodePointBuffer struct {
	arr []rune
	pos int
	src string
}

// Consume implements (antlr.CharStream).Consume.
func (s *supplementalCodePointBuffer) Consume() {
	if s.pos >= len(s.arr) {
		panic("cannot consume EOF")
	}
	s.pos++
}

// LA implements (antlr.CharStream).LA.
func (s *supplementalCodePointBuffer) LA(offset int) int {
	if offset == 0 {
		return 0
	}
	if offset < 0 {
		offset++
	}
	pos := s.pos + offset - 1
	if pos < 0 || pos >= len(s.arr) {
		return antlr.TokenEOF
	}
	return int(uint(s.arr[pos]))
}

// LT mimics (*antlr.InputStream).LT.
func (s *supplementalCodePointBuffer) LT(offset int) int {
	return s.LA(offset)
}

// Mark implements (antlr.CharStream).Mark.
func (s *supplementalCodePointBuffer) Mark() int {
	return -1
}

// Release implements (antlr.CharStream).Release.
func (s *supplementalCodePointBuffer) Release(marker int) {}

// Index implements (antlr.CharStream).Index.
func (s *supplementalCodePointBuffer) Index() int {
	return s.pos
}

// Seek implements (antlr.CharStream).Seek.
func (s *supplementalCodePointBuffer) Seek(index int) {
	if index <= s.pos {
		s.pos = index
		return
	}
	if index < len(s.arr) {
		s.pos = index
	} else {
		s.pos = len(s.arr)
	}
}

// Size implements (antlr.CharStream).Size.
func (s *supplementalCodePointBuffer) Size() int {
	return len(s.arr)
}

// GetSourceName implements (antlr.CharStream).GetSourceName.
func (s *supplementalCodePointBuffer) GetSourceName() string {
	return s.src
}

// GetText implements (antlr.CharStream).GetText.
func (s *supplementalCodePointBuffer) GetText(start, stop int) string {
	if stop >= len(s.arr) {
		stop = len(s.arr) - 1
	}
	if start >= len(s.arr) {
		return ""
	}
	return string(s.arr[start : stop+1])
}

// GetTextFromTokens implements (antlr.CharStream).GetTextFromTokens.
func (s *supplementalCodePointBuffer) GetTextFromTokens(start, stop antlr.Token) string {
	if start != nil && stop != nil {
		return s.GetText(start.GetTokenIndex(), stop.GetTokenIndex())
	}
	return ""
}

// GetTextFromInterval implements (antlr.CharStream).GetTextFromInterval.
func (s *supplementalCodePointBuffer) GetTextFromInterval(i *antlr.Interval) string {
	return s.GetText(i.Start, i.Stop)
}

// String mimics (*antlr.InputStream).String.
func (s *supplementalCodePointBuffer) String() string {
	return string(s.arr)
}

var _ antlr.CharStream = &supplementalCodePointBuffer{}

// newCodePointBuffer returns an efficient implementation of antlr.CharStream for the given text
// based on the ranges of the encoded code points contained within.
//
// Code points are represented as an array of byte, uint16, or rune. This approach ensures that
// each index represents a code point by itself without needing to use an array of rune. At first
// we assume all code points are less than or equal to '\u007f'. If this holds true, the
// underlying storage is a byte array containing only ASCII characters. If we encountered a code
// point above this range but less than or equal to '\uffff' we allocate a uint16 array, copy the
// elements of previous byte array to the uint16 array, and continue. If this holds true, the
// underlying storage is a uint16 array containing only Unicode characters in the Basic Multilingual
// Plane. If we encounter a code point above '\uffff' we allocate an rune array, copy the previous
// elements of the byte or uint16 array, and continue. The underlying storage is an rune array
// containing any Unicode character.
//
// TODO: upstream this to ANTLRv4 Go runtime?
func newCodePointBuffer(text, desc string) antlr.CharStream {
	var (
		idx   = 0
		buf8  = make([]byte, 0, len(text))
		buf16 []uint16
		buf32 []rune
	)
	for idx < len(text) {
		r, s := utf8.DecodeRuneInString(text[idx:])
		idx += s
		if r < utf8.RuneSelf {
			buf8 = append(buf8, byte(r))
			continue
		}
		if r <= 0xffff {
			buf16 = make([]uint16, len(buf8), len(text))
			for i, v := range buf8 {
				buf16[i] = uint16(v)
			}
			buf8 = nil
			buf16 = append(buf16, uint16(r))
			goto copy16
		}
		buf32 = make([]rune, len(buf8), len(text))
		for i, v := range buf8 {
			buf32[i] = rune(uint32(v))
		}
		buf8 = nil
		buf32 = append(buf32, r)
		goto copy32
	}
	return &asciiCodePointBuffer{
		arr: buf8,
		src: desc,
	}
copy16:
	for idx < len(text) {
		r, s := utf8.DecodeRuneInString(text[idx:])
		idx += s
		if r <= 0xffff {
			buf16 = append(buf16, uint16(r))
			continue
		}
		buf32 = make([]rune, len(buf16), len(text))
		for i, v := range buf16 {
			buf32[i] = rune(uint32(v))
		}
		buf16 = nil
		buf32 = append(buf32, r)
		goto copy32
	}
	return &basicCodePointBuffer{
		arr: buf16,
		src: desc,
	}
copy32:
	for idx < len(text) {
		r, s := utf8.DecodeRuneInString(text[idx:])
		idx += s
		buf32 = append(buf32, r)
	}
	return &supplementalCodePointBuffer{
		arr: buf32,
		src: desc,
	}
}
