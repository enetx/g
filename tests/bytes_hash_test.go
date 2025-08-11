package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestBytes_Hash_MD5(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().MD5()

	if hash.Empty() {
		t.Error("MD5 hash should not be empty")
	}

	if len(hash) != 32 {
		t.Errorf("MD5 hash length should be 32, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA1(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA1()

	if hash.Empty() {
		t.Error("SHA1 hash should not be empty")
	}

	if len(hash) != 40 {
		t.Errorf("SHA1 hash length should be 40, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA256(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA256()

	if hash.Empty() {
		t.Error("SHA256 hash should not be empty")
	}

	if len(hash) != 64 {
		t.Errorf("SHA256 hash length should be 64, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA512(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA512()

	if hash.Empty() {
		t.Error("SHA512 hash should not be empty")
	}

	if len(hash) != 128 {
		t.Errorf("SHA512 hash length should be 128, got %d", len(hash))
	}
}

func TestBytes_Hash_Consistency(t *testing.T) {
	testBytes := g.Bytes("test data")

	hash1 := testBytes.Hash().MD5()
	hash2 := testBytes.Hash().MD5()

	if !hash1.Eq(hash2) {
		t.Error("Hash results should be consistent for the same input")
	}
}
