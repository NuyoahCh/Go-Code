// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//版权所有：The Go Authors。版权所有。 
//该源代码的使用由bsd风格管理 
//可以在license文件中找到的许可证。

/*
	Package builtin provides documentation for Go's predeclared identifiers.
	The items documented here are not actually in package builtin
	but their descriptions here allow godoc to present documentation
	for the language's special identifiers.
    包内置为Go的预声明标识符提供了文档。 
    这里记录的项目实际上不在包内构建中 
    但是它们在这里的描述允许godoc提供文档 
    语言的特殊标识符。
*/
package builtin

// bool is the set of boolean values, true and false.
// bool是布尔值的集合，true和false。
type bool bool

// true and false are the two untyped boolean values.
// true和false是两个未类型化的布尔值。
const (
    true  = 0 == 0 // Untyped bool.
    false = 0 != 0 // Untyped bool.
)

/**
    在34-109行定义了 Go 语言中的有/无符号的基本数据类型
*/
// uint8 is the set of all unsigned 8-bit integers.
// uint8 是所有的 8 bit 无符号整数的集合
// Range: 0 through 255.
type uint8 uint8

// uint16 is the set of all unsigned 16-bit integers.
// Range: 0 through 65535.
type uint16 uint16

// uint32 is the set of all unsigned 32-bit integers.
// Range: 0 through 4294967295.
type uint32 uint32

// uint64 is the set of all unsigned 64-bit integers.
// Range: 0 through 18446744073709551615.
type uint64 uint64

// 有符号整型，范围的划定是要去考虑 0 的存在的
// int8 is the set of all signed 8-bit integers.
// Range: -128 through 127.
type int8 int8

// int16 is the set of all signed 16-bit integers.
// Range: -32768 through 32767.
type int16 int16

// int32 is the set of all signed 32-bit integers.
// Range: -2147483648 through 2147483647.
type int32 int32

// int64 is the set of all signed 64-bit integers.
// Range: -9223372036854775808 through 9223372036854775807.
type int64 int64

// float32 is the set of all IEEE-754 32-bit floating-point numbers.
type float32 float32

// float64 is the set of all IEEE-754 64-bit floating-point numbers.
type float64 float64

// 复数类型
// complex64 is the set of all complex numbers with float32 real and
// imaginary parts.
type complex64 complex64

// complex128 is the set of all complex numbers with float64 real and
// imaginary parts.
type complex128 complex128

// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string

// int is a signed integer type that is at least 32 bits in size. It is a
// distinct type, however, and not an alias for, say, int32.
type int int

// uint is an unsigned integer type that is at least 32 bits in size. It is a
// distinct type, however, and not an alias for, say, uint32.
type uint uint


// uintptr是一个足够大的整数类型，可以保存类型的位模式，任意指针。
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr

// byte is an alias for uint8 and is equivalent to uint8 in all ways. It is
// used, by convention, to distinguish byte values from 8-bit unsigned
// integer values.
type byte = uint8

// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32

// any is an alias for interface{} and is equivalent to interface{} in all ways.
// any是interface{}的别名，在所有方面都等同于interface{}。
type any = interface{}

// comparable is an interface that is implemented by all comparable types
// (booleans, numbers, strings, pointers, channels, arrays of comparable types,
// structs whose fields are all comparable types).
// The comparable interface may only be used as a type parameter constraint,
// not as the type of a variable.
// comparable 是由所有可比较类型实现的接口
// (布尔值、数字、字符串、指针、通道、可比较类型的数组； 字段都是可比较类型的结构体)。 
// 可比较的接口只能用作类型参数约束， 
// 不是作为变量的类型
type comparable interface{ comparable }

// iota is a predeclared identifier representing the untyped integer ordinal
// number of the current const specification in a (usually parenthesized)
// const declaration. It is zero-indexed.
// iota是一个预先声明的标识符，表示未类型的整数序数 
// 当前const规范的编号（通常带括号） 
// const声明。它是零索引的。
const iota = 0 // Untyped int. 没有形式的整型

// nil is a predeclared identifier representing the zero value for a
// pointer, channel, func, interface, map, or slice type.
// nil 是预先声明的标识符，表示a的零值 
// 指针、通道、函数、接口、映射或切片类型。 
// 类型必须是指针、通道、函数、接口、映射或切片类型
var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

/**
    这些类型定义在 Go 代码中通常是为了文档目的和代码的可读性，帮助开发者理解代码的意图。
    它们并不真正提供新的数据类型，而是作为“占位符”类型，用来标明某些函数或方法可以接受不同的类型，
    具体类型可以在实现中根据实际需要替换。通过这种方式，可以使文档或代码更加清晰和易于理解
*/
// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
// 这里的 Type 仅用于文档目的。这是一个替身 
// 用于任何Go语言类型，但表示任何给定函数的相同类型 
// 调用。
type Type int

