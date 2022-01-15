package udefiner

import (
	"strings"
	"testing"

	"github.com/anaskhan96/soup"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestGetDefinitionsFromHTML(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(buttUrbanDictionaryHtml))
	ud, _ := getDefinitionsFromHTML(soup.Root{Pointer: node})
	assert.Equal(t, 5, len(ud), "amount of scraped definitions should be 5")
	expectedFirstDefinition := UrbanDefinition{
		"Potent weed, usually a sativa, that gives you an energetic high",
		"\"That butt gave me a crazy head high\"\"I was geeked after smoking butt yesterday\"",
		"274",
		"23",
	}
	assert.Equal(t, expectedFirstDefinition, ud[0], "first definition should be scraped properly")
}
