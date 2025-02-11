// 该包实现了字节切片的操作，在golang中字节切片等价于字符串string，它的一些操作和字符串相似。
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bytes implements functions for the manipulation of byte slices.
// It is analogous to the facilities of the strings package.

package bytes

import (
	"internal/bytealg"
	"unicode"
	"unicode/utf8"
)

// Equal reports whether a and b
// are the same length and contain the same bytes.
// An nil argument is equivalent to an empty slice.
// Compare whether two byte arrays are equal
func Equal(a, b []byte) bool {
	// Neither cmd/compile nor gccgo allocates for these string conversions.
	// 强制转化成为 string 类型进行比较
	return string(a) == string(b)
}

// Compare returns an integer comparing two byte slices lexicographically.
// 相比较两个字符切片的字典序
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b. 返回值的区别
// A nil argument is equivalent to an empty slice. 空的参数是等同于一个空的切片的
func Compare(a, b []byte) int {
	return bytealg.Compare(a, b)
}

// explode splits s into a slice of UTF-8 sequences, one per Unicode code point (still slices of bytes),
// up to a maximum of n byte slices. Invalid UTF-8 sequences are chopped into individual bytes.
// explode 将 s 分割成 UTF-8 序列的切片，每个 Unicode 代码点一个切片（仍然是字节的切片），最多 n 个字节切片。
// 无效的 UTF-8 序列将被切割成单个字节。
func explode(s []byte, n int) [][]byte {
	/*
		该函数的作用是将一个字节切片 s 分割成多个子切片，每个子切片包含一个 UTF-8 编码的 Unicode 代码点。函数的参数 n 指定了返回的子切片的最大数量。
	*/
	if n <= 0 {
		n = len(s)
	}
	a := make([][]byte, n)
	var size int
	na := 0
	for len(s) > 0 {
		if na+1 >= n {
			a[na] = s
			na++
			break
		}
		_, size = utf8.DecodeRune(s)
		a[na] = s[0:size:size]
		s = s[size:]
		na++
	}
	return a[0:na]
}

// Count counts the number of non-overlapping instances of sep in s.
// 统计给定的字节切片 s 中，非重叠的子切片 sep 出现的次数
// If sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.
func Count(s, sep []byte) int {
	// special case
	if len(sep) == 0 {
		return utf8.RuneCount(s) + 1
	}
	if len(sep) == 1 {
		return bytealg.Count(s, sep[0])
	}
	n := 0
	for {
		i := Index(s, sep)
		if i == -1 {
			return n
		}
		n++
		s = s[i+len(sep):]
	}
}

// Contains reports whether subslice is within b.
// 判断这个切片是否是另一个切片的字切片
func Contains(b, subslice []byte) bool {
	return Index(b, subslice) != -1
}

// ContainsAny reports whether any of the UTF-8-encoded code points in chars are within b.
// 第一个方法是下面这几个方法的父类方法，本质上都是上面的一个方法实现的
// Array 数组类型的
func ContainsAny(b []byte, chars string) bool {
	return IndexAny(b, chars) >= 0
}

// ContainsRune reports whether the rune is contained in the UTF-8-encoded byte slice b.
// Rune 类型
func ContainsRune(b []byte, r rune) bool {
	return IndexRune(b, r) >= 0
}

// IndexByte returns the index of the first instance of c in b, or -1 if c is not present in b.
// 下面的方法都是在索引字节数组的基础上进行实现的
func IndexByte(b []byte, c byte) int {
	return bytealg.IndexByte(b, c)
}

