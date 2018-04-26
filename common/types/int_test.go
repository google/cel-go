package types

import(
	"testing"
)

func TestInt_Add(t *testing.T) {
	if !Int(4).Add(Int(-3)).Equal(Int(1)).(Bool) {
		t.Error("Adding two ints did not match expected value.")
	}
	if !IsError(Int(-1).Add(String("-1"))) {
		t.Error("Adding non-int to int was not an error.")
	}
}

func TestInt_Compare(t *testing.T) {
	lt := Int(-1300)
	gt := Int(204)
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

func TestInt_ConvertToType(t *testing.T) {
	if !Int(-4).ConvertToType(IntType).Equal(Int(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to int")
	}
	if !Int(-4).ConvertToType(UintType).Equal(Uint(18446744073709551612)).(Bool) {
		t.Error("Unsuccessful type conversion to uint")
	}
	if !Int(-4).ConvertToType(DoubleType).Equal(Double(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to double")
	}
	if !Int(-4).ConvertToType(StringType).Equal(String("-4")).(Bool) {
		t.Error("Unsuccessful type conversion to string")
	}
	if !Int(-4).ConvertToType(TypeType).Equal(IntType).(Bool) {
		t.Error("Unsuccessful type conversion to type")
	}
}

func TestInt_Divide(t *testing.T) {
	if !Int(3).Divide(Int(2)).Equal(Int(1)).(Bool) {
		t.Error("Dividing two ints did not match expectations.")
	}
	if !IsError(IntZero.Divide(IntZero)) {
		t.Error("Divide by zero did not cause error.")
	}
	if !IsError(Int(1).Divide(Double(-1))) {
		t.Error("Division permitted without express type-conversion.")
	}
}

func TestInt_Equal(t *testing.T) {
	if Int(0).Equal(False).(Bool) {
		t.Error("Int equal to non-int type")
	}
}

func TestInt_Modulo(t *testing.T) {
	if !Int(21).Modulo(Int(2)).Equal(Int(1)).(Bool) {
		t.Error("Unexpected result from modulus operator.")
	}
	if !IsError(Int(21).Modulo(IntZero)) {
		t.Error("Modulus by zero did not cause error.")
	}
	if !IsError(Int(21).Modulo(uintZero)) {
		t.Error("Modulus permitted between different types without type conversion.")
	}
}

func TestInt_Multiply(t *testing.T) {
	if !Int(2).Multiply(Int(-2)).Equal(Int(-4)).(Bool) {
		t.Error("Multiplying two values did not match expectations.")
	}
	if !IsError(Int(1).Multiply(Double(-4.0))) {
		t.Error("Multiplication permitted without express type-conversion.")
	}
}

func TestInt_Negate(t *testing.T) {
	if !Int(1).Negate().Equal(Int(-1)).(Bool) {
		t.Error("Negating int value did not succeed")
	}
}

func TestInt_Subtract(t *testing.T) {
	if !Int(4).Subtract(Int(-3)).Equal(Int(7)).(Bool) {
		t.Error("Subtracting two ints did not match expected value.")
	}
	if !IsError(Int(1).Subtract(Uint(1))) {
		t.Error("Subtraction permitted without express type-conversion.")
	}
}
