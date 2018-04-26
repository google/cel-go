package types

import(
	"testing"
)

func TestDouble_Add(t *testing.T) {
	if !Double(4).Add(Double(-3.5)).Equal(Double(0.5)).(Bool) {
		t.Error("Adding two doubles did not match expected value.")
	}
	if !IsError(Double(-1).Add(String("-1"))) {
		t.Error("Adding non-double to double was not an error.")
	}
}

func TestDouble_Compare(t *testing.T) {
	lt := Double(-1300)
	gt := Double(204)
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

func TestDouble_ConvertToType(t *testing.T) {
	if !Double(-4.5).ConvertToType(IntType).Equal(Int(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to int")
	}
	if !Double(-4.5).ConvertToType(UintType).Equal(Uint(18446744073709551612)).(Bool) {
		t.Error("Unsuccessful type conversion to uint")
	}
	if !Double(-4.5).ConvertToType(DoubleType).Equal(Double(-4.5)).(Bool) {
		t.Error("Unsuccessful type conversion to double")
	}
	if !Double(-4.5).ConvertToType(StringType).Equal(String("-4.5")).(Bool) {
		t.Error("Unsuccessful type conversion to string")
	}
	if !Double(-4.5).ConvertToType(TypeType).Equal(DoubleType).(Bool) {
		t.Error("Unsuccessful type conversion to type")
	}
}

func TestDouble_Divide(t *testing.T) {
	if !Double(3).Divide(Double(1.5)).Equal(Double(2)).(Bool) {
		t.Error("Dividing two doubles did not match expectations.")
	}
	if !IsError(Double(1.1).Divide(Double(0))) {
		t.Error("Division by zero did not cause error.")
	}
	if !IsError(Double(1.1).Divide(IntNegOne)) {
		t.Error("Division permitted without express type-conversion.")
	}
}

func TestDouble_Equal(t *testing.T) {
	if Double(0).Equal(False).(Bool) {
		t.Error("Double equal to non-double type")
	}
}

func TestDouble_Multiply(t *testing.T) {
	if !Double(1.1).Multiply(Double(-1.2)).Equal(Double(-1.32)).(Bool) {
		t.Error("Multiplying two doubles did not match expectations.")
	}
	if !IsError(Double(1.1).Multiply(IntNegOne)) {
		t.Error("Multiplication permitted without express type-conversion.")
	}
}

func TestDouble_Negate(t *testing.T) {
	if !Double(1.1).Negate().Equal(Double(-1.1)).(Bool) {
		t.Error("Negating double value did not succeed")
	}
}

func TestDouble_Subtract(t *testing.T) {
	if !Double(4).Subtract(Double(-3.5)).Equal(Double(7.5)).(Bool) {
		t.Error("Subtracting two doubles did not match expected value.")
	}
	if !IsError(Double(1.1).Subtract(IntNegOne)) {
		t.Error("Subtraction permitted without express type-conversion.")
	}
}