// 查找切片出现的第一次的位置
func indexBytePortable(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.
// 这段代码实现了 LastIndex 函数，目的是返回字节切片 s 中最后一次出现 sep 的位置索引。如果找不到 sep，则返回 -1。
func LastIndex(s, sep []byte) int {
	n := len(sep)
	switch {
	case n == 0:
		return len(s)
	case n == 1:
		return LastIndexByte(s, sep[0])
	case n == len(s):
		if Equal(s, sep) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}
	// Rabin-Karp search from the end of the string
	// Rabin-Karp 算法的应用
	hashss, pow := bytealg.HashStrRevBytes(sep)
	last := len(s) - n
	// 通过 Rabin-Karp 算法进行字符串匹配。首先计算 sep 的哈希值 hashss 和 pow（用于滑动窗口）。
	var h uint32
	for i := len(s) - 1; i >= last; i-- {
		h = h*bytealg.PrimeRK + uint32(s[i])
	}
	if h == hashss && Equal(s[last:], sep) {
		return last
	}
	for i := last - 1; i >= 0; i-- {
		h *= bytealg.PrimeRK
		h += uint32(s[i])
		h -= pow * uint32(s[i+n])
		if h == hashss && Equal(s[i:i+n], sep) {
			return i
		}
	}
	return -1
}

// 倒序寻找字符出现的位置
// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func LastIndexByte(s []byte, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// IndexRune interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index of the first occurrence in s of the given rune.
// It returns -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the first instance of any
// invalid UTF-8 byte sequence.
/*
这段代码实现了 IndexRune 函数，目的是在字节切片 s 中查找给定的 Unicode 字符 r（以 UTF-8 编码的方式），返回第一次出现该字符的字节索引。
如果没有找到该字符，返回 -1。如果 r 是无效的 UTF-8 字符，或者 r 是 utf8.RuneError，则根据不同情况处理。
*/
func IndexRune(s []byte, r rune) int {
	switch {
	case 0 <= r && r < utf8.RuneSelf:
		return IndexByte(s, byte(r))
	case r == utf8.RuneError:
		for i := 0; i < len(s); {
			r1, n := utf8.DecodeRune(s[i:])
			if r1 == utf8.RuneError {
				return i
			}
			i += n
		}
		return -1
	case !utf8.ValidRune(r):
		return -1
	default:
		var b [utf8.UTFMax]byte
		n := utf8.EncodeRune(b[:], r)
		return Index(s, b[:n])
	}
}

// IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points.
// It returns the byte index of the first occurrence in s of any of the Unicode
// code points in chars. It returns -1 if chars is empty or if there is no code
// point in common.
func IndexAny(s []byte, chars string) int {
	if chars == "" {
		// Avoid scanning all of s.
		return -1
	}
	if len(s) == 1 {
		r := rune(s[0])
		if r >= utf8.RuneSelf {
			// search utf8.RuneError.
			for _, r = range chars {
				if r == utf8.RuneError {
					return 0
				}
			}
			return -1
		}
		// Rabin-Karp 哈希算法：bytealg 采用 Rabin-Karp 算法（基于滑动窗口和哈希的匹配算法）来优化字节序列的查找。
		// 相比传统的逐个字符比较，哈希算法在长字符串查找时可以显著提高性能。
		if bytealg.IndexByteString(chars, s[0]) >= 0 {
			return 0
		}
		return -1
	}
	if len(chars) == 1 {
		r := rune(chars[0])
		if r >= utf8.RuneSelf {
			r = utf8.RuneError
		}
		return IndexRune(s, r)
	}
	if len(s) > 8 {
		if as, isASCII := makeASCIISet(chars); isASCII {
			for i, c := range s {
				if as.contains(c) {
					return i
				}
			}
			return -1
		}
	}
	var width int
	for i := 0; i < len(s); i += width {
		r := rune(s[i])
		if r < utf8.RuneSelf {
			if bytealg.IndexByteString(chars, s[i]) >= 0 {
				return i
			}
			width = 1
			continue
		}
		r, width = utf8.DecodeRune(s[i:])
		if r != utf8.RuneError {
			// r is 2 to 4 bytes
			if len(chars) == width {
				if chars == string(r) {
					return i
				}
				continue
			}
			// Use bytealg.IndexString for performance if available.
			if bytealg.MaxLen >= width {
				if bytealg.IndexString(chars, string(r)) >= 0 {
					return i
				}
				continue
			}
		}
		for _, ch := range chars {
			if r == ch {
				return i
			}
		}
	}
	return -1
}

// LastIndexAny interprets s as a sequence of UTF-8-encoded Unicode code
// points. It returns the byte index of the last occurrence in s of any of
// the Unicode code points in chars. It returns -1 if chars is empty or if
// there is no code point in common.
func LastIndexAny(s []byte, chars string) int {
	// 为了提高性能优先进特殊值的判断
	if chars == "" {
		// Avoid scanning all of s.
		return -1
	}
	if len(s) > 8 {
		if as, isASCII := makeASCIISet(chars); isASCII {
			for i := len(s) - 1; i >= 0; i-- {
				if as.contains(s[i]) {
					return i
				}
			}
			return -1
		}
	}
	if len(s) == 1 {
		r := rune(s[0])
		if r >= utf8.RuneSelf {
			for _, r = range chars {
				if r == utf8.RuneError {
					return 0
				}
			}
			return -1
		}
		if bytealg.IndexByteString(chars, s[0]) >= 0 {
			return 0
		}
		return -1
	}
	if len(chars) == 1 {
		cr := rune(chars[0])
		if cr >= utf8.RuneSelf {
			cr = utf8.RuneError
		}
		for i := len(s); i > 0; {
			r, size := utf8.DecodeLastRune(s[:i])
			i -= size
			if r == cr {
				return i
			}
		}
		return -1
	}
	for i := len(s); i > 0; {
		r := rune(s[i-1])
		if r < utf8.RuneSelf {
			if bytealg.IndexByteString(chars, s[i-1]) >= 0 {
				return i - 1
			}
			i--
			continue
		}
		r, size := utf8.DecodeLastRune(s[:i])
		i -= size
		if r != utf8.RuneError {
			// r is 2 to 4 bytes
			if len(chars) == size {
				if chars == string(r) {
					return i
				}
				continue
			}
			// Use bytealg.IndexString for performance if available.
			if bytealg.MaxLen >= size {
				if bytealg.IndexString(chars, string(r)) >= 0 {
					return i
				}
				continue
			}
		}
		for _, ch := range chars {
			if r == ch {
				return i
			}
		}
	}
	return -1
}

// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subslices.
// 通用的分割函数，功能是基于分隔符 sep 将一个字节切片 s 分割成多个子切片，并且会将分隔符本身（包括 sepSave 字节）保留在结果子切片中。
func genSplit(s, sep []byte, sepSave, n int) [][]byte {
	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		// 调用了刚刚在上面所使用的函数
		return explode(s, n)
	}
	if n < 0 {
		n = Count(s, sep) + 1
	}

	a := make([][]byte, n)
	n--
	i := 0
	for i < n {
		m := Index(s, sep)
		if m < 0 {
			break
		}
		a[i] = s[: m+sepSave : m+sepSave]
		s = s[m+len(sep):]
		i++
	}
	a[i] = s
	return a[:i+1]
}

