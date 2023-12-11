package render

import (
	"fmt"
	"testing"
)

func Test_makeMarkdownTitleLink(t *testing.T) {
	tests := []struct {
		v   string
		exp string
	}{
		{
			v:   "(archived) Make sure `if` statements using short assignment",
			exp: "#-archived-make-sure-if-statements-using-short-assignment",
		},
		{
			v:   ":derelict_house: Interactively visualize packages",
			exp: "#-derelict_house-interactively-visualize-packages",
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			got := makeMarkdownTitleLink(tc.v)
			if got != tc.exp {
				t.Error("got", got, "exp", tc.exp)
			}
		})
	}
}
