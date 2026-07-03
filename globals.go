package g

import (
	"os"
)

const (
	// ASCII_LETTERS is the set of all ASCII letters (lowercase + uppercase).
	ASCII_LETTERS String = ASCII_LOWERCASE + ASCII_UPPERCASE
	// ASCII_LOWERCASE is the set of lowercase ASCII letters.
	ASCII_LOWERCASE String = "abcdefghijklmnopqrstuvwxyz"
	// ASCII_UPPERCASE is the set of uppercase ASCII letters.
	ASCII_UPPERCASE String = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// DIGITS is the set of decimal digit characters.
	DIGITS String = "0123456789"
	// HEXDIGITS is the set of hexadecimal digit characters (both cases).
	HEXDIGITS String = "0123456789abcdefABCDEF"
	// OCTDIGITS is the set of octal digit characters.
	OCTDIGITS String = "01234567"
	// PUNCTUATION is the set of ASCII punctuation characters.
	PUNCTUATION String = `!"#$%&'()*+,-./:;<=>?@[\]^{|}~` + "`"

	// FileDefault is the default permission mode (0o644) used when writing files.
	FileDefault os.FileMode = 0o644
	// FileCreate is the permission mode (0o666) used when creating files.
	FileCreate os.FileMode = 0o666
	// DirDefault is the default permission mode (0o755) used when creating directories.
	DirDefault os.FileMode = 0o755
	// FullAccess is the permission mode (0o777) granting read, write and execute to everyone.
	FullAccess os.FileMode = 0o777

	// PathSeparator is the OS-specific path separator as a String.
	PathSeparator = String(os.PathSeparator)
)
