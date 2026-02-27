package g

import (
	"fmt"
	"math"
	"strconv"
)

type fmtAlign int

const (
	alignNone   fmtAlign = iota
	alignLeft            // '<'
	alignRight           // '>'
	alignCenter          // '^'
)

type fmtSign int

const (
	signNone  fmtSign = iota
	signPlus          // '+'
	signMinus         // '-'
	signSpace         // ' '
)

type fmtSpec struct {
	fill      rune
	align     fmtAlign
	sign      fmtSign
	alternate bool
	zeroPad   bool
	width     Int
	precision Int
	verb      byte
}

func splitFmtSpec(placeholder String) (String, String) {
	if i := findUnparenColon(placeholder); i >= 0 {
		return placeholder[:i], placeholder[i+1:]
	}

	return placeholder, ""
}

func findUnparenColon(placeholder String) Int {
	depth := 0

	for i := Int(0); i < placeholder.Len(); i++ {
		switch placeholder[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ':':
			if depth == 0 {
				return i
			}
		}
	}

	return -1
}

func parseDigits(runes []rune, pos int) (int, int, bool) {
	n, has := 0, false

	for pos < len(runes) && runes[pos] >= '0' && runes[pos] <= '9' {
		n = n*10 + int(runes[pos]-'0')
		pos++
		has = true
	}

	return n, pos, has
}

func parseFmtSpec(spec String) fmtSpec {
	fs := fmtSpec{fill: ' ', width: -1, precision: -1}

	runes := []rune(spec)
	pos := 0

	// parse [[fill]align]
	if pos < len(runes) {
		if pos+1 < len(runes) && isAlign(runes[pos+1]) {
			fs.fill = runes[pos]
			fs.align = toAlign(runes[pos+1])
			pos += 2
		} else if isAlign(runes[pos]) {
			fs.align = toAlign(runes[pos])
			pos++
		}
	}

	// parse [sign]
	if pos < len(runes) {
		switch runes[pos] {
		case '+':
			fs.sign = signPlus
			pos++
		case '-':
			fs.sign = signMinus
			pos++
		case ' ':
			fs.sign = signSpace
			pos++
		}
	}

	// parse [#]
	if pos < len(runes) && runes[pos] == '#' {
		fs.alternate = true
		pos++
	}

	// parse [0]
	if pos < len(runes) && runes[pos] == '0' && fs.align == alignNone {
		fs.zeroPad = true
		pos++
	}

	// parse [width]
	if n, np, ok := parseDigits(runes, pos); ok {
		fs.width = Int(n)
		pos = np
	}

	// parse [.precision]
	if pos < len(runes) && runes[pos] == '.' {
		pos++

		if n, np, ok := parseDigits(runes, pos); ok {
			fs.precision = Int(n)
			pos = np
		} else {
			fs.precision = 0
		}
	}

	// parse [type]
	if pos < len(runes) {
		fs.verb = byte(runes[pos])
	}

	return fs
}

func isAlign(r rune) bool { return r == '<' || r == '>' || r == '^' }

func toAlign(r rune) fmtAlign {
	switch r {
	case '<':
		return alignLeft
	case '>':
		return alignRight
	default:
		return alignCenter
	}
}

func applyFmtSpec(value any, spec fmtSpec) String {
	var s String

	switch spec.verb {
	case 'x':
		s = fmtIntBase(value, 16, fmtAltPrefix("0x", spec), spec, false)
	case 'X':
		s = fmtIntBase(value, 16, fmtAltPrefix("0x", spec), spec, true)
	case 'o':
		s = fmtIntBase(value, 8, fmtAltPrefix("0o", spec), spec, false)
	case 'b':
		s = fmtIntBase(value, 2, fmtAltPrefix("0b", spec), spec, false)
	case 'e':
		s = fmtExponential(value, spec, 'e')
	case 'E':
		s = fmtExponential(value, spec, 'E')
	case 'T':
		s = String(fmt.Sprintf("%T", value))
	case 'p':
		s = String(fmt.Sprintf("%p", value))
	case '?':
		if spec.alternate {
			s = fmtPrettyDebug(value)
		} else {
			s = String(fmt.Sprintf("%#v", value))
		}
	default:
		s = fmtDefault(value, spec)
	}

	return fmtPad(s, spec, fmtIsNumeric(value))
}

func fmtAltPrefix(prefix String, spec fmtSpec) String {
	if spec.alternate {
		return prefix
	}

	return ""
}

func fmtIntBase(value any, base int, prefix String, spec fmtSpec, upper bool) String {
	format := func(s String) String {
		if upper {
			return prefix + s.Upper()
		}

		return prefix + s
	}

	if u, ok := fmtToUint64(value); ok {
		return format(String(strconv.FormatUint(u, base)))
	}

	if i, ok := fmtToInt64(value); ok {
		s := format(String(strconv.FormatInt(abs64(i), base)))

		if i < 0 {
			return "-" + s
		}

		return fmtApplySign(s, true, spec)
	}

	return String(fmt.Sprint(value))
}

