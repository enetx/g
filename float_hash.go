package g

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an Float for hashing.
type fhash struct{ float Float }

// Hash returns a fhash struct wrapping the given Float.
func (f Float) Hash() fhash { return fhash{f} }

// MD5 computes the MD5 hash of the wrapped Float and returns the hash as an String.
func (fh fhash) MD5() String { return floatHasher(md5.New(), fh.float) }

// SHA1 computes the SHA1 hash of the wrapped Float and returns the hash as an String.
func (fh fhash) SHA1() String { return floatHasher(sha1.New(), fh.float) }

// SHA256 computes the SHA256 hash of the wrapped Float and returns the hash as an String.
func (fh fhash) SHA256() String { return floatHasher(sha256.New(), fh.float) }

// SHA512 computes the SHA512 hash of the wrapped Float and returns the hash as an String.
func (fh fhash) SHA512() String { return floatHasher(sha512.New(), fh.float) }

// floatHasher a helper function that computes the hash of the given Float using the specified
// hash.Hash algorithm and returns the hash as an String.
func floatHasher(algorithm hash.Hash, f Float) String {
	_, _ = algorithm.Write(f.Bytes())
	return String(hex.EncodeToString(algorithm.Sum(nil)))
}
