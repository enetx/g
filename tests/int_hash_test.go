package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestInt_Hash_MD5(t *testing.T) {
	testInt := g.Int(42)
	hash := testInt.Hash().MD5()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("MD5 hash should not be empty")
	}

	// Hash should be 32 characters (16 bytes * 2 hex chars)
	if len(hash) != 32 {
		t.Errorf("MD5 hash length should be 32, got %d", len(hash))
	}
}

func TestInt_Hash_SHA1(t *testing.T) {
	testInt := g.Int(100)
	hash := testInt.Hash().SHA1()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA1 hash should not be empty")
	}

	// Hash should be 40 characters (20 bytes * 2 hex chars)
	if len(hash) != 40 {
		t.Errorf("SHA1 hash length should be 40, got %d", len(hash))
	}
}

func TestInt_Hash_SHA256(t *testing.T) {
	testInt := g.Int(256)
	hash := testInt.Hash().SHA256()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA256 hash should not be empty")
	}

	// Hash should be 64 characters (32 bytes * 2 hex chars)
	if len(hash) != 64 {
		t.Errorf("SHA256 hash length should be 64, got %d", len(hash))
	}
}

func TestInt_Hash_SHA512(t *testing.T) {
	testInt := g.Int(512)
	hash := testInt.Hash().SHA512()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA512 hash should not be empty")
	}

	// Hash should be 128 characters (64 bytes * 2 hex chars)
	if len(hash) != 128 {
		t.Errorf("SHA512 hash length should be 128, got %d", len(hash))
	}
}

func TestInt_Hash_Zero(t *testing.T) {
	zeroInt := g.Int(0)

	md5Hash := zeroInt.Hash().MD5()
	sha1Hash := zeroInt.Hash().SHA1()
	sha256Hash := zeroInt.Hash().SHA256()
	sha512Hash := zeroInt.Hash().SHA512()

	// Check that zero produces non-empty hashes
	if md5Hash.Empty() {
		t.Error("MD5 hash of zero should not be empty")
	}
	if sha1Hash.Empty() {
		t.Error("SHA1 hash of zero should not be empty")
	}
	if sha256Hash.Empty() {
		t.Error("SHA256 hash of zero should not be empty")
	}
	if sha512Hash.Empty() {
		t.Error("SHA512 hash of zero should not be empty")
	}
}

func TestInt_Hash_Consistency(t *testing.T) {
	testInt := g.Int(999)

	// Hash the same integer multiple times
	hash1 := testInt.Hash().MD5()
	hash2 := testInt.Hash().MD5()
	hash3 := testInt.Hash().MD5()

	// Results should be consistent
	if !hash1.Eq(hash2) || !hash2.Eq(hash3) {
		t.Error("Hash results should be consistent for the same input")
	}
}

func TestInt_Hash_DifferentInputs(t *testing.T) {
	int1 := g.Int(123)
	int2 := g.Int(456)

	hash1 := int1.Hash().MD5()
	hash2 := int2.Hash().MD5()

	// Different inputs should produce different hashes
	if hash1.Eq(hash2) {
		t.Error("Different inputs should produce different hashes")
	}
}

func TestInt_Hash_NegativeNumbers(t *testing.T) {
	positiveInt := g.Int(42)
	negativeInt := g.Int(-42)

	positiveHash := positiveInt.Hash().SHA256()
	negativeHash := negativeInt.Hash().SHA256()

	// Positive and negative versions should produce different hashes
	if positiveHash.Eq(negativeHash) {
		t.Error("Positive and negative integers should produce different hashes")
	}

	// Both should be valid hashes
	if positiveHash.Empty() || negativeHash.Empty() {
		t.Error("Hashes should not be empty")
	}
}

func TestInt_Hash_LargeNumbers(t *testing.T) {
	largeInt := g.Int(9223372036854775807) // max int64

	hash := largeInt.Hash().SHA512()

	if hash.Empty() {
		t.Error("Hash of large integer should not be empty")
	}

	// Should produce valid length hash
	if len(hash) != 128 {
		t.Errorf("SHA512 hash length should be 128, got %d", len(hash))
	}
}

func TestInt_Hash_AllAlgorithms_SameInput(t *testing.T) {
	testInt := g.Int(777)

	md5Hash := testInt.Hash().MD5()
	sha1Hash := testInt.Hash().SHA1()
	sha256Hash := testInt.Hash().SHA256()
	sha512Hash := testInt.Hash().SHA512()

	// All hashes should be different from each other
	hashes := []g.String{md5Hash, sha1Hash, sha256Hash, sha512Hash}

	for i, hash1 := range hashes {
		for j, hash2 := range hashes {
			if i != j && hash1.Eq(hash2) {
				t.Errorf("Different hash algorithms should produce different results")
			}
		}
	}
}
