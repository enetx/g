package g_test

import (
	"math"
	"testing"

	"github.com/enetx/g"
)

func TestFloat_Hash_MD5(t *testing.T) {
	testFloat := g.Float(3.14159)
	hash := testFloat.Hash().MD5()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("MD5 hash should not be empty")
	}

	// Hash should be 32 characters (16 bytes * 2 hex chars)
	if len(hash) != 32 {
		t.Errorf("MD5 hash length should be 32, got %d", len(hash))
	}
}

func TestFloat_Hash_SHA1(t *testing.T) {
	testFloat := g.Float(2.71828)
	hash := testFloat.Hash().SHA1()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA1 hash should not be empty")
	}

	// Hash should be 40 characters (20 bytes * 2 hex chars)
	if len(hash) != 40 {
		t.Errorf("SHA1 hash length should be 40, got %d", len(hash))
	}
}

func TestFloat_Hash_SHA256(t *testing.T) {
	testFloat := g.Float(1.41421356)
	hash := testFloat.Hash().SHA256()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA256 hash should not be empty")
	}

	// Hash should be 64 characters (32 bytes * 2 hex chars)
	if len(hash) != 64 {
		t.Errorf("SHA256 hash length should be 64, got %d", len(hash))
	}
}

func TestFloat_Hash_SHA512(t *testing.T) {
	testFloat := g.Float(1.73205080)
	hash := testFloat.Hash().SHA512()

	// Hash should not be empty
	if hash.Empty() {
		t.Error("SHA512 hash should not be empty")
	}

	// Hash should be 128 characters (64 bytes * 2 hex chars)
	if len(hash) != 128 {
		t.Errorf("SHA512 hash length should be 128, got %d", len(hash))
	}
}

func TestFloat_Hash_Zero(t *testing.T) {
	zeroFloat := g.Float(0.0)

	md5Hash := zeroFloat.Hash().MD5()
	sha1Hash := zeroFloat.Hash().SHA1()
	sha256Hash := zeroFloat.Hash().SHA256()
	sha512Hash := zeroFloat.Hash().SHA512()

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

func TestFloat_Hash_Consistency(t *testing.T) {
	testFloat := g.Float(42.42)

	// Hash the same float multiple times
	hash1 := testFloat.Hash().MD5()
	hash2 := testFloat.Hash().MD5()
	hash3 := testFloat.Hash().MD5()

	// Results should be consistent
	if !hash1.Eq(hash2) || !hash2.Eq(hash3) {
		t.Error("Hash results should be consistent for the same input")
	}
}

func TestFloat_Hash_DifferentInputs(t *testing.T) {
	float1 := g.Float(1.23)
	float2 := g.Float(4.56)

	hash1 := float1.Hash().MD5()
	hash2 := float2.Hash().MD5()

	// Different inputs should produce different hashes
	if hash1.Eq(hash2) {
		t.Error("Different inputs should produce different hashes")
	}
}

func TestFloat_Hash_NegativeNumbers(t *testing.T) {
	positiveFloat := g.Float(3.14)
	negativeFloat := g.Float(-3.14)

	positiveHash := positiveFloat.Hash().SHA256()
	negativeHash := negativeFloat.Hash().SHA256()

	// Positive and negative versions should produce different hashes
	if positiveHash.Eq(negativeHash) {
		t.Error("Positive and negative floats should produce different hashes")
	}

	// Both should be valid hashes
	if positiveHash.Empty() || negativeHash.Empty() {
		t.Error("Hashes should not be empty")
	}
}

func TestFloat_Hash_SmallDifferences(t *testing.T) {
	float1 := g.Float(1.0000001)
	float2 := g.Float(1.0000002)

	hash1 := float1.Hash().SHA256()
	hash2 := float2.Hash().SHA256()

	// Even small differences should produce different hashes
	if hash1.Eq(hash2) {
		t.Error("Small differences in floats should produce different hashes")
	}
}

func TestFloat_Hash_SpecialValues(t *testing.T) {
	tests := []struct {
		name  string
		value g.Float
	}{
		{"positive infinity", g.Float(math.Inf(1))},
		{"negative infinity", g.Float(math.Inf(-1))},
		{"NaN", g.Float(math.NaN())},
		{"very large", g.Float(1.7976931348623157e+308)},
		{"very small", g.Float(2.2250738585072014e-308)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := tt.value.Hash().MD5()

			// Should produce a hash without panicking
			if hash.Empty() {
				t.Errorf("Hash of %s should not be empty", tt.name)
			}

			if len(hash) != 32 {
				t.Errorf("Hash of %s should have length 32, got %d", tt.name, len(hash))
			}
		})
	}
}

func TestFloat_Hash_AllAlgorithms_SameInput(t *testing.T) {
	testFloat := g.Float(123.456)

	md5Hash := testFloat.Hash().MD5()
	sha1Hash := testFloat.Hash().SHA1()
	sha256Hash := testFloat.Hash().SHA256()
	sha512Hash := testFloat.Hash().SHA512()

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