// SplitN slices s into subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, SplitN splits after each UTF-8 sequence.
// The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
//
// To split around the first instance of a separator, see Cut.
// N 的实现设定了自切片的相对应的数量
func SplitN(s, sep []byte, n int) [][]byte { return genSplit(s, sep, 0, n) }

// SplitAfterN slices s into subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfterN splits after each UTF-8 sequence.
// The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
//
// SplitAfterN 是 Go 标准库中用于字节切片分割的一个函数，功能是将字节切片 s 按照分隔符 sep 切割成多个子切片，
// 并在每个分隔符后进行分割。与 SplitN 类似，它也使用了 n 参数来控制返回的子切片的数量，
// 不同之处在于它是 在每个分隔符之后 进行切割。
func SplitAfterN(s, sep []byte, n int) [][]byte {
	return genSplit(s, sep, len(sep), n)
}

// Split slices s into all subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.
//
// To split around the first instance of a separator, see Cut.
// Split 是 Go 标准库中用于分割字节切片的一个函数，其功能是将字节切片 s 按照分隔符 sep 切割成多个子切片，并返回这些子切片。
// 与 SplitN 和 SplitAfterN 不同，Split 会 返回所有的子切片，即分隔符 sep 每次出现的位置都会产生一个子切片。
func Split(s, sep []byte) [][]byte { return genSplit(s, sep, 0, -1) }

