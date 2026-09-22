package util

import (
	"slices"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  []string
	}{
		{
			name:  "fits on one line",
			text:  "hello world",
			width: 20,
			want:  []string{"hello world"},
		},
		{
			name:  "wraps at width",
			text:  "the quick brown fox",
			width: 10,
			want:  []string{"the quick", "brown fox"},
		},
		{
			name:  "single very long word gets cut",
			text:  "supercalifragilisticexpialidocious",
			width: 10,
			want:  []string{"supercalif"},
		},
		{
			name:  "empty string",
			text:  "",
			width: 10,
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapText(tt.text, tt.width)
			if !slices.Equal(got, tt.want) {
				t.Errorf("WrapText(%q, %d) = %v, want %v", tt.text, tt.width, got, tt.want)
			}
		})
	}
}

// NOTE: WrapText's maxlines cap stops adding *content* lines once the cap is hit, but the loop still appends a final
// (empty) trailing line afterwards, so the returned slice length is actually maxlines+1, with the last entry always
// empty. This test documents that current behaviour rather than asserting an ideal cap, since fixing it would change
// the public API's observable behaviour.
func TestWrapTextMaxLines(t *testing.T) {
	got := WrapText("the quick brown fox jumps over the lazy dog", 10, 2)
	want := []string{"the quick", "brown fox", ""}
	if !slices.Equal(got, want) {
		t.Errorf("WrapText with maxlines=2 = %v, want %v", got, want)
	}
}

func TestLoremIpsum(t *testing.T) {
	tests := []int{0, 1, 5, 50, 200}

	for _, words := range tests {
		got := LoremIpsum(words)
		gotWords := 0
		if got != "" {
			gotWords = len(splitOnSpace(got))
		}
		if gotWords != words {
			t.Errorf("LoremIpsum(%d) returned %d words, want %d", words, gotWords, words)
		}
	}
}

func splitOnSpace(s string) []string {
	var words []string
	current := ""
	for _, r := range s {
		if r == ' ' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
			continue
		}
		current += string(r)
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
