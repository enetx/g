package g_test

import (
	"math"
	"testing"

	. "github.com/enetx/g"
)

func TestFloatMathClassify(t *testing.T) {
	nan := Float(math.NaN())
	pinf := Float(math.Inf(1))
	ninf := Float(math.Inf(-1))
	subnormal := Float(5e-324)
	zero := Float(0)
	negzero := Float(math.Copysign(0, -1))

	if !nan.IsNaN() {
		t.Error("IsNaN(NaN): expected true")
	}

	if Float(1.5).IsNaN() {
		t.Error("IsNaN(1.5): expected false")
	}

	if nan.IsInf() || nan.IsFinite() || nan.IsNormal() {
		t.Error("NaN: expected IsInf, IsFinite, IsNormal all false")
	}

	if !pinf.IsInf() || !ninf.IsInf() {
		t.Error("IsInf(±Inf): expected true")
	}

	if pinf.IsNaN() || pinf.IsFinite() || pinf.IsNormal() {
		t.Error("+Inf: expected IsNaN, IsFinite, IsNormal all false")
	}

	if ninf.IsNaN() || ninf.IsFinite() || ninf.IsNormal() {
		t.Error("-Inf: expected IsNaN, IsFinite, IsNormal all false")
	}

	if !Float(1.5).IsFinite() || !zero.IsFinite() || !subnormal.IsFinite() {
		t.Error("IsFinite(1.5, 0, subnormal): expected true")
	}

	if !Float(1.5).IsNormal() || !Float(-2.75).IsNormal() {
		t.Error("IsNormal(1.5, -2.75): expected true")
	}

	if zero.IsNormal() {
		t.Error("IsNormal(0): expected false")
	}

	if negzero.IsNormal() {
		t.Error("IsNormal(-0): expected false")
	}

	if subnormal.IsNormal() {
		t.Error("IsNormal(5e-324): expected false, subnormal is not normal")
	}

	if !Float(math.SmallestNonzeroFloat64).IsFinite() || Float(math.SmallestNonzeroFloat64).IsNormal() {
		t.Error("SmallestNonzeroFloat64: expected finite but not normal")
	}
}

func TestFloatMathSignum(t *testing.T) {
	if !Float(math.NaN()).Signum().IsNaN() {
		t.Error("Signum(NaN): expected NaN")
	}

	if Float(3.5).Signum() != 1 {
		t.Errorf("Signum(3.5): expected 1, got %v", Float(3.5).Signum())
	}

	if Float(-2).Signum() != -1 {
		t.Errorf("Signum(-2): expected -1, got %v", Float(-2).Signum())
	}

	if Float(0).Signum() != 1 {
		t.Errorf("Signum(+0): expected 1, got %v", Float(0).Signum())
	}

	negzero := Float(math.Copysign(0, -1))
	if negzero.Signum() != -1 {
		t.Errorf("Signum(-0): expected -1, got %v", negzero.Signum())
	}

	if Float(math.Inf(1)).Signum() != 1 {
		t.Error("Signum(+Inf): expected 1")
	}

	if Float(math.Inf(-1)).Signum() != -1 {
		t.Error("Signum(-Inf): expected -1")
	}
}

func TestFloatIsSignPositiveNegative(t *testing.T) {
	tests := []struct {
		name         string
		in           Float
		signPositive bool
	}{
		{"positive", Float(3.5), true},
		{"negative", Float(-2), false},
		{"+0", Float(0), true},
		{"-0", Float(math.Copysign(0, -1)), false},
		{"+Inf", Float(math.Inf(1)), true},
		{"-Inf", Float(math.Inf(-1)), false},
		{"NaN", Float(math.NaN()), true},
		{"-NaN", Float(math.Copysign(math.NaN(), -1)), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.IsSignPositive(); got != tc.signPositive {
				t.Errorf("IsSignPositive(%v) = %v, want %v", tc.in, got, tc.signPositive)
			}
			if got := tc.in.IsSignNegative(); got != !tc.signPositive {
				t.Errorf("IsSignNegative(%v) = %v, want %v", tc.in, got, !tc.signPositive)
			}
		})
	}
}

func TestFloatNeg(t *testing.T) {
	if Float(3.5).Neg() != -3.5 {
		t.Errorf("Neg(3.5): expected -3.5, got %v", Float(3.5).Neg())
	}

	if Float(-2).Neg() != 2 {
		t.Errorf("Neg(-2): expected 2, got %v", Float(-2).Neg())
	}

	if !math.Signbit(Float(0).Neg().Std()) {
		t.Error("Neg(+0): expected -0 (sign bit set)")
	}

	if Float(math.Inf(1)).Neg() != Float(math.Inf(-1)) {
		t.Error("Neg(+Inf): expected -Inf")
	}
}

func TestFloatMathCeilFloorTrunc(t *testing.T) {
	if Float(2.3).Ceil() != 3 || Float(-2.3).Ceil() != -2 || Float(2).Ceil() != 2 {
		t.Error("Ceil: unexpected results")
	}

	if Float(2.7).Floor() != 2 || Float(-2.7).Floor() != -3 || Float(2).Floor() != 2 {
		t.Error("Floor: unexpected results")
	}

	if Float(2.7).Trunc() != 2 || Float(-2.7).Trunc() != -2 {
		t.Error("Trunc: unexpected results")
	}
}

