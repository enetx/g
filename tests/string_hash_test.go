package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestStringMD5(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want String
	}{
		{
			name: "empty",
			s:    String("").Hash().MD5(),
			want: String("d41d8cd98f00b204e9800998ecf8427e"),
		},
		{
			name: "hello",
			s:    String("hello").Hash().MD5(),
			want: String("5d41402abc4b2a76b9719d911017c592"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s; got != tt.want {
				t.Errorf("String.MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSHA1(t *testing.T) {
	s := String("Hello, world!")
	expected := "943a702d06f34599aee1f8da8ef9f7296031d699"
	actual := s.Hash().SHA1().Std()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestStringSHA256(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want String
	}{
		{
			"empty",
			String(""),
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{"a", String("a"), "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
		{
			"abc",
			String("abc"),
			"ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			"message digest",
			String("message digest"),
			"f7846f55cf23e14eebeab5b4e1550cad5b509e3348fbc4efa3a1413d393cb650",
		},
		{
			"secure hash algorithm",
			String("secure hash algorithm"),
			"f30ceb2bb2829e79e4ca9753d35a8ecc00262d164cc077080295381cbd643f0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Hash().SHA256(); got != tt.want {
				t.Errorf("String.SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSHA512(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want String
	}{
		{
			"empty",
			String(""),
			"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			"hello",
			String("hello"),
			"9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			"hello world",
			String("hello world"),
			"309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Hash().SHA512(); got != tt.want {
				t.Errorf("String.SHA512() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringHMACSHA256(t *testing.T) {
	tests := []struct {
		name string
		data String
		key  String
		want String
	}{
		{
			"hello-secret",
			"hello",
			"secret",
			"88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b",
		},
		{
			"empty-data",
			"",
			"secret",
			"f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			"empty-key",
			"hello",
			"",
			"4352b26e33fe0d769a8922a6ba29004109f01688e26acc9e6cb347e5a5afc4da",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Hash().HMACSHA256(tt.key); got != tt.want {
				t.Errorf("String.HMACSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringHMACSHA512(t *testing.T) {
	tests := []struct {
		name string
		data String
		key  String
		want String
	}{
		{
			"hello-secret",
			"hello",
			"secret",
			"db1595ae88a62fd151ec1cba81b98c39df82daae7b4cb9820f446d5bf02f1dcfca6683d88cab3e273f5963ab8ec469a746b5b19086371239f67d1e5f99a79440",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Hash().HMACSHA512(tt.key); got != tt.want {
				t.Errorf("String.HMACSHA512() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringHashRawDigestLength(t *testing.T) {
	data := String("hello world")

	tests := []struct {
		name string
		raw  Bytes
		want int
	}{
		{"MD5Raw", data.Hash().MD5Raw(), 16},
		{"SHA1Raw", data.Hash().SHA1Raw(), 20},
		{"SHA256Raw", data.Hash().SHA256Raw(), 32},
		{"SHA512Raw", data.Hash().SHA512Raw(), 64},
		{"HMACSHA256Raw", data.Hash().HMACSHA256Raw("key"), 32},
		{"HMACSHA512Raw", data.Hash().HMACSHA512Raw("key"), 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.raw) != tt.want {
				t.Errorf("%s length = %d, want %d", tt.name, len(tt.raw), tt.want)
			}
		})
	}
}

func TestStringHashRawHexConsistency(t *testing.T) {
	data := String("hello world")
	key := String("secret")

	tests := []struct {
		name string
		hex  String
		raw  Bytes
	}{
		{"MD5", data.Hash().MD5(), data.Hash().MD5Raw()},
		{"SHA1", data.Hash().SHA1(), data.Hash().SHA1Raw()},
		{"SHA256", data.Hash().SHA256(), data.Hash().SHA256Raw()},
		{"SHA512", data.Hash().SHA512(), data.Hash().SHA512Raw()},
		{"HMACSHA256", data.Hash().HMACSHA256(key), data.Hash().HMACSHA256Raw(key)},
		{"HMACSHA512", data.Hash().HMACSHA512(key), data.Hash().HMACSHA512Raw(key)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.raw.Hex().StringUnsafe(); got != tt.hex {
				t.Errorf("%s() and %sRaw().Hex() mismatch\ngot:  %s\nwant: %s", tt.name, tt.name, got, tt.hex)
			}
		})
	}
}

func TestStringHMACDifferentKeys(t *testing.T) {
	data := String("hello world")

	hash1 := data.Hash().HMACSHA256("key1")
	hash2 := data.Hash().HMACSHA256("key2")

	if hash1 == hash2 {
		t.Error("HMAC with different keys should produce different results")
	}
}

func TestStringHMACConsistency(t *testing.T) {
	data := String("test data")
	key := String("key")

	hash1 := data.Hash().HMACSHA256(key)
	hash2 := data.Hash().HMACSHA256(key)

	if hash1 != hash2 {
		t.Error("HMAC results should be consistent for the same input and key")
	}
}
