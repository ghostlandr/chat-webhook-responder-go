package udefiner

import (
	"strings"
	"testing"

	"github.com/anaskhan96/soup"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestUrbanDefinition(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(buttUrbanDictionaryHtml))
	ud, _ := getDefinitionsFromHTML(soup.Root{Pointer: node})
	first := ud[0]

	assert.Equal(t, first.upvotes, "274")
	assert.Equal(t, first.downvotes, "23")
}
