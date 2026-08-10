package g_test

import (
	"os"
	"testing"

	. "github.com/enetx/g"
)

func TestConstants_ASCII(t *testing.T) {
	tests := []struct {
		name     string
		constant String
		expected string
		minLen   int
	}{
		{"ASCIILowercase", ASCIILowercase, "abcdefghijklmnopqrstuvwxyz", 26},
		{"ASCIIUppercase", ASCIIUppercase, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", 26},
		{"ASCIILetters", ASCIILetters, "", 52}, // combination of upper + lower
		{"Digits", Digits, "0123456789", 10},
		{"HexDigits", HexDigits, "", 22}, // 0-9 + a-f + A-F
		{"OctDigits", OctDigits, "01234567", 8},
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

func TestConstants_ASCIILetters_Composition(t *testing.T) {
	expected := string(ASCIILowercase) + string(ASCIIUppercase)
	actual := string(ASCIILetters)

	if actual != expected {
		t.Errorf("ASCIILetters should be ASCIILowercase + ASCIIUppercase")
		t.Errorf("Got: %q", actual)
		t.Errorf("Want: %q", expected)
	}
}

func TestConstants_HEXDIGITS_Content(t *testing.T) {
	hexdigits := string(HexDigits)

	// Should contain all digits 0-9
	for i := '0'; i <= '9'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HexDigits should contain digit %c", i)
		}
	}

	// Should contain lowercase a-f
	for i := 'a'; i <= 'f'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HexDigits should contain lowercase hex digit %c", i)
		}
	}

	// Should contain uppercase A-F
	for i := 'A'; i <= 'F'; i++ {
		if !contains(hexdigits, string(i)) {
			t.Errorf("HexDigits should contain uppercase hex digit %c", i)
		}
	}
}

func TestConstants_PUNCTUATION(t *testing.T) {
	punctuation := string(Punctuation)

	// Test that it contains expected punctuation characters
	expectedChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "=", "+"}

	for _, char := range expectedChars {
		if !contains(punctuation, char) {
			t.Errorf("Punctuation should contain %q", char)
		}
	}

	// Should not be empty
	if len(punctuation) == 0 {
		t.Error("Punctuation should not be empty")
	}
}

func TestConstants_FileModes(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected os.FileMode
	}{
		{"FileDefault", FileDefault, 0o644},
		{"FileCreate", FileCreate, 0o666},
		{"DirDefault", DirDefault, 0o755},
		{"FullAccess", FullAccess, 0o777},
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
	expected := String(os.PathSeparator)

	if PathSeparator != expected {
		t.Errorf("PathSeperator = %q, want %q", PathSeparator, expected)
	}

	// Should not be empty
	if len(PathSeparator) == 0 {
		t.Error("PathSeperator should not be empty")
	}
}

func TestConstants_FileModes_Permissions(t *testing.T) {
	// Test that file modes have correct permission patterns

	// FileDefault (644) - owner read/write, group/other read
	if FileDefault&0o200 == 0 { // owner write
		t.Error("FileDefault should have owner write permission")
	}
	if FileDefault&0o400 == 0 { // owner read
		t.Error("FileDefault should have owner read permission")
	}
	if FileDefault&0o044 != 0o044 { // group/other read
		t.Error("FileDefault should have group and other read permissions")
	}

	// DirDefault (755) - owner read/write/execute, group/other read/execute
	if DirDefault&0o700 != 0o700 { // owner all permissions
		t.Error("DirDefault should have all owner permissions")
	}
	if DirDefault&0o055 != 0o055 { // group/other read/execute
		t.Error("DirDefault should have group and other read/execute permissions")
	}

	// FullAccess (777) - all permissions for all
	if FullAccess != 0o777 {
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
