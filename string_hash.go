package g

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an String for hashing.
type shash struct{ str String }

// Hash returns a shash struct wrapping the given String.
func (s String) Hash() shash { return shash{s} }

// MD5 computes the MD5 hash of the wrapped String and returns the hash as an String.
func (sh shash) MD5() String { return stringHasher(md5.New(), sh.str) }

// SHA1 computes the SHA1 hash of the wrapped String and returns the hash as an String.
func (sh shash) SHA1() String { return stringHasher(sha1.New(), sh.str) }

// SHA256 computes the SHA256 hash of the wrapped String and returns the hash as an String.
func (sh shash) SHA256() String { return stringHasher(sha256.New(), sh.str) }

// SHA512 computes the SHA512 hash of the wrapped String and returns the hash as an String.
func (sh shash) SHA512() String { return stringHasher(sha512.New(), sh.str) }

// stringHasher a helper function that computes the hash of the given String using the specified
// hash.Hash algorithm and returns the hash as an String.
func stringHasher(algorithm hash.Hash, hstr String) String {
	_, _ = algorithm.Write(hstr.Bytes())
	return String(hex.EncodeToString(algorithm.Sum(nil)))
}
