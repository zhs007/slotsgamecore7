package lowcode

import "testing"

func TestIntPow(t *testing.T) {
	tests := []struct {
		name string
		base int
		exp  int
		want int
	}{
		{"2^10", 2, 10, 1024},
		{"5^0", 5, 0, 1},
		{"2^-1", 2, -1, 1},
		{"0^0", 0, 0, 1},
		{"1^1000", 1, 1000, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intPow(tt.base, tt.exp)
			if got != tt.want {
				t.Fatalf("intPow(%d, %d) = %d, want %d", tt.base, tt.exp, got, tt.want)
			}
		})
	}
}
