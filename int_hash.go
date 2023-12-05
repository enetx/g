package g

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an Int for hashing.
type ihash struct{ int Int }

// Hash returns a ihash struct wrapping the given Int.
func (i Int) Hash() ihash { return ihash{i} }

// MD5 computes the MD5 hash of the wrapped Int and returns the hash as an String.
func (ih ihash) MD5() String { return intHasher(md5.New(), ih.int) }

// SHA1 computes the SHA1 hash of the wrapped Int and returns the hash as an String.
func (ih ihash) SHA1() String { return intHasher(sha1.New(), ih.int) }

// SHA256 computes the SHA256 hash of the wrapped Int and returns the hash as an String.
func (ih ihash) SHA256() String { return intHasher(sha256.New(), ih.int) }

// SHA512 computes the SHA512 hash of the wrapped Int and returns the hash as an String.
func (ih ihash) SHA512() String { return intHasher(sha512.New(), ih.int) }

// intHasher a helper function that computes the hash of the given Int using the specified hash.Hash algorithm and returns the hash as an String.
func intHasher(algorithm hash.Hash, ih Int) String {
	_, _ = algorithm.Write(ih.Bytes())
	return String(hex.EncodeToString(algorithm.Sum(nil)))
}
