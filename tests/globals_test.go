package g_test

import (
	"os"
	"testing"

	"github.com/enetx/g"
)

func TestConstants_ASCII(t *testing.T) {
	tests := []struct {
		name     string
		constant g.String
		expected string
		minLen   int
	}{
		{"ASCII_LOWERCASE", g.ASCII_LOWERCASE, "abcdefghijklmnopqrstuvwxyz", 26},
		{"ASCII_UPPERCASE", g.ASCII_UPPERCASE, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", 26},
		{"ASCII_LETTERS", g.ASCII_LETTERS, "", 52}, // combination of upper + lower
		{"DIGITS", g.DIGITS, "0123456789", 10},
		{"HEXDIGITS", g.HEXDIGITS, "", 22}, // 0-9 + a-f + A-F
		{"OCTDIGITS", g.OCTDIGITS, "01234567", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constant := string(tt.constant)

			if tt.expected != "" && constant != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, constant, tt.expected)
			}

			if len(constant) < tt.minLen {
				t.Errorf("%s length = %d, want at least %d", tt.name, len(constant), tt.minLen)
			}
		})
	}
}

func TestConstants_ASCII_LETTERS_Composition(t *testing.T) {
	expected := string(g.ASCII_LOWERCASE) + string(g.ASCII_UPPERCASE)
	actual := string(g.ASCII_LETTERS)

	if actual != expected {
		t.Errorf("ASCII_LETTERS should be ASCII_LOWERCASE + ASCII_UPPERCASE")
		t.Errorf("Got: %q", actual)
		t.Errorf("Want: %q", expected)
	}
}

func TestConstants_HEXDIGITS_Content(t *testing.T) {
	hexdigits := string(g.HEXDIGITS)

	// Should contain all digits 0-9
	for i := '0'; i <= '9'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HEXDIGITS should contain digit %c", i)
		}
	}

	// Should contain lowercase a-f
	for i := 'a'; i <= 'f'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HEXDIGITS should contain lowercase hex digit %c", i)
		}
	}

	// Should contain uppercase A-F
	for i := 'A'; i <= 'F'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HEXDIGITS should contain uppercase hex digit %c", i)
		}
	}
}

func TestConstants_PUNCTUATION(t *testing.T) {
	punctuation := string(g.PUNCTUATION)

	// Test that it contains expected punctuation characters
	expectedChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "=", "+"}

	for _, char := range expectedChars {
		if !contains(punctuation, char) {
			t.Errorf("PUNCTUATION should contain %q", char)
		}
	}

	// Should not be empty
	if len(punctuation) == 0 {
		t.Error("PUNCTUATION should not be empty")
	}
}

func TestConstants_FileModes(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected os.FileMode
	}{
		{"FileDefault", g.FileDefault, 0o644},
		{"FileCreate", g.FileCreate, 0o666},
		{"DirDefault", g.DirDefault, 0o755},
		{"FullAccess", g.FullAccess, 0o777},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode != tt.expected {
				t.Errorf("%s = %o, want %o", tt.name, tt.mode, tt.expected)
			}
		})
	}
}

func TestConstants_PathSeparator(t *testing.T) {
	expected := g.String(os.PathSeparator)

	if g.PathSeperator != expected {
		t.Errorf("PathSeperator = %q, want %q", g.PathSeperator, expected)
	}

	// Should not be empty
	if len(g.PathSeperator) == 0 {
		t.Error("PathSeperator should not be empty")
	}
}

func TestConstants_FileModes_Permissions(t *testing.T) {
	// Test that file modes have correct permission patterns

	// FileDefault (644) - owner read/write, group/other read
	if g.FileDefault&0o200 == 0 { // owner write
		t.Error("FileDefault should have owner write permission")
	}
	if g.FileDefault&0o400 == 0 { // owner read
		t.Error("FileDefault should have owner read permission")
	}
	if g.FileDefault&0o044 != 0o044 { // group/other read
		t.Error("FileDefault should have group and other read permissions")
	}

	// DirDefault (755) - owner read/write/execute, group/other read/execute
	if g.DirDefault&0o700 != 0o700 { // owner all permissions
		t.Error("DirDefault should have all owner permissions")
	}
	if g.DirDefault&0o055 != 0o055 { // group/other read/execute
		t.Error("DirDefault should have group and other read/execute permissions")
	}

	// FullAccess (777) - all permissions for all
	if g.FullAccess != 0o777 {
		t.Error("FullAccess should be 0o777")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
