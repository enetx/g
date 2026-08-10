package mimesniff_test

import (
	"bytes"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/enetx/g/internal/mimesniff"
)

// The sniffer here is a verbatim copy of the standard library's, kept out of
// the production closure so that the fs package does not import net/http. Only
// a test may import net/http, so the equivalence is asserted here rather than
// assumed.
//
// If this ever fails, the copy has drifted from the toolchain: re-copy
// sniff.go from GOROOT instead of patching it.
func TestDetectContentTypeMatchesNetHTTP(t *testing.T) {
	cases := map[string][]byte{
		"empty":            {},
		"nil":              nil,
		"ascii text":       []byte("test content"),
		"leading space":    []byte("   \n\t hello"),
		"html document":    []byte("<!DOCTYPE HTML><html><body>hi</body></html>"),
		"html fragment":    []byte("<HTML> <head>"),
		"xml":              []byte(`<?xml version="1.0" encoding="UTF-8"?><a/>`),
		"pdf":              []byte("%PDF-1.7\n%\xe2\xe3\xcf\xd3"),
		"postscript":       []byte("%!PS-Adobe-3.0"),
		"png":              append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 32)...),
		"jpeg":             append([]byte("\xff\xd8\xff"), make([]byte, 32)...),
		"gif87":            []byte("GIF87a....."),
		"gif89":            []byte("GIF89a....."),
		"webp":             []byte("RIFF\x00\x00\x00\x00WEBPVP8 "),
		"wav":              []byte("RIFF\x00\x00\x00\x00WAVEfmt "),
		"mp4":              []byte("\x00\x00\x00\x18ftypmp42\x00\x00\x00\x00mp42isom"),
		"gzip":             []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00"),
		"zip":              []byte("PK\x03\x04\x14\x00\x00\x00"),
		"utf8 bom":         []byte("\xef\xbb\xbfhello"),
		"utf16 le bom":     []byte("\xff\xfeh\x00i\x00"),
		"utf16 be bom":     []byte("\xfe\xffh\x00i\x00"),
		"binary":           {0x00, 0x01, 0x02, 0x03, 0x7f, 0xff},
		"ogg":              []byte("OggS\x00some audio payload"),
		"woff":             []byte("wOFF\x00\x01\x00\x00"),
		"truncated prefix": []byte("\x89PN"),
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			want := http.DetectContentType(data)
			if got := mimesniff.DetectContentType(data); got != want {
				t.Fatalf("DetectContentType(%q) = %q, want %q", data, got, want)
			}
		})
	}
}

// Beyond the curated table, the copy must agree on arbitrary bytes: the
// sniffing algorithm branches on byte patterns the table cannot enumerate.
func TestDetectContentTypeMatchesNetHTTPOnRandomInput(t *testing.T) {
	source := rand.New(rand.NewPCG(0x5EED, 0xC0FFEE))
	buffer := make([]byte, mimesniff.SniffLen+64)

	for range 2_000 {
		size := source.IntN(len(buffer))
		data := buffer[:size]
		for i := range data {
			data[i] = byte(source.UintN(256))
		}

		want := http.DetectContentType(data)
		if got := mimesniff.DetectContentType(data); got != want {
			t.Fatalf("DetectContentType(%x) = %q, want %q", data, got, want)
		}
	}
}

// The prefix bound is part of the contract: net/http considers at most the
// first 512 bytes, so trailing content must never change the verdict.
func TestSniffLenMatchesNetHTTPPrefixBound(t *testing.T) {
	head := append([]byte("<!DOCTYPE HTML>"), bytes.Repeat([]byte(" "), mimesniff.SniffLen)...)
	tail := append(append([]byte{}, head...), 0x00, 0x01, 0x02)

	if got, want := mimesniff.DetectContentType(tail), http.DetectContentType(tail); got != want {
		t.Fatalf("beyond-prefix bytes changed the verdict: %q != %q", got, want)
	}
	if got, want := mimesniff.DetectContentType(tail), mimesniff.DetectContentType(head); got != want {
		t.Fatalf("copy is not bounded by SniffLen: %q != %q", got, want)
	}
}