// SplitAfter slices s into all subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfter splits after each UTF-8 sequence.
// It is equivalent to SplitAfterN with a count of -1.
func SplitAfter(s, sep []byte) [][]byte {
	return genSplit(s, sep, len(sep), -1)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// Fields interprets s as a sequence of UTF-8-encoded code points.
// It splits the slice s around each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, returning a slice of subslices of s or an
// empty slice if s contains only white space.
// Fields 是 Go 标准库中用于按空白字符分割字节切片的函数。
// 它将输入的字节切片 s 按照一个或多个连续的空白字符进行切割，返回切割后的字节切片数组。
// 空白字符由 unicode.IsSpace 定义，通常包括空格、制表符、换行符等。
func Fields(s []byte) [][]byte {
	// First count the fields.
	// This is an exact count if s is ASCII, otherwise it is an approximation.
	n := 0
	wasSpace := 1
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(s); i++ {
		r := s[i]
		setBits |= r
		isSpace := int(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	if setBits >= utf8.RuneSelf {
		// Some runes in the input slice are not ASCII.
		return FieldsFunc(s, unicode.IsSpace)
	}

	// ASCII fast path
	a := make([][]byte, n)
	na := 0
	fieldStart := 0
	i := 0
	// Skip spaces in the front of the input.
	for i < len(s) && asciiSpace[s[i]] != 0 {
		i++
	}
	fieldStart = i
	for i < len(s) {
		if asciiSpace[s[i]] == 0 {
			i++
			continue
		}
		a[na] = s[fieldStart:i:i]
		na++
		i++
		// Skip spaces in between fields.
		for i < len(s) && asciiSpace[s[i]] != 0 {
			i++
		}
		fieldStart = i
	}
	if fieldStart < len(s) { // Last field might end at EOF.
		a[na] = s[fieldStart:len(s):len(s)]
	}
	return a
}

// FieldsFunc interprets s as a sequence of UTF-8-encoded code points.
// It splits the slice s at each run of code points c satisfying f(c) and
// returns a slice of subslices of s. If all code points in s satisfy f(c), or
// len(s) == 0, an empty slice is returned.
//
// FieldsFunc makes no guarantees about the order in which it calls f(c)
// and assumes that f always returns the same value for a given c.
// FieldsFunc 是 Go 标准库中的一个函数，用于根据自定义的条件函数 f 将字节切片 s 分割为多个子切片。
// 与 Fields 函数不同，FieldsFunc 允许开发者提供一个自定义的判断条件 f，用来决定哪些字符应该作为分隔符。
func FieldsFunc(s []byte, f func(rune) bool) [][]byte {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 32)

	// Find the field start and end indices.
	// Doing this in a separate pass (rather than slicing the string s
	// and collecting the result substrings right away) is significantly
	// more efficient, possibly due to cache effects.
	start := -1 // valid span start if >= 0
	for i := 0; i < len(s); {
		size := 1
		r := rune(s[i])
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(s[i:])
		}
		if f(r) {
			if start >= 0 {
				spans = append(spans, span{start, i})
				start = -1
			}
		} else {
			if start < 0 {
				start = i
			}
		}
		i += size
	}

	// Last field might end at EOF.
	if start >= 0 {
		spans = append(spans, span{start, len(s)})
	}

	// Create subslices from recorded field indices.
	a := make([][]byte, len(spans))
	for i, span := range spans {
		a[i] = s[span.start:span.end:span.end]
	}

	return a
}

// Join concatenates the elements of s to create a new byte slice. The separator
// sep is placed between elements in the resulting slice.
// Join 是一个将多个字节切片合并为一个字节切片的函数，并在每两个元素之间插入一个分隔符。它的功能类似于字符串操作中的 strings.Join，但针对的是字节切片。
func Join(s [][]byte, sep []byte) []byte {
	if len(s) == 0 {
		return []byte{}
	}
	if len(s) == 1 {
		// Just return a copy.
		return append([]byte(nil), s[0]...)
	}
	n := len(sep) * (len(s) - 1)
	for _, v := range s {
		n += len(v)
	}

	b := make([]byte, n)
	bp := copy(b, s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], v)
	}
	return b
}

/*
这段代码定义了两个函数：HasPrefix 和 HasSuffix，它们分别用于检查一个字节切片是否以指定的前缀或后缀开始或结束。
*/
// HasPrefix tests whether the byte slice s begins with prefix.
func HasPrefix(s, prefix []byte) bool {
	return len(s) >= len(prefix) && Equal(s[0:len(prefix)], prefix)
}

// HasSuffix tests whether the byte slice s ends with suffix.
func HasSuffix(s, suffix []byte) bool {
	return len(s) >= len(suffix) && Equal(s[len(s)-len(suffix):], suffix)
}

