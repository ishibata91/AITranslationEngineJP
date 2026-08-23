package japanesetext

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{source: "Hello", want: false},
		{source: "Aあ", want: true},
		{source: "Aア", want: true},
		{source: "A漢", want: true},
		{source: "<tag>123</tag>", want: false},
	}
	for _, tt := range tests {
		if got := Contains(tt.source); got != tt.want {
			t.Errorf("Contains(%q) = %t, want %t", tt.source, got, tt.want)
		}
	}
}
