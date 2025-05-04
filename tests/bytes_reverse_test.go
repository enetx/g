package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

var reverseTests = []struct {
	name  string
	input Bytes
	want  Bytes
}{
	{"Empty", Bytes(""), Bytes("")},
	{"Single ASCII", Bytes("A"), Bytes("A")},
	{"ASCII", Bytes("ABCdef"), Bytes("fedCBA")},
	{"Single Unicode", Bytes("Ж"), Bytes("Ж")},
	{"Unicode", Bytes("Привет"), Bytes("тевирП")},
	{"Chinese", Bytes("你好世界"), Bytes("界世好你")},
	{"Hindi", Bytes("नमस्ते"), Bytes("ेत्समन")},
	{"Mixed ASCII+Unicode", Bytes("Go🚀Lang"), Bytes("gnaL🚀oG")},
	{"Family Emoji", Bytes("👨‍👩‍👧‍👦"), Bytes("👦‍👧‍👩‍👨")},
	{"Combining", Bytes("é"), Bytes("é")},
	{"Emoji Sequence", Bytes("🙂🙃🙂"), Bytes("🙂🙃🙂")},
	{"Raw Bytes", Bytes([]byte{0, 1, 2, 3}), Bytes([]byte{3, 2, 1, 0})},
	{"Variation Selector", Bytes("✈️"), Bytes("️✈")},
	{"Invalid UTF-8 Bytes", Bytes([]byte{0xff, 0xfe, 0xfd}), Bytes([]byte{0xfd, 0xfe, 0xff})},
}

func TestReverse(t *testing.T) {
	for _, tt := range reverseTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reverse()
			if string(got) != string(tt.want) {
				t.Errorf("Reverse(%q) = %q; want %q", string(tt.input), string(got), string(tt.want))
			}
		})
	}
}