// 这个 Map 函数的作用是将给定字节切片中的每个字符根据提供的映射函数进行转换，并返回一个新的字节切片。
// 函数会对每个字符进行逐个处理，并根据映射函数的结果生成新的字节切片。
// 如果映射函数返回一个负值，则该字符会被丢弃（即从结果中移除）。映射后的字符仍然会被解释为 UTF-8 编码的代码点。
// Map returns a copy of the byte slice s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the byte slice with no replacement. The characters in s and the
// output are interpreted as UTF-8-encoded code points.
func Map(mapping func(r rune) rune, s []byte) []byte {
	/*
		内存扩展： 在某些情况下，映射后的字节切片可能会比原字节切片大。因此，代码会根据需要动态扩展字节切片的容量。
					虽然这种情况较为罕见，但考虑到所有的字符都可能被修改或丢弃，扩展是必须的。
		字符丢弃： 如果映射函数返回负值（即 rune < 0），则该字符会从结果字节切片中移除。这在需要移除某些字符的情况下非常有用。
	*/
	// In the worst case, the slice can grow when mapped, making
	// things unpleasant. But it's so rare we barge in assuming it's
	// fine. It could also shrink but that falls out naturally.
	maxbytes := len(s) // length of b
	nbytes := 0        // number of bytes encoded in b
	b := make([]byte, maxbytes)
	for i := 0; i < len(s); {
		wid := 1
		r := rune(s[i])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[i:])
		}
		r = mapping(r)
		if r >= 0 {
			rl := utf8.RuneLen(r)
			if rl < 0 {
				rl = len(string(utf8.RuneError))
			}
			if nbytes+rl > maxbytes {
				// Grow the buffer.
				maxbytes = maxbytes*2 + utf8.UTFMax
				nb := make([]byte, maxbytes)
				copy(nb, b[0:nbytes])
				b = nb
			}
			nbytes += utf8.EncodeRune(b[nbytes:maxbytes], r)
		}
		i += wid
	}
	return b[0:nbytes]
}

// 这个 Repeat 函数的作用是生成一个新的字节切片，该字节切片由原始字节切片 b 的 count 次重复组成。
// 该函数执行前会检查是否存在可能的溢出情况，并根据条件进行处理。如果 count 为负值或者重复后字节切片长度的计算会导致溢出，函数会触发 panic。
// Repeat returns a new byte slice consisting of count copies of b.
//
// It panics if count is negative or if
// the result of (len(b) * count) overflows.
func Repeat(b []byte, count int) []byte {
	if count == 0 {
		return []byte{}
	}
	// Since we cannot return an error on overflow,
	// we should panic if the repeat will generate
	// an overflow.
	// See Issue golang.org/issue/16237.
	if count < 0 {
		panic("bytes: negative Repeat count")
	} else if len(b)*count/count != len(b) {
		panic("bytes: Repeat count causes overflow")
	}

	nb := make([]byte, len(b)*count)
	bp := copy(nb, b)
	for bp < len(nb) {
		copy(nb[bp:], nb[:bp])
		bp *= 2
	}
	return nb
}

// 所谓的大小写的转化就是 ASCII 字节的大小写转化，不是 UTF-8 字节的大小写转化。
// ToUpper returns a copy of the byte slice s with all Unicode letters mapped to
// their upper case.
func ToUpper(s []byte) []byte {
	isASCII, hasLower := true, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasLower = hasLower || ('a' <= c && c <= 'z')
	}

	if isASCII { // optimize for ASCII-only byte slices.
		if !hasLower {
			// Just return a copy.
			return append([]byte(""), s...)
		}
		b := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			b[i] = c
		}
		return b
	}
	return Map(unicode.ToUpper, s)
}

// ToLower returns a copy of the byte slice s with all Unicode letters mapped to
// their lower case.
func ToLower(s []byte) []byte {
	isASCII, hasUpper := true, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}

	if isASCII { // optimize for ASCII-only byte slices.
		if !hasUpper {
			return append([]byte(""), s...)
		}
		b := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}
			b[i] = c
		}
		return b
	}
	return Map(unicode.ToLower, s)
}

// ToTitle treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their title case.
func ToTitle(s []byte) []byte { return Map(unicode.ToTitle, s) }

// ToUpperSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// upper case, giving priority to the special casing rules.
func ToUpperSpecial(c unicode.SpecialCase, s []byte) []byte {
	return Map(c.ToUpper, s)
}

// ToLowerSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// lower case, giving priority to the special casing rules.
func ToLowerSpecial(c unicode.SpecialCase, s []byte) []byte {
	return Map(c.ToLower, s)
}

// ToTitleSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// title case, giving priority to the special casing rules.
func ToTitleSpecial(c unicode.SpecialCase, s []byte) []byte {
	return Map(c.ToTitle, s)
}

// ToValidUTF8 treats s as UTF-8-encoded bytes and returns a copy with each run of bytes
// representing invalid UTF-8 replaced with the bytes in replacement, which may be empty.
func ToValidUTF8(s, replacement []byte) []byte {
	b := make([]byte, 0, len(s)+len(replacement))
	invalid := false // previous byte was from an invalid UTF-8 sequence
	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			i++
			invalid = false
			b = append(b, c)
			continue
		}
		_, wid := utf8.DecodeRune(s[i:])
		if wid == 1 {
			i++
			if !invalid {
				invalid = true
				b = append(b, replacement...)
			}
			continue
		}
		invalid = false
		b = append(b, s[i:i+wid]...)
		i += wid
	}
	return b
}

