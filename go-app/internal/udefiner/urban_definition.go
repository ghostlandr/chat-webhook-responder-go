package udefiner

import (
	"fmt"

	"github.com/anaskhan96/soup"
)

type UrbanDefinition struct {
	DefID     string
	meaning   string
	example   string
	upvotes   string
	downvotes string
}

func (ud UrbanDefinition) String() string {
	var thumbArea string
	if ud.upvotes != "" && ud.downvotes != "" {
		thumbArea = fmt.Sprintf(":+1: %s :-1: %s", ud.upvotes, ud.downvotes)
	}
	return fmt.Sprintf("%s\n%s", ud.meaning, thumbArea)
}

func (ud *UrbanDefinition) SetVotes(upvotes, downvotes string) {
	ud.upvotes = upvotes
	ud.downvotes = downvotes
}

// NewUrbanDefinition expects HTML that contains all the definition material for one definition
func NewUrbanDefinition(definition soup.Root) UrbanDefinition {
	if definition.Error != nil {
		return UrbanDefinition{}
	}

	defID := definition.Attrs()["data-defid"]

	meaning := definition.Find("div", "class", "meaning").FullText()

	example := definition.Find("div", "class", "example").FullText()

	// Urban Dictionary populates the up and down votes after page load, so they
	// aren't coming back when we scrape :(
	// for _, thumbArea := range thumbsArea {
	// 	switch thumbArea.Attrs()["data-direction"] {
	// 	case "up":
	// 		upvotes = thumbArea.FullText()
	// 	case "down":
	// 		downvotes = thumbArea.FullText()
	// 	}
	// }

	// fmt.Printf("Upvotes: %s, downvotes: %s\n", upvotes, downvotes)

	return UrbanDefinition{
		DefID:   defID,
		meaning: meaning,
		example: example,
	}
}
