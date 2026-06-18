package sensors

import "testing"

func TestValueValidateRejectsEmptyKind(t *testing.T) {
	if err := (Value{}).Validate(); err == nil {
		t.Fatal("expected empty value kind to be invalid")
	}
}

func TestValueValidateRejectsUnknownKind(t *testing.T) {
	if err := (Value{Kind: "banana"}).Validate(); err == nil {
		t.Fatal("expected unknown value kind to be invalid")
	}
}

func TestNumericZeroIsValidNumericValue(t *testing.T) {
	value := NewNumericValue(0, "km/h")
	if err := value.Validate(); err != nil {
		t.Fatalf("Validate numeric zero: %v", err)
	}
	number, ok := value.Numeric()
	if !ok || number != 0 {
		t.Fatalf("Numeric() = %v %v, want 0 true", number, ok)
	}
}
