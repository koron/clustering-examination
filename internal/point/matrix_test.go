package point

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUTMatrix(t *testing.T) {
	m := NewUTMatrix(4)
	if m.n != 4 {
		t.Fatalf("m.n mismatch: want=%d got=%d", 4, m.n)
	}
	m.Set(0, 0, 1)
	m.Set(0, 1, 2)
	m.Set(0, 2, 3)
	m.Set(0, 3, 4)
	m.Set(1, 1, 5)
	m.Set(1, 2, 6)
	m.Set(1, 3, 7)
	m.Set(2, 2, 8)
	m.Set(2, 3, 9)
	m.Set(3, 3, 10)
	if d := cmp.Diff(m.d, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}); d != "" {
		t.Errorf("m.d mismatch: -want +got\n%s", d)
	}

	m = NewUTMatrix(4)
	m.Set(0, 0, 1)
	m.Set(1, 0, 2)
	m.Set(2, 0, 3)
	m.Set(3, 0, 4)
	m.Set(1, 1, 5)
	m.Set(2, 1, 6)
	m.Set(3, 1, 7)
	m.Set(2, 2, 8)
	m.Set(3, 2, 9)
	m.Set(3, 3, 10)
	if d := cmp.Diff(m.d, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}); d != "" {
		t.Errorf("m.d mismatch: -want +got\n%s", d)
	}
}
