package term

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringJoinsOnHyphen(t *testing.T) {
	ter := Term("define: just long things")
	want := "just-long-things"
	got := ter.String()
	assert.Equal(t, want, got, "String() should join term on hyphens")
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		want string
		term string
	}{
		{
			name: "String() should slugify with hyphens",
			want: "just-long-things",
			term: "define: just long things",
		},
		{
			name: "String() should slugify even without a colon-delimited command",
			want: "more-long-things",
			term: "More long things",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Term(tc.term).String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRawTerms(t *testing.T) {
	tests := []struct {
		name string
		want string
		term string
	}{
		{
			name: "Raw() should split on colon",
			want: "butt",
			term: "define: butt",
		},
		{
			name: "Raw() returns original string if no colon present",
			want: "butt",
			term: "butt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Term(test.term).Raw()
			assert.Equal(t, test.want, got, test.name)
		})
	}
}