func TestFloatMathFract(t *testing.T) {
	if Float(3.75).Fract() != 0.75 {
		t.Errorf("Fract(3.75): expected 0.75, got %v", Float(3.75).Fract())
	}

	if Float(-3.75).Fract() != -0.75 {
		t.Errorf("Fract(-3.75): expected -0.75, got %v", Float(-3.75).Fract())
	}

	if Float(5).Fract() != 0 {
		t.Errorf("Fract(5): expected 0, got %v", Float(5).Fract())
	}

	if !Float(math.Inf(1)).Fract().IsNaN() {
		t.Error("Fract(+Inf): expected NaN")
	}

	if !Float(math.Inf(-1)).Fract().IsNaN() {
		t.Error("Fract(-Inf): expected NaN")
	}

	if !Float(math.NaN()).Fract().IsNaN() {
		t.Error("Fract(NaN): expected NaN")
	}
}

func TestFloatMathClamp(t *testing.T) {
	if Float(5).Clamp(1, 10) != 5 {
		t.Error("Clamp(5, 1, 10): expected 5")
	}

	if Float(-3).Clamp(1, 10) != 1 {
		t.Error("Clamp(-3, 1, 10): expected 1")
	}

	if Float(42).Clamp(1, 10) != 10 {
		t.Error("Clamp(42, 1, 10): expected 10")
	}

	if Float(1).Clamp(1, 10) != 1 || Float(10).Clamp(1, 10) != 10 {
		t.Error("Clamp at boundaries: expected boundary values unchanged")
	}

	if Float(3).Clamp(3, 3) != 3 {
		t.Error("Clamp(3, 3, 3): expected 3")
	}

	if !Float(math.NaN()).Clamp(1, 10).IsNaN() {
		t.Error("Clamp(NaN, 1, 10): expected NaN")
	}

	if Float(math.Inf(1)).Clamp(1, 10) != 10 || Float(math.Inf(-1)).Clamp(1, 10) != 1 {
		t.Error("Clamp(±Inf, 1, 10): expected boundary values")
	}
}

func TestFloatMathRecipCopysign(t *testing.T) {
	if Float(4).Recip() != 0.25 {
		t.Errorf("Recip(4): expected 0.25, got %v", Float(4).Recip())
	}

	if !Float(0).Recip().IsInf() {
		t.Error("Recip(0): expected +Inf")
	}

	if Float(3).Copysign(-1) != -3 {
		t.Errorf("Copysign(3, -1): expected -3, got %v", Float(3).Copysign(-1))
	}

	if Float(-3).Copysign(1) != 3 {
		t.Errorf("Copysign(-3, 1): expected 3, got %v", Float(-3).Copysign(1))
	}

	negzero := Float(math.Copysign(0, -1))
	if Float(7).Copysign(negzero) != -7 {
		t.Error("Copysign(7, -0): expected -7")
	}
}

func TestFloatMathMulAdd(t *testing.T) {
	// FMA performs a single rounding: 0.1*10 rounds to exactly 1.0 when done
	// as a separate multiply, so the naive expression yields 0, while the fused
	// operation keeps the exact product and yields the nonzero residue.
	x, y, z := Float(0.1), Float(10), Float(-1)

	naive := x*y + z
	if naive != 0 {
		t.Errorf("naive 0.1*10-1: expected 0 due to double rounding, got %v", naive)
	}

	fused := x.MulAdd(y, z)
	if fused == 0 {
		t.Error("MulAdd(0.1, 10, -1): expected nonzero residue from single rounding")
	}

	if fused.Std() != math.FMA(0.1, 10, -1) {
		t.Errorf("MulAdd: expected %v, got %v", math.FMA(0.1, 10, -1), fused)
	}

	if Float(2).MulAdd(3, 4) != 10 {
		t.Errorf("MulAdd(2, 3, 4): expected 10, got %v", Float(2).MulAdd(3, 4))
	}
}

func TestFloatMathHypotCbrt(t *testing.T) {
	if Float(3).Hypot(4) != 5 {
		t.Errorf("Hypot(3, 4): expected 5, got %v", Float(3).Hypot(4))
	}

	if Float(1e300).Hypot(1e300).IsInf() {
		t.Error("Hypot(1e300, 1e300): expected no overflow")
	}

	if Float(27).Cbrt() != 3 {
		t.Errorf("Cbrt(27): expected 3, got %v", Float(27).Cbrt())
	}

	if Float(-8).Cbrt() != -2 {
		t.Errorf("Cbrt(-8): expected -2, got %v", Float(-8).Cbrt())
	}
}