// Type1 is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type1 int

// IntegerType is here for the purposes of documentation only. It is a stand-in
// for any integer type: int, uint, int8 etc.
type IntegerType int

// FloatType is here for the purposes of documentation only. It is a stand-in
// for either float type: float32 or float64.
type FloatType float32

// ComplexType is here for the purposes of documentation only. It is a
// stand-in for either complex type: complex64 or complex128.
type ComplexType complex64

// The append built-in function appends elements to the end of a slice. If
// it has sufficient capacity, the destination is resliced to accommodate the
// new elements. If it does not, a new underlying array will be allocated.
// Append returns the updated slice. It is therefore necessary to store the
// result of append, often in the variable holding the slice itself:
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
// As a special case, it is legal to append a string to a byte slice, like this:
//	slice = append([]byte("hello "), "world"...)
// 主要是用于切片的追加
func append(slice []Type, elems ...Type) []Type

// The copy built-in function copies elements from a source slice into a
// destination slice. (As a special case, it also will copy bytes from a
// string to a slice of bytes.) The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
// 这个 copy 函数是 Go 语言的一个内置函数，
// 用来将一个源切片（src）的元素复制到目标切片（dst）中
// dst: 目标切片（接收数据的切片）。
// src: 源切片（提供数据的切片）。
// copy 函数是 Go 中一个非常有用的内置函数，
// 能够高效、灵活地将数据从一个切片复制到另一个切片。
// 它可以处理切片之间的重叠情况，也能从字符串复制到字节切片。
// 它返回复制的元素数量，使得开发者可以根据返回值进行后续的操作。
func copy(dst, src []Type) int

// The delete built-in function deletes the element with the specified key
// (m[key]) from the map. If m is nil or there is no such element, delete
// is a no-op.
// delete 函数用于从映射中删除指定键的元素。
// 如果指定键存在，元素会被删除；如果不存在，什么也不做。
// 即使映射为空（nil）或映射中没有指定的键，delete 也不会引发错误，处理是安全的。
func delete(m map[Type]Type1, key Type)

// The len built-in function returns the length of v, according to its type:
//	Array: the number of elements in v.
//	Pointer to array: the number of elements in *v (even if v is nil).
//	Slice, or map: the number of elements in v; if v is nil, len(v) is zero.
//	String: the number of bytes in v.
//	Channel: the number of elements queued (unread) in the channel buffer;
//	         if v is nil, len(v) is zero.
// For some arguments, such as a string literal or a simple array expression, the
// result can be a constant. See the Go language specification's "Length and
// capacity" section for details.
// v: 传入的变量，可以是数组、切片、映射（map）、字符串或通道（channel）等。
// 返回值: 返回一个 int 类型的值，表示 v 的长度，具体长度的定义取决于 v 的类型。
// 对于字符串，len 返回字符串中的字节数。需要注意的是，Go 中的字符串是以 UTF-8 编码存储的，
// len 返回的是字符串的字节数，而不是字符的数量。如果字符串包含多字节字符（如中文字符），
// 它们的字节数会大于 1。“你好” - - 就是看作6个字符
func len(v Type) int

// The cap built-in function returns the capacity of v, according to its type:
//	Array: the number of elements in v (same as len(v)).
//	Pointer to array: the number of elements in *v (same as len(v)).
//	Slice: the maximum length the slice can reach when resliced;
//	if v is nil, cap(v) is zero.
//	Channel: the channel buffer capacity, in units of elements;
//	if v is nil, cap(v) is zero.
// For some arguments, such as a simple array expression, the result can be a
// constant. See the Go language specification's "Length and capacity" section for
// details.
// 作用: cap 函数用于返回给定数据类型（如数组、切片、通道等）的容量（capacity）。
// 返回值: 数组和指针指向的数组：返回数组的元素数量（即 len(v)）。
// 切片：返回切片的最大容量，表示切片的底层数组能容纳的元素个数。
// 通道：返回通道的缓冲区容量，即通道中可以存放的元素数量。
// 如果是 nil，返回 0。
func cap(v Type) int

// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
// make 用于初始化并分配一个切片、映射（map）或通道（channel）。
// make 的返回类型是与类型 t 相同的对象，而不是指向该类型的指针。
// 对于切片：size 指定切片的长度，第二个参数可以指定切片的容量。
// 对于映射：size 用于指定映射的初始大小（即存储桶的数量）。
// 对于通道：size 用于指定通道的缓冲区大小，若为 0 或省略，则为无缓冲通道。
// 不会返回指针，而是初始化一个默认值
func make(t Type, size ...IntegerType) Type

// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
// 作用: new 用于为给定类型分配内存，并返回该类型的指针，初始化为该类型的零值。
// 返回值: 返回一个指向新分配的零值的指针。
func new(Type) *Type

// The complex built-in function constructs a complex value from two
// floating-point values. The real and imaginary parts must be of the same
// size, either float32 or float64 (or assignable to them), and the return
// value will be the corresponding complex type (complex64 for float32,
// complex128 for float64).
// 作用: complex 用于构造一个复数值。它接受两个浮点数参数，表示复数的实部和虚部。
// 返回值: 返回一个复数类型（complex64 或 complex128）。
func complex(r, i FloatType) ComplexType

// The real built-in function returns the real part of the complex number c.
// The return value will be floating point type corresponding to the type of c.
// 作用: 返回复数的实部。
// 返回值: 返回与复数 c 相对应的浮点类型的实部。
func real(c ComplexType) FloatType

// The imag built-in function returns the imaginary part of the complex
// number c. The return value will be floating point type corresponding to
// the type of c.
// 作用: 返回复数的虚部。
// 返回值: 返回与复数 c 相对应的浮点类型的虚部。
func imag(c ComplexType) FloatType

// The close built-in function closes a channel, which must be either
// bidirectional or send-only. It should be executed only by the sender,
// never the receiver, and has the effect of shutting down the channel after
// the last sent value is received. After the last value has been received
// from a closed channel c, any receive from c will succeed without
// blocking, returning the zero value for the channel element. The form
//	x, ok := <-c
// will also set ok to false for a closed channel.
// 作用: close 用于关闭一个通道。关闭通道后，不能再向通道发送数据，但可以接收数据。
// 如果通道已经关闭，接收操作会立即返回通道元素的零值，并且 ok 为 false。
// 注意: 关闭通道的操作只能由发送方执行，接收方不应该关闭通道。
// 关闭了之后接收方也只能是接收到0 （即对应类型的零值）
func close(c chan<- Type)

// The panic built-in function stops normal execution of the current
// goroutine. When a function F calls panic, normal execution of F stops
// immediately. Any functions whose execution was deferred by F are run in
// the usual way, and then F returns to its caller. To the caller G, the
// invocation of F then behaves like a call to panic, terminating G's
// execution and running any deferred functions. This continues until all
// functions in the executing goroutine have stopped, in reverse order. At
// that point, the program is terminated with a non-zero exit code. This
// termination sequence is called panicking and can be controlled by the
// built-in function recover.
// 作用: panic 用于停止当前 goroutine 的正常执行，并触发 panic 过程。
// panic 会导致当前函数的执行立即终止，执行已 defer 的函数，然后返回到上层调用者，直到程序终止。
// 应用: 通常用于程序中出现不可恢复的错误时。
func panic(v any)

// The recover built-in function allows a program to manage behavior of a
// panicking goroutine. Executing a call to recover inside a deferred
// function (but not any function called by it) stops the panicking sequence
// by restoring normal execution and retrieves the error value passed to the
// call of panic. If recover is called outside the deferred function it will
// not stop a panicking sequence. In this case, or when the goroutine is not
// panicking, or if the argument supplied to panic was nil, recover returns
// nil. Thus the return value from recover reports whether the goroutine is
// panicking.
// 作用: recover 允许在 panic 发生时，恢复 goroutine 的正常执行。
// recover 只能在 defer 中调用。如果在其他地方调用，无法停止 panic 过程。
// 返回值:如果在 panic 过程中调用，recover 返回 panic 的值；如果当前没有 panic，返回 nil。
func recover() any

// The print built-in function formats its arguments in an
// implementation-specific way and writes the result to standard error.
// Print is useful for bootstrapping and debugging; it is not guaranteed
// to stay in the language.
// 作用: recover 允许在 panic 发生时，恢复 goroutine 的正常执行。
// recover 只能在 defer 中调用。如果在其他地方调用，无法停止 panic 过程。
// 返回值: 如果在 panic 过程中调用，recover 返回 panic 的值；
// 如果当前没有 panic，返回 nil。
func print(args ...Type)

// The println built-in function formats its arguments in an
// implementation-specific way and writes the result to standard error.
// Spaces are always added between arguments and a newline is appended.
// Println is useful for bootstrapping and debugging; it is not guaranteed
// to stay in the language.
func println(args ...Type)

// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
// 作用: error 是一个标准的接口类型，用于表示错误。Go 中的所有错误都实现了这个接口，
// 其方法 Error() 返回一个字符串，描述错误的详细信息。
type error interface {
	Error() string
}