// isSeparator 函数用于判断给定的字符（rune 类型）是否能标记一个单词边界。在很多文本处理任务中，我们需要识别出哪些字符作为单词的分隔符，这个函数正是实现了这个功能。
// isSeparator reports whether the rune could mark a word boundary.
// TODO: update when package unicode captures more of the properties.
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

// Title 函数用于将给定的 UTF-8 编码字节切片 s 中的每个单词的首字母转为标题格式（即大写字母），并返回一个新的字节切片。标题格式是指每个单词的首字母大写，其他字母小写。
// 方便的形式就可以提供给外界进行访问的
// Title treats s as UTF-8-encoded bytes and returns a copy with all Unicode letters that begin
// words mapped to their title case.
//
// Deprecated: The rule Title uses for word boundaries does not handle Unicode
// punctuation properly. Use golang.org/x/text/cases instead.
func Title(s []byte) []byte {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		s)
}

// 它们根据不同的条件去掉字节切片的前缀、后缀或满足特定条件的字符，然后返回一个新的字节切片。
// TrimLeftFunc returns a subslice of s by slicing off
// TrimLeftFunc treats s as UTF-8-encoded bytes and returns a subslice of s by slicing off
// all leading UTF-8-encoded code points c that satisfy f(c).
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte {
	i := indexFunc(s, f, false)
	if i == -1 {
		return nil
	}
	return s[i:]
}

// TrimRightFunc returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points c that satisfy f(c).
func TrimRightFunc(s []byte, f func(r rune) bool) []byte {
	i := lastIndexFunc(s, f, false)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRune(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

// TrimFunc returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded code points c that satisfy f(c).
func TrimFunc(s []byte, f func(r rune) bool) []byte {
	return TrimRightFunc(TrimLeftFunc(s, f), f)
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix []byte) []byte {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func TrimSuffix(s, suffix []byte) []byte {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
// TrimLeftFunc 和 TrimRightFunc 分别去掉前导和后缀符合条件的字符，而 TrimFunc 一次性去掉两端符合条件的字符。
// TrimPrefix 和 TrimSuffix 则分别针对字节切片的前缀和后缀进行处理。

// IndexFunc interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index in s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func IndexFunc(s []byte, f func(r rune) bool) int {
	return indexFunc(s, f, true)
}

// LastIndexFunc interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index in s of the last Unicode
// code point satisfying f(c), or -1 if none do.
func LastIndexFunc(s []byte, f func(r rune) bool) int {
	return lastIndexFunc(s, f, true)
}

// indexFunc is the same as IndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func indexFunc(s []byte, f func(r rune) bool, truth bool) int {
	start := 0
	for start < len(s) {
		wid := 1
		r := rune(s[start])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[start:])
		}
		if f(r) == truth {
			return start
		}
		start += wid
	}
	return -1
}

// lastIndexFunc is the same as LastIndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
// 通过这个参数，你可以灵活地控制查找的条件，避免额外编写多个函数来反转条件。
// 只需要通过设置 truth 为 true 或 false，就可以改变查找行为，查找符合条件或不符合条件的字符。
func lastIndexFunc(s []byte, f func(r rune) bool, truth bool) int {
	for i := len(s); i > 0; {
		r, size := rune(s[i-1]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeLastRune(s[0:i])
		}
		i -= size
		if f(r) == truth {
			return i
		}
	}
	return -1
}

// asciiSet is a 32-byte value, where each bit represents the presence of a
// given ASCII character in the set. The 128-bits of the lower 16 bytes,
// starting with the least-significant bit of the lowest word to the
// most-significant bit of the highest word, map to the full range of all
// 128 ASCII characters. The 128-bits of the upper 16 bytes will be zeroed,
// ensuring that any non-ASCII character will be reported as not in the set.
// This allocates a total of 32 bytes even though the upper half
// is unused to avoid bounds checks in asciiSet.contains.
// asciiSet 通过位运算加速字符查找： 通过将字符映射到位图中的特定位，能在常数时间内确定字符是否存在，而不像传统的字符比较需要逐一匹配字符。
// 避免内存分配： containsRune 函数避免了使用额外的字符串拷贝或其他内存分配操作，这对于频繁的字符检查尤为重要。
type asciiSet [8]uint32

// makeASCIISet creates a set of ASCII characters and reports whether all
// characters in chars are ASCII.
func makeASCIISet(chars string) (as asciiSet, ok bool) {
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c/32] |= 1 << (c % 32)
	}
	return as, true
}

