// Package fieldhasher provides a transformer that hashes one or more log
// entry fields and writes the resulting hex digest into a configurable
// destination field.
//
// Supported algorithms: md5, sha1, sha256 (default).
//
// Example usage:
//
//	h := fieldhasher.New([]fieldhasher.Rule{
//		{
//			Fields: []string{"email"},
//			Dest:   "email_hash",
//			Algo:   fieldhasher.SHA256,
//		},
//	})
//	result := h.Apply(entry)
//
// When multiple Fields are specified their string representations are
// concatenated in order before hashing, allowing composite keys.
package fieldhasher
