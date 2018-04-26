package types

import(
	"testing"
)

func TestUint_Add(t *testing.T) {
	if !Uint(4).Add(Uint(3)).Equal(Uint(7)).(Bool) {
		t.Error("Adding two uints did not match expected value.")
	}
	if !IsError(Uint(1).Add(String("-1"))) {
		t.Error("Adding non-uint to uint was not an error.")
	}
}

func TestUint_Compare(t *testing.T) {
	lt := Uint(204)
	gt := Uint(1300)
	if !lt.Compare(gt).Equal(IntNegOne).(Bool) {
		t.Error("Comparison did not yield - 1")
	}
	if !gt.Compare(lt).Equal(IntOne).(Bool) {
		t.Error("Comparison did not yield 1")
	}
	if !gt.Compare(gt).Equal(IntZero).(Bool) {
		t.Error(("Comparison did not yield 0"))
	}
	if !IsError(gt.Compare(TypeType)) {
		t.Error("Types not comparable")
	}
}

func TestUint_ConvertToType(t *testing.T) {
	if !Uint(18446744073709551612).ConvertToType(IntType).Equal(Int(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to int")
	}
	if !Uint(4).ConvertToType(UintType).Equal(Uint(4)).(Bool) {
		t.Error("Unsuccessful type conversion to uint")
	}
	if !Uint(4).ConvertToType(DoubleType).Equal(Double(4)).(Bool) {
		t.Error("Unsuccessful type conversion to double")
	}
	if !Uint(4).ConvertToType(StringType).Equal(String("4")).(Bool) {
		t.Error("Unsuccessful type conversion to string")
	}
	if !Uint(4).ConvertToType(TypeType).Equal(UintType).(Bool) {
		t.Error("Unsuccessful type conversion to type")
	}
}

func TestUint_Divide(t *testing.T) {
	if !Uint(3).Divide(Uint(2)).Equal(Uint(1)).(Bool) {
		t.Error("Dividing two uints did not match expectations.")
	}
	if !IsError(uintZero.Divide(uintZero)) {
		t.Error("Divide by zero did not cause error.")
	}
	if !IsError(Uint(1).Divide(Double(-1))) {
		t.Error("Division permitted without express type-conversion.")
	}
}

func TestUint_Equal(t *testing.T) {
	if Uint(0).Equal(False).(Bool) {
		t.Error("Uint equal to non-uint type")
	}
}

func TestUint_Modulo(t *testing.T) {
	if !Uint(21).Modulo(Uint(2)).Equal(Uint(1)).(Bool) {
		t.Error("Unexpected result from modulus operator.")
	}
	if !IsError(Uint(21).Modulo(uintZero)) {
		t.Error("Modulus by zero did not cause error.")
	}
	if !IsError(Uint(21).Modulo(IntOne)) {
		t.Error("Modulus permitted between different types without type conversion.")
	}
}

func TestUint_Multiply(t *testing.T) {
	if !Uint(2).Multiply(Uint(2)).Equal(Uint(4)).(Bool) {
		t.Error("Multiplying two values did not match expectations.")
	}
	if !IsError(Uint(1).Multiply(Double(-4.0))) {
		t.Error("Multiplication permitted without express type-conversion.")
	}
}

func TestUint_Subtract(t *testing.T) {
	if !Uint(4).Subtract(Uint(3)).Equal(Uint(1)).(Bool) {
		t.Error("Subtracting two uints did not match expected value.")
	}
	if !IsError(Uint(1).Subtract(Int(1))) {
		t.Error("Subtraction permitted without express type-conversion.")
	}
}
