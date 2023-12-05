package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
)

func TestStringMD5(t *testing.T) {
	tests := []struct {
		name string
		s    g.String
		want g.String
	}{
		{
			name: "empty",
			s:    g.NewString("").Hash().MD5(),
			want: g.String("d41d8cd98f00b204e9800998ecf8427e"),
		},
		{
			name: "hello",
			s:    g.NewString("hello").Hash().MD5(),
			want: g.String("5d41402abc4b2a76b9719d911017c592"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s; got != tt.want {
				t.Errorf("g.String.MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSHA1(t *testing.T) {
	s := g.NewString("Hello, world!")
	expected := "943a702d06f34599aee1f8da8ef9f7296031d699"

	actual := s.Hash().SHA1().Std()
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestStringSHA256(t *testing.T) {
	tests := []struct {
		name string
		s    g.String
		want g.String
	}{
		{
			"empty",
			g.String(""),
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{"a", g.String("a"), "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
		{
			"abc",
			g.String("abc"),
			"ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			"message digest",
			g.String("message digest"),
			"f7846f55cf23e14eebeab5b4e1550cad5b509e3348fbc4efa3a1413d393cb650",
		},
		{
			"secure hash algorithm",
			g.String("secure hash algorithm"),
			"f30ceb2bb2829e79e4ca9753d35a8ecc00262d164cc077080295381cbd643f0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Hash().SHA256(); got != tt.want {
				t.Errorf("g.String.SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSHA512(t *testing.T) {
	tests := []struct {
		name string
		s    g.String
		want g.String
	}{
		{
			"empty",
			g.String(""),
			"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			"hello",
			g.String("hello"),
			"9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			"hello world",
			g.String("hello world"),
			"309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Hash().SHA512(); got != tt.want {
				t.Errorf("g.String.SHA512() = %v, want %v", got, tt.want)
			}
		})
	}
}