func TestFloatMathExpLog(t *testing.T) {
	if Float(0).Exp() != 1 {
		t.Error("Exp(0): expected 1")
	}

	if Float(1).Exp() != Float(math.E) {
		t.Errorf("Exp(1): expected e, got %v", Float(1).Exp())
	}

	if Float(3).Exp2() != 8 {
		t.Errorf("Exp2(3): expected 8, got %v", Float(3).Exp2())
	}

	if Float(0).ExpM1() != 0 {
		t.Error("ExpM1(0): expected 0")
	}

	if Float(1e-10).ExpM1().Std() != math.Expm1(1e-10) {
		t.Error("ExpM1(1e-10): expected math.Expm1 result")
	}

	if Float(math.E).Ln() != 1 {
		t.Errorf("Ln(e): expected 1, got %v", Float(math.E).Ln())
	}

	if !Float(-1).Ln().IsNaN() {
		t.Error("Ln(-1): expected NaN")
	}

	if !Float(0).Ln().IsInf() || Float(0).Ln() > 0 {
		t.Error("Ln(0): expected -Inf")
	}

	if Float(0).Ln1p() != 0 {
		t.Error("Ln1p(0): expected 0")
	}

	if !Float(-2).Ln1p().IsNaN() {
		t.Error("Ln1p(-2): expected NaN")
	}

	if Float(8).Log2() != 3 {
		t.Errorf("Log2(8): expected 3, got %v", Float(8).Log2())
	}

	if Float(1000).Log10() != 3 {
		t.Errorf("Log10(1000): expected 3, got %v", Float(1000).Log10())
	}
}

func TestFloatMathTrig(t *testing.T) {
	const eps = 1e-15

	if Float(0).Sin() != 0 || Float(0).Cos() != 1 || Float(0).Tan() != 0 {
		t.Error("Sin/Cos/Tan(0): unexpected results")
	}

	if math.Abs(Float(math.Pi/2).Sin().Std()-1) > eps {
		t.Error("Sin(Pi/2): expected 1")
	}

	if Float(1).Asin() != Float(math.Pi/2) {
		t.Errorf("Asin(1): expected Pi/2, got %v", Float(1).Asin())
	}

	if Float(1).Acos() != 0 {
		t.Errorf("Acos(1): expected 0, got %v", Float(1).Acos())
	}

	if !Float(2).Asin().IsNaN() || !Float(2).Acos().IsNaN() {
		t.Error("Asin/Acos(2): expected NaN")
	}

	if Float(1).Atan() != Float(math.Pi/4) {
		t.Errorf("Atan(1): expected Pi/4, got %v", Float(1).Atan())
	}

	if Float(0).Sinh() != 0 || Float(0).Cosh() != 1 || Float(0).Tanh() != 0 {
		t.Error("Sinh/Cosh/Tanh(0): unexpected results")
	}

	if Float(0).Asinh() != 0 || Float(1).Acosh() != 0 || Float(0).Atanh() != 0 {
		t.Error("Asinh(0)/Acosh(1)/Atanh(0): unexpected results")
	}

	if !Float(0.5).Acosh().IsNaN() {
		t.Error("Acosh(0.5): expected NaN")
	}

	if !Float(2).Atanh().IsNaN() {
		t.Error("Atanh(2): expected NaN")
	}
}

func TestFloatMathAtan2Quadrants(t *testing.T) {
	const eps = 1e-15

	testCases := []struct {
		y, x     Float
		expected float64
	}{
		{1, 1, math.Pi / 4},        // quadrant I
		{1, -1, 3 * math.Pi / 4},   // quadrant II
		{-1, -1, -3 * math.Pi / 4}, // quadrant III
		{-1, 1, -math.Pi / 4},      // quadrant IV
		{0, 1, 0},
		{1, 0, math.Pi / 2},
		{-1, 0, -math.Pi / 2},
	}

	for _, tc := range testCases {
		result := tc.y.Atan2(tc.x).Std()
		if math.Abs(result-tc.expected) > eps {
			t.Errorf("Atan2(%v, %v): expected %v, got %v", tc.y, tc.x, tc.expected, result)
		}
	}

	if Float(0).Atan2(-1).Std() != math.Pi {
		t.Errorf("Atan2(0, -1): expected Pi, got %v", Float(0).Atan2(-1))
	}
}

func TestFloatMathAngles(t *testing.T) {
	if Float(math.Pi).ToDegrees() != 180 {
		t.Errorf("ToDegrees(Pi): expected 180, got %v", Float(math.Pi).ToDegrees())
	}

	if Float(math.Pi/2).ToDegrees() != 90 {
		t.Errorf("ToDegrees(Pi/2): expected 90, got %v", Float(math.Pi/2).ToDegrees())
	}

	if Float(180).ToRadians() != Float(math.Pi) {
		t.Errorf("ToRadians(180): expected Pi, got %v", Float(180).ToRadians())
	}

	if Float(0).ToDegrees() != 0 || Float(0).ToRadians() != 0 {
		t.Error("ToDegrees/ToRadians(0): expected 0")
	}

	roundtrip := Float(45).ToRadians().ToDegrees()
	if math.Abs(roundtrip.Std()-45) > 1e-13 {
		t.Errorf("45 deg roundtrip: expected ~45, got %v", roundtrip)
	}
}