// contains reports whether c is inside the set.
func (as *asciiSet) contains(c byte) bool {
	return (as[c/32] & (1 << (c % 32))) != 0
}

// containsRune is a simplified version of strings.ContainsRune
// to avoid importing the strings package.
// We avoid bytes.ContainsRune to avoid allocating a temporary copy of s.
func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
// 这是一个高层次的修剪函数，用于去掉字符串两端的指定字符。它首先会检查 cutset 是否为空或单个字节字符。如果是，它会直接调用专门的函数来处理这些情况。
// 否则，如果 cutset 中的所有字符都是 ASCII 字符，它会使用 asciiSet 来提高效率。对于非 ASCII 字符，它会使用 Unicode 字符串来进行修剪。
// Trim returns a subslice of s by slicing off all leading and
// trailing UTF-8-encoded code points contained in cutset.
func Trim(s []byte, cutset string) []byte {
	if len(s) == 0 || cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimLeftByte(trimRightByte(s, cutset[0]), cutset[0])
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimLeftASCII(trimRightASCII(s, &as), &as)
	}
	return trimLeftUnicode(trimRightUnicode(s, cutset), cutset)
}

// TrimLeft returns a subslice of s by slicing off all leading
// UTF-8-encoded code points contained in cutset.
func TrimLeft(s []byte, cutset string) []byte {
	if len(s) == 0 || cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimLeftByte(s, cutset[0])
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimLeftASCII(s, &as)
	}
	return trimLeftUnicode(s, cutset)
}

func trimLeftByte(s []byte, c byte) []byte {
	for len(s) > 0 && s[0] == c {
		s = s[1:]
	}
	return s
}

func trimLeftASCII(s []byte, as *asciiSet) []byte {
	for len(s) > 0 {
		if !as.contains(s[0]) {
			break
		}
		s = s[1:]
	}
	return s
}

func trimLeftUnicode(s []byte, cutset string) []byte {
	for len(s) > 0 {
		r, n := rune(s[0]), 1
		if r >= utf8.RuneSelf {
			r, n = utf8.DecodeRune(s)
		}
		if !containsRune(cutset, r) {
			break
		}
		s = s[n:]
	}
	return s
}

// TrimRight returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points that are contained in cutset.
func TrimRight(s []byte, cutset string) []byte {
	if len(s) == 0 || cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimRightByte(s, cutset[0])
	}
	if as, ok := makeASCIISet(cutset); ok {
		return trimRightASCII(s, &as)
	}
	return trimRightUnicode(s, cutset)
}

func trimRightByte(s []byte, c byte) []byte {
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

func trimRightASCII(s []byte, as *asciiSet) []byte {
	for len(s) > 0 {
		if !as.contains(s[len(s)-1]) {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}

func trimRightUnicode(s []byte, cutset string) []byte {
	for len(s) > 0 {
		r, n := rune(s[len(s)-1]), 1
		if r >= utf8.RuneSelf {
			r, n = utf8.DecodeLastRune(s)
		}
		if !containsRune(cutset, r) {
			break
		}
		s = s[:len(s)-n]
	}
	return s
}

// TrimSpace returns a subslice of s by slicing off all leading and
// trailing white space, as defined by Unicode.
func TrimSpace(s []byte) []byte {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// Now look for the first ASCII non-space byte from the end
	stop := len(s)
	for ; stop > start; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return TrimFunc(s[start:stop], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	if start == stop {
		// Special case to preserve previous TrimLeftFunc behavior,
		// returning nil instead of empty slice if all spaces.
		return nil
	}
	return s[start:stop]
}

// Runes interprets s as a sequence of UTF-8-encoded code points.
// It returns a slice of runes (Unicode code points) equivalent to s.
// 将 UTF-8 编码的字节切片 s 解释为 Unicode 代码点序列，并返回一个包含这些 Unicode 代码点（rune）的切片。
func Runes(s []byte) []rune {
	t := make([]rune, utf8.RuneCount(s))
	i := 0
	for len(s) > 0 {
		r, l := utf8.DecodeRune(s)
		t[i] = r
		i++
		s = s[l:]
	}
	return t
}

// Replace returns a copy of the slice s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the slice
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune slice.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new []byte, n int) []byte {
	m := 0
	if n != 0 {
		// Compute number of replacements.
		m = Count(s, old)
	}
	if m == 0 {
		// Just return a copy.
		return append([]byte(nil), s...)
	}
	if n < 0 || m < n {
		n = m
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRune(s[start:])
				j += wid
			}
		} else {
			j += Index(s[start:], old)
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], new)
		start = j + len(old)
	}
	w += copy(t[w:], s[start:])
	return t[0:w]
}

