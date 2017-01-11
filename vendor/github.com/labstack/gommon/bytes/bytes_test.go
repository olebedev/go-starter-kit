package bytes

import "testing"

func TestBytes(t *testing.T) {
	// B
	b := Format(515)
	if b != "515 B" {
		t.Errorf("expected `515 B`, got %s", b)
	}

	// MB
	b = Format(13231323)
	if b != "13.23 MB" {
		t.Errorf("expected `13.23 MB`, got %s", b)
	}

	// Exact
	b = Format(1000 * 1000 * 1000)
	if b != "1.00 GB" {
		t.Errorf("expected `1.00 GB`, got %s", b)
	}

	// Binary prefix
	BinaryPrefix(true)
	b = Format(1323)
	if b != "1.29 KiB" {
		t.Errorf("expected `1.29 KiB`, got %s", b)
	}
}
