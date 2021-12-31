package term

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawSplitsOnColon(t *testing.T) {
	ter := Term("define: butt")
	want := "butt"
	got := ter.Raw()
	assert.Equal(t, want, got, "Raw() should split on colon")
}

func TestRawReturnsOriginalStringIfNoColon(t *testing.T) {
	ter := Term("butt")
	want := "butt"
	got := ter.Raw()
	assert.Equal(t, want, got, "Raw() returns original string if no colon present")
}

func TestStringJoinsOnHyphen(t *testing.T) {
	ter := Term("define: just long things")
	want := "just-long-things"
	got := ter.String()
	assert.Equal(t, want, got, "String() should join term on hyphens")
}