// ReplaceAll returns a copy of the slice s with all
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the slice
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune slice.
func ReplaceAll(s, old, new []byte) []byte {
	return Replace(s, old, new, -1)
}

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, which is a more general
// form of case-insensitivity.
// 判断两个字节切片 s 和 t 是否在 Unicode 上进行大小写不敏感的比较时相等。它考虑了大小写折叠（即忽略大小写差异）。
func EqualFold(s, t []byte) bool {
	for len(s) != 0 && len(t) != 0 {
		// Extract first rune from each.
		var sr, tr rune
		if s[0] < utf8.RuneSelf {
			sr, s = rune(s[0]), s[1:]
		} else {
			r, size := utf8.DecodeRune(s)
			sr, s = r, s[size:]
		}
		if t[0] < utf8.RuneSelf {
			tr, t = rune(t[0]), t[1:]
		} else {
			r, size := utf8.DecodeRune(t)
			tr, t = r, t[size:]
		}

		// If they match, keep going; if not, return false.

		// Easy case.
		if tr == sr {
			continue
		}

		// Make sr < tr to simplify what follows.
		if tr < sr {
			tr, sr = sr, tr
		}
		// Fast check for ASCII.
		if tr < utf8.RuneSelf {
			// ASCII only, sr/tr must be upper/lower case
			if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
				continue
			}
			return false
		}

		// General case. SimpleFold(x) returns the next equivalent rune > x
		// or wraps around to smaller values.
		r := unicode.SimpleFold(sr)
		for r != sr && r < tr {
			r = unicode.SimpleFold(r)
		}
		if r == tr {
			continue
		}
		return false
	}

	// One string is empty. Are both?
	return len(s) == len(t)
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
// 功能：返回字节序列 s 中第一次出现字节序列 sep 的索引位置，如果未找到，返回 -1。
func Index(s, sep []byte) int {
	n := len(sep)
	switch {
	case n == 0:
		return 0
	case n == 1:
		return IndexByte(s, sep[0])
	case n == len(s):
		if Equal(sep, s) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	case n <= bytealg.MaxLen:
		// Use brute force when s and sep both are small
		if len(s) <= bytealg.MaxBruteForce {
			return bytealg.Index(s, sep)
		}
		c0 := sep[0]
		c1 := sep[1]
		i := 0
		t := len(s) - n + 1
		fails := 0
		for i < t {
			if s[i] != c0 {
				// IndexByte is faster than bytealg.Index, so use it as long as
				// we're not getting lots of false positives.
				o := IndexByte(s[i+1:t], c0)
				if o < 0 {
					return -1
				}
				i += o + 1
			}
			if s[i+1] == c1 && Equal(s[i:i+n], sep) {
				return i
			}
			fails++
			i++
			// Switch to bytealg.Index when IndexByte produces too many false positives.
			if fails > bytealg.Cutover(i) {
				r := bytealg.Index(s[i:], sep)
				if r >= 0 {
					return r + i
				}
				return -1
			}
		}
		return -1
	}
	c0 := sep[0]
	c1 := sep[1]
	i := 0
	fails := 0
	t := len(s) - n + 1
	for i < t {
		if s[i] != c0 {
			o := IndexByte(s[i+1:t], c0)
			if o < 0 {
				break
			}
			i += o + 1
		}
		if s[i+1] == c1 && Equal(s[i:i+n], sep) {
			return i
		}
		i++
		fails++
		if fails >= 4+i>>4 && i < t {
			// Give up on IndexByte, it isn't skipping ahead
			// far enough to be better than Rabin-Karp.
			// Experiments (using IndexPeriodic) suggest
			// the cutover is about 16 byte skips.
			// TODO: if large prefixes of sep are matching
			// we should cutover at even larger average skips,
			// because Equal becomes that much more expensive.
			// This code does not take that effect into account.
			j := bytealg.IndexRabinKarpBytes(s[i:], sep)
			if j < 0 {
				return -1
			}
			return i + j
		}
	}
	return -1
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, nil, false.
//
// Cut returns slices of the original slice s, not copies.
// 将字节序列 s 按照第一次出现的 sep 切割成两部分。返回切割后的两部分以及是否找到了分隔符。
func Cut(s, sep []byte) (before, after []byte, found bool) {
	if i := Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, nil, false
}
