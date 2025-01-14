// Package conv
// Example for StrToBytes:
// str := "example"
// bytes := StrToBytes(str)
// fmt.Printf("%v\\n", bytes)
//
// Example for BytesToStr:
// bytes := []byte("example")
// str := BytesToStr(bytes)
// fmt.Printf("%s\\n", str)
package conv

import "unsafe"

// StrToBytes converts a string to a byte slice without copying data.
// The returned []byte shares the same underlying memory as the input string.
// WARNING: Modifying the []byte can lead to undefined behavior as strings are immutable in Go.
// Use only in performance-critical scenarios where immutability can be guaranteed.
func StrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToStr converts a byte slice to a string without copying data.
// The returned string shares the same underlying memory as the input []byte.
// WARNING: The input []byte mustn't be modified after the conversion, as strings are immutable.
// Use only when you can ensure the byte slice's immutability after conversion.
func BytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
