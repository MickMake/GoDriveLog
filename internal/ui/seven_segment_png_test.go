package ui

import "testing"

func TestFormatDigitCellsAppliesDecimalToPreviousCell(t *testing.T) {
	cells := formatDigitCells("82.5", 0)
	want := []digitCell{
		{Symbol: "8"},
		{Symbol: "2", Decimal: true},
		{Symbol: "5"},
	}
	assertDigitCells(t, cells, want)
}

func TestFormatDigitCellsHandlesSignedDecimal(t *testing.T) {
	cells := formatDigitCells("-5.0", 0)
	want := []digitCell{
		{Symbol: "dash"},
		{Symbol: "5", Decimal: true},
		{Symbol: "0"},
	}
	assertDigitCells(t, cells, want)
}

func TestFormatDigitCellsPadsAndTrims(t *testing.T) {
	assertDigitCells(t, formatDigitCells("0.9", 4), []digitCell{
		{Symbol: "blank"},
		{Symbol: "blank"},
		{Symbol: "0", Decimal: true},
		{Symbol: "9"},
	})

	assertDigitCells(t, formatDigitCells("100.0", 3), []digitCell{
		{Symbol: "0"},
		{Symbol: "0", Decimal: true},
		{Symbol: "0"},
	})
}

func assertDigitCells(t *testing.T, got, want []digitCell) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cell %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}
