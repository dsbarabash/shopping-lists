package model

type State int

const (
	Archived = iota
	Active
)

var stateName = map[State]string{
	Active:   "Active",
	Archived: "Archived",
}

func (s State) String() string {
	return stateName[s]
}
