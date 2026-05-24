// Package fieldencoding provides a transformer that encodes or decodes
// string fields in a log entry using common schemes such as base64 and hex.
//
// Each Rule targets a named field and specifies a Mode:
//
//	"base64-encode" — encode the field value as standard base64
//	"base64-decode" — decode a base64-encoded field value
//	"hex-encode"    — encode the field value as lowercase hexadecimal
//	"hex-decode"    — decode a hexadecimal field value
//
// An optional Dest field may be specified; if omitted the transformation
// is applied in-place.
package fieldencoding
