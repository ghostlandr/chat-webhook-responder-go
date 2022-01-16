package term

import "strings"

type Term string

func (t Term) Raw() string {
	split := strings.Split(string(t), ": ")
	if len(split) == 1 {
		return string(t)
	}
	return strings.Split(string(t), ": ")[1]
}

func (t Term) String() string {
	term := t.Raw()
	st := strings.Split(strings.ToLower(term), " ")
	return strings.Join(st, "-")
}
