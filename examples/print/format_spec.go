package main

import (
	. "github.com/enetx/g"
)

func main() {
	// --- Format specifiers use {:spec} syntax (Rust-style) ---
	// The colon separates the placeholder from the format spec.

	// --- Integer bases ---
	// {:x} hex lowercase, {:X} hex uppercase, {:o} octal, {:b} binary
	Println("hex:  {:x}", 255) // ff
	Println("HEX:  {:X}", 255) // FF
	Println("oct:  {:o}", 255) // 377
	Println("bin:  {:b}", 42)  // 101010

	// --- Alternate form (#) adds base prefixes ---
	Println("alt hex: {:#x}", 255) // 0xff
	Println("alt HEX: {:#X}", 255) // 0xFF
	Println("alt oct: {:#o}", 255) // 0o377
	Println("alt bin: {:#b}", 42)  // 0b101010

	// --- Width and alignment ---
	// > right-align, < left-align, ^ center
	Println("[{:>10}]", "hello") // [     hello]
	Println("[{:<10}]", "hello") // [hello     ]
	Println("[{:^10}]", "hello") // [  hello   ]

	// --- Custom fill character ---
	// Any character before the alignment specifier becomes the fill.
	Println("[{:-^20}]", "hello") // [-------hello--------]
	Println("[{:*>10}]", "hi")    // [********hi]
	Println("[{:.<10}]", "hi")    // [hi........]

	// --- Default alignment ---
	// Numbers default to right-align, strings default to left-align.
	Println("[{:10}]", 42)   // [        42]
	Println("[{:10}]", "hi") // [hi        ]

	// --- Precision ---
	// For floats: number of decimal places. For strings: truncate.
	Println("{:.2}", 3.14159) // 3.14
	Println("{:.0}", 3.14159) // 3
	Println("{:.3}", "hello") // hel

	// --- Sign ---
	// + always show sign, ' ' (space) adds space for positive numbers.
	Println("{:+}", 42)  // +42
	Println("{:+}", -42) // -42
	Println("{: }", 42)  //  42

	// --- Zero-padding ---
	// 0 before width pads with zeros. Sign-aware: sign stays in front.
	Println("{:05}", 42)  // 00042
	Println("{:+05}", 42) // +0042
	Println("{:05}", -42) // -0042

	// --- Combining specifiers ---
	// Zero-pad + alternate hex
	Println("{:#010x}", 255) // 0x000000ff
	// Zero-pad binary
	Println("{:08b}", 42) // 00101010
	// Width + precision + alignment
	Println("[{:>10.2}]", 3.14159) // [      3.14]
	// Sign + precision
	Println("{:+.2}", 3.14) // +3.14

	// --- Exponential notation ---
	Println("{:e}", 3.14)      // 3.140000e+00
	Println("{:E}", 3.14)      // 3.140000E+00
	Println("{:.2e}", 3.14159) // 3.14e+00

	// --- Debug ---
	// {:?} uses %#v, {:#?} pretty-prints with indentation.
	Println("{:?}", "hello") // "hello"
	Println("{:?}", 42)      // 42

	// --- Works with g.Int and g.Float ---
	Println("{:x}", Int(255))        // ff
	Println("{:08b}", Int(42))       // 00101010
	Println("{:.2}", Float(3.14159)) // 3.14
	Println("{:05}", Int(42))        // 00042

	// --- Modifiers + format spec ---
	// Modifiers (method calls) run first, then the spec formats the result.
	Println("{1.Abs:05}", Int(-42))          // 00042
	Println("{.Upper:>10}", String("hello")) // "     HELLO"

	// --- Named + spec ---
	data := Named{"port": 8080, "name": String("go")}
	Println("port: {port:x}", data)         // port: 1f90
	Println("name: {name.Upper:>10}", data) // name:         GO

	// --- Auto-index + spec ---
	Println("{:x} {:o} {:b}", 255, 255, 255) // ff 377 11111111
}
