package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestBytes_Hash_MD5(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().MD5()
	if hash.IsEmpty() {
		t.Error("MD5 hash should not be empty")
	}
	if len(hash) != 32 {
		t.Errorf("MD5 hash length should be 32, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA1(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA1()
	if hash.IsEmpty() {
		t.Error("SHA1 hash should not be empty")
	}
	if len(hash) != 40 {
		t.Errorf("SHA1 hash length should be 40, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA256(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA256()
	if hash.IsEmpty() {
		t.Error("SHA256 hash should not be empty")
	}
	if len(hash) != 64 {
		t.Errorf("SHA256 hash length should be 64, got %d", len(hash))
	}
}

func TestBytes_Hash_SHA512(t *testing.T) {
	testBytes := g.Bytes("hello world")
	hash := testBytes.Hash().SHA512()
	if hash.IsEmpty() {
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

func TestBytes_Hash_MD5Raw(t *testing.T) {
	raw := g.Bytes("hello world").Hash().MD5Raw()
	if len(raw) != 16 {
		t.Errorf("MD5 raw digest length should be 16, got %d", len(raw))
	}
}

func TestBytes_Hash_SHA1Raw(t *testing.T) {
	raw := g.Bytes("hello world").Hash().SHA1Raw()
	if len(raw) != 20 {
		t.Errorf("SHA1 raw digest length should be 20, got %d", len(raw))
	}
}

func TestBytes_Hash_SHA256Raw(t *testing.T) {
	raw := g.Bytes("hello world").Hash().SHA256Raw()
	if len(raw) != 32 {
		t.Errorf("SHA256 raw digest length should be 32, got %d", len(raw))
	}
}

func TestBytes_Hash_SHA512Raw(t *testing.T) {
	raw := g.Bytes("hello world").Hash().SHA512Raw()
	if len(raw) != 64 {
		t.Errorf("SHA512 raw digest length should be 64, got %d", len(raw))
	}
}

func TestBytes_Hash_RawHexConsistency(t *testing.T) {
	data := g.Bytes("hello world")

	hex := data.Hash().SHA256()
	rawHex := data.Hash().SHA256Raw().Hex()

	if !hex.Eq(rawHex) {
		t.Errorf("SHA256() and SHA256Raw().Hex() should be equal\ngot:  %s\nwant: %s", rawHex, hex)
	}
}

func TestBytes_Hash_HMACSHA256(t *testing.T) {
	data := g.Bytes("hello world")
	key := g.Bytes("secret")
	hash := data.Hash().HMACSHA256(key)

	if hash.IsEmpty() {
		t.Error("HMACSHA256 hash should not be empty")
	}
	if len(hash) != 64 {
		t.Errorf("HMACSHA256 hex length should be 64, got %d", len(hash))
	}
}

func TestBytes_Hash_HMACSHA512(t *testing.T) {
	data := g.Bytes("hello world")
	key := g.Bytes("secret")
	hash := data.Hash().HMACSHA512(key)

	if hash.IsEmpty() {
		t.Error("HMACSHA512 hash should not be empty")
	}
	if len(hash) != 128 {
		t.Errorf("HMACSHA512 hex length should be 128, got %d", len(hash))
	}
}

func TestBytes_Hash_HMACSHA256Raw(t *testing.T) {
	data := g.Bytes("hello world")
	key := g.Bytes("secret")
	raw := data.Hash().HMACSHA256Raw(key)

	if len(raw) != 32 {
		t.Errorf("HMACSHA256 raw digest length should be 32, got %d", len(raw))
	}
}

func TestBytes_Hash_HMACSHA512Raw(t *testing.T) {
	data := g.Bytes("hello world")
	key := g.Bytes("secret")
	raw := data.Hash().HMACSHA512Raw(key)

	if len(raw) != 64 {
		t.Errorf("HMACSHA512 raw digest length should be 64, got %d", len(raw))
	}
}

func TestBytes_Hash_HMACRawHexConsistency(t *testing.T) {
	data := g.Bytes("hello world")
	key := g.Bytes("secret")

	hex := data.Hash().HMACSHA256(key)
	rawHex := data.Hash().HMACSHA256Raw(key).Hex()

	if !hex.Eq(rawHex) {
		t.Errorf("HMACSHA256() and HMACSHA256Raw().Hex() should be equal\ngot:  %s\nwant: %s", rawHex, hex)
	}
}

func TestBytes_Hash_HMACConsistency(t *testing.T) {
	data := g.Bytes("test data")
	key := g.Bytes("key")

	hash1 := data.Hash().HMACSHA256(key)
	hash2 := data.Hash().HMACSHA256(key)

	if !hash1.Eq(hash2) {
		t.Error("HMAC results should be consistent for the same input and key")
	}
}

func TestBytes_Hash_HMACDifferentKeys(t *testing.T) {
	data := g.Bytes("hello world")

	hash1 := data.Hash().HMACSHA256(g.Bytes("key1"))
	hash2 := data.Hash().HMACSHA256(g.Bytes("key2"))

	if hash1.Eq(hash2) {
		t.Error("HMAC with different keys should produce different results")
	}
}

func TestBytes_Hex(t *testing.T) {
	data := g.Bytes("hello")
	hex := data.Hex()

	expected := g.Bytes("68656c6c6f")
	if !hex.Eq(expected) {
		t.Errorf("Hex() should be %s, got %s", expected, hex)
	}
}

func TestBytes_Hex_Empty(t *testing.T) {
	data := g.Bytes("")
	hex := data.Hex()

	if !hex.IsEmpty() {
		t.Errorf("Hex() of empty Bytes should be empty, got %s", hex)
	}
}

func TestBytes_Hash_KnownSHA256(t *testing.T) {
	// SHA256("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	hash := g.Bytes("hello").Hash().SHA256()
	expected := g.Bytes("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824")

	if !hash.Eq(expected) {
		t.Errorf("SHA256 known value mismatch\ngot:  %s\nwant: %s", hash, expected)
	}
}

func TestBytes_Hash_KnownHMACSHA256(t *testing.T) {
	// HMAC-SHA256("hello", "secret") = 88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b
	hash := g.Bytes("hello").Hash().HMACSHA256(g.Bytes("secret"))
	expected := g.Bytes("88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b")

	if !hash.Eq(expected) {
		t.Errorf("HMACSHA256 known value mismatch\ngot:  %s\nwant: %s", hash, expected)
	}
}