func fmtDefault(value any, spec fmtSpec) String {
	prec := -1
	if spec.precision >= 0 {
		prec = spec.precision.Std()
	}

	if prec >= 0 || spec.sign != signNone {
		if f, ok := fmtToFloat64(value); ok {
			return fmtApplySign(
				String(strconv.FormatFloat(f, 'f', prec, 64)), f >= 0, spec)
		}

		if spec.sign != signNone {
			if i, ok := fmtToInt64(value); ok {
				return fmtApplySign(
					String(strconv.FormatInt(abs64(i), 10)), i >= 0, spec)
			}
		}
	}

	s := String(fmt.Sprint(value))
	if prec >= 0 && s.LenRunes() > spec.precision {
		return String([]rune(s)[:spec.precision])
	}

	return s
}

func fmtExponential(value any, spec fmtSpec, verb byte) String {
	prec := 6
	if spec.precision >= 0 {
		prec = spec.precision.Std()
	}

	f, ok := fmtToFloat64(value)
	if !ok {
		var i int64
		if i, ok = fmtToInt64(value); ok {
			f = float64(i)
		}
	}

	if ok {
		return fmtApplySign(
			String(strconv.FormatFloat(f, verb, prec, 64)), f >= 0, spec)
	}

	return String(fmt.Sprint(value))
}

func fmtApplySign(s String, positive bool, spec fmtSpec) String {
	if !positive {
		if s[0] != '-' {
			return "-" + s
		}

		return s
	}

	switch spec.sign {
	case signPlus:
		return "+" + s
	case signSpace:
		return " " + s
	}

	return s
}

func fmtPad(s String, spec fmtSpec, isNum bool) String {
	if spec.width.IsNegative() || s.LenRunes() >= spec.width {
		return s
	}

	fill := String(string(spec.fill))

	// zero-pad: sign-aware
	if spec.zeroPad && spec.align == alignNone && isNum {
		return fmtZeroPad(s, spec.width)
	}

	// default alignment: right for numbers, left for strings
	align := spec.align
	if align == alignNone {
		if isNum {
			align = alignRight
		} else {
			align = alignLeft
		}
	}

	switch align {
	case alignLeft:
		return s.LeftJustify(spec.width, fill)
	case alignRight:
		return s.RightJustify(spec.width, fill)
	default:
		return s.Center(spec.width, fill)
	}
}

func fmtZeroPad(s String, width Int) String {
	var sign byte

	body := s
	if s.Len() > 0 && (s[0] == '+' || s[0] == '-' || s[0] == ' ') {
		sign = s[0]
		body = s[1:]
	}

	var prefix String
	if body.Len() > 1 && body[0] == '0' && (body[1] == 'x' || body[1] == 'X' || body[1] == 'b' || body[1] == 'o') {
		prefix = body[:2]
		body = body[2:]
	}

	padLen := width - body.LenRunes() - prefix.Len()
	if sign != 0 {
		padLen--
	}

	var b Builder
	if sign != 0 {
		b.WriteByte(sign)
	}

	b.WriteString(prefix)
	b.WriteString(String("0").Repeat(padLen))
	b.WriteString(body)

	return b.String()
}

func fmtPrettyDebug(value any) String {
	return fmtIndentDebug(String(fmt.Sprintf("%#v", value)))
}

func fmtIndentDebug(s String) String {
	var b Builder

	indent := Int(0)
	inString := false
	pad := String("  ")

	for i := Int(0); i < s.Len(); i++ {
		ch := s[i]

		if ch == '"' && (i == 0 || s[i-1] != '\\') {
			inString = !inString
			b.WriteByte(ch)

			continue
		}

		if inString {
			b.WriteByte(ch)
			continue
		}

		switch ch {
		case '{', '[':
			indent++
			b.WriteByte(ch)
			b.WriteByte('\n')
			b.WriteString(pad.Repeat(indent))
		case '}', ']':
			indent--
			b.WriteByte('\n')
			b.WriteString(pad.Repeat(indent))
			b.WriteByte(ch)
		case ',':
			b.WriteByte(ch)
			b.WriteByte('\n')
			b.WriteString(pad.Repeat(indent))
		default:
			b.WriteByte(ch)
		}
	}

	return b.String()
}

func fmtIsNumeric(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		Int, Float:
		return true
	}

	return false
}

func fmtToInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case Int:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v <= math.MaxInt64 {
			return int64(v), true
		}

		return 0, false
	case Float:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	}

	return 0, false
}

func fmtToUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	}

	return 0, false
}

func fmtToFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case Float:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}

	return 0, false
}

func abs64(i int64) int64 {
	if i < 0 {
		return -i
	}

	return i
}